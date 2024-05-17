package token

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"goweb/models"
	"goweb/repositories"
	"goweb/server"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

const ExpireAccessMinutes = 30
const ExpireRefreshMinutes = 2 * 60
const AutoLogoffMinutes = 10

type JwtCustomClaims struct {
	ID   uint   `json:"id"`
	UUID string `json:"uuid"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

type CachedTokens struct {
	AccessUID  string `json:"access"`
	RefreshUID string `json:"refresh"`
}

type ServiceWrapper interface {
	GenerateTokenPair(user *models.User) (accessToken, refreshToken string, exp int64, err error)
	ParseToken(tokenString, secret string) (claims *jwt.MapClaims, err error)
	ValidateToken(claims *JwtCustomClaims, isRefresh bool) error
}

type Service struct {
	server *server.Server
}

func NewTokenService(server *server.Server) *Service {
	return &Service{
		server: server,
	}
}

func (tokenService *Service) GenerateTokenPair(user *models.User) (accessToken, refreshToken string, exp int64, err error) {
	var accessUID, refreshUID string

	if accessToken, accessUID, exp, err = tokenService.createToken(user, ExpireAccessMinutes,
		tokenService.server.Config.Auth.AccessSecret); err != nil {
		return
	}

	if refreshToken, refreshUID, _, err = tokenService.createToken(user, ExpireRefreshMinutes,
		tokenService.server.Config.Auth.RefreshSecret); err != nil {
		return
	}

	cacheJSON, err := json.Marshal(CachedTokens{
		AccessUID:  accessUID,
		RefreshUID: refreshUID,
	})

	tokenService.server.Redis.Set(context.Background(), fmt.Sprintf("token-%d", user.ID), string(cacheJSON), time.Minute*AutoLogoffMinutes)

	return
}

func (tokenService *Service) ParseToken(tokenString, secret string) (claims *JwtCustomClaims, err error) {
	str := "ParseToken"
	log.Info().Msg("Inside function" + str)
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
	if err != nil {
		return
	}

	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func (tokenService *Service) ValidateToken(claims *JwtCustomClaims, isRefresh bool) (user *models.User, err error) {
	var g errgroup.Group

	g.Go(func() error {
		cacheJSON, _ := tokenService.server.Redis.Get(context.Background(), fmt.Sprintf("token-%d", claims.ID)).Result()
		cachedTokens := new(CachedTokens)
		err = json.Unmarshal([]byte(cacheJSON), cachedTokens)

		var tokenUID string
		if isRefresh {
			tokenUID = cachedTokens.RefreshUID
		} else {
			tokenUID = cachedTokens.AccessUID
		}

		if err != nil || tokenUID != claims.UUID {
			return errors.New("token not found")
		}

		return nil
	})

	g.Go(func() error {
		user = new(models.User)
		userRepository := repositories.NewUserRepository(tokenService.server.DB)
		userRepository.GetUser(user, int(claims.ID))
		if user.ID == 0 {
			return errors.New("user not found")
		}

		return nil
	})

	err = g.Wait()

	return user, err
}

func (tokenService *Service) createToken(user *models.User, expireMinutes int, secret string) (token, uid string, exp int64, err error) {
	expiry := time.Now().Add(time.Minute * time.Duration(expireMinutes))
	uid = uuid.New().String()

	claims := &JwtCustomClaims{
		user.ID,
		uid,
		user.Name,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
		},
	}
	exp = expiry.Unix()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = jwtToken.SignedString([]byte(secret))

	return
}
