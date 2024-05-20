package services

import (
	"context"
	"encoding/json"
	"fmt"
	"goweb/api"
	"goweb/models"
	"goweb/repositories"
	"goweb/server"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

const ExpireAccessMinutes = 30
const ExpireRefreshMinutes = 2 * 60
const AutoLogoffMinutes = 10

type Domain struct {
	UUID string
	Name string
}

type JwtCustomClaims struct {
	ID      uint     `json:"id"`
	UUID    string   `json:"uuid"`
	Name    string   `json:"name"`
	Tenants []Domain `json:"tenants"`
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

type TokenService struct {
	server *server.Server
}

func NewTokenService(server *server.Server) *TokenService {
	return &TokenService{
		server: server,
	}
}

func (tokenService *TokenService) GenerateTokenPair(user *models.User) (accessToken, refreshToken string, exp int64, err error) {
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

func (tokenService *TokenService) ParseToken(tokenString, secret string) (claims *JwtCustomClaims, err error) {
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

func (tokenService *TokenService) ValidateToken(claims *JwtCustomClaims, isRefresh bool) (user *models.User, err error) {
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
			return api.TOKEN_EXPIRED()
		}

		return nil
	})

	g.Go(func() error {
		user = new(models.User)
		userRepository := repositories.NewUserRepository(tokenService.server.DB)
		userRepository.GetUser(user, int(claims.ID))
		if user.ID == 0 {
			return api.USER_NOT_FOUND()
		}

		return nil
	})

	err = g.Wait()

	return user, err
}

func (tokenService *TokenService) createToken(user *models.User, expireMinutes int, secret string) (token, tokenUuid string, exp int64, err error) {
	expiry := time.Now().Add(time.Minute * time.Duration(expireMinutes))
	tokenUuid = uuid.New().String()

	tokenService.server.DB.Preload("Tenants").Where(user).Find(user)

	var tenants []Domain
	for _, e := range user.Tenants {
		tenants = append(tenants, Domain{e.Name, e.UUID.String()})
	}

	claims := &JwtCustomClaims{
		user.ID,
		tokenUuid,
		user.Name,
		tenants,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
		},
	}
	exp = expiry.Unix()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err = jwtToken.SignedString([]byte(secret))

	return
}
