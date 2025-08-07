package services

import (
	"context"
	"encoding/json"
	"fmt"
	"goweb/api"
	"goweb/models"
	"goweb/server"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

const ExpireAccessMinutes = 30
const ExpireRefreshMinutes = 2 * 60
const AutoLogoffMinutes = 10

const TokenUserCacheKey = "tokens:user:%d"

type Domain struct {
	UUID string
	Name string
}

type JwtCustomClaims struct {
	UUID     string   `json:"uuid"`
	UserID   uint64   `json:"userid"`
	UserUUID string   `json:"useruuid"`
	UserName string   `json:"username"`
	Domains  []Domain `json:"domains"`
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
	if !domainsLoaded(user) {
		tokenService.server.DB.Preload("Domains").First(user)
	}

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

	tokenService.server.Redis.Set(context.Background(), fmt.Sprintf(TokenUserCacheKey, user.ID), string(cacheJSON), time.Minute*AutoLogoffMinutes)

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
		cacheJSON, _ := tokenService.server.Redis.Get(context.Background(), fmt.Sprintf(TokenUserCacheKey, claims.UserID)).Result()
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
		err := tokenService.server.DB.Preload("Domains").First(user, claims.UserID).Error
		if err != nil || user.ID == 0 {
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

	if !domainsLoaded(user) {
		tokenService.server.DB.Preload("Domains").First(user)
	}

	var domains []Domain
	for _, e := range user.Domains {
		domains = append(domains, Domain{e.UUID.String(), e.Name})
	}

	claims := &JwtCustomClaims{
		tokenUuid,
		user.ID,
		user.UUID.String(),
		user.Name,
		domains,
		jwt.RegisteredClaims{
			Issuer:    "Ohmex",
			ExpiresAt: jwt.NewNumericDate(expiry),
		},
	}
	exp = expiry.Unix()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err = jwtToken.SignedString([]byte(secret))

	return
}

func domainsLoaded(user *models.User) bool {
	return len(user.Domains) > 0 && user.Domains[0] != nil && user.Domains[0].UUID.String() != ""
}
