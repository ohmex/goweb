package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"goweb/api"
	"goweb/models"
	"goweb/requests"
	"goweb/server"
	"goweb/util"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserService struct {
	DB    *gorm.DB
	Redis *server.Server
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

// SetRedis sets the Redis client for caching
func (service *UserService) SetRedis(server *server.Server) {
	service.Redis = server
}

// Cache keys
const (
	UserCacheKey       = "user:%d"
	UserEmailCacheKey  = "user:email:%s"
	UserUUIDCacheKey   = "user:uuid:%s"
	UserDomainCacheKey = "users:domain:%d"
	CacheTTL           = 15 * time.Minute
)

// getCachedUser retrieves user from cache
func (service *UserService) getCachedUser(key string) (*models.User, error) {
	if service.Redis == nil {
		return nil, errors.New("redis not configured")
	}

	ctx := context.Background()
	data, err := service.Redis.Redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// setCachedUser stores user in cache
func (service *UserService) setCachedUser(key string, user *models.User) error {
	if service.Redis == nil {
		return nil
	}

	ctx := context.Background()
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return service.Redis.Redis.Set(ctx, key, data, CacheTTL).Err()
}

// invalidateUserCache removes user from cache
func (service *UserService) invalidateUserCache(userID uint64, email, uuid string) {
	if service.Redis == nil {
		return
	}

	ctx := context.Background()
	keys := []string{
		fmt.Sprintf(UserCacheKey, userID),
		fmt.Sprintf(UserEmailCacheKey, email),
		fmt.Sprintf(UserUUIDCacheKey, uuid),
	}

	for _, key := range keys {
		service.Redis.Redis.Del(ctx, key)
	}
}

func (service *UserService) Register(e echo.Context, request *requests.RegisterRequest, domain *models.Domain) error {
	user := models.User{}

	service.GetUserByEmail(&user, request.Email)

	if user.ID != 0 {
		return api.WebResponse(e, http.StatusBadRequest, api.USER_EXISTS())
	}

	encryptedPassword, err := util.HashPassword(request.Password)

	if err != nil {
		return api.WebResponse(e, http.StatusInternalServerError, api.RESOURCE_CREATION_FAILED("Resource creation failed - error creating password"))
	}

	newUser := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: encryptedPassword,
		Domains:  []*models.Domain{domain},
	}

	ok := service.DB.Create(&newUser).Error

	if ok != nil {
		return api.WebResponse(e, http.StatusInternalServerError, api.RESOURCE_CREATION_FAILED())
	}

	service.DB.Save(&models.DomainUser{UserID: newUser.ID, DomainID: domain.ID, Active: true})

	return api.WebResponse(e, http.StatusCreated, api.RESOURCE_CREATED("User created"))
}

func (service *UserService) UpdateUser(user *models.User) error {
	err := service.DB.Save(user).Error
	if err == nil {
		// Invalidate cache after update
		service.invalidateUserCache(user.ID, user.Email, user.UUID.String())
	}
	return err
}

func (service *UserService) GetUser(user *models.User, id int) error {
	return service.DB.
		First(user, id).Error
}

func (service *UserService) GetUserByUUID(user *models.User, uuid string) error {
	// Try cache first
	if cachedUser, err := service.getCachedUser(fmt.Sprintf(UserUUIDCacheKey, uuid)); err == nil {
		*user = *cachedUser
		return nil
	}

	// Fallback to database
	err := service.DB.
		Where("uuid = ?", uuid).
		First(user).Error

	// Cache the result if found
	if err == nil && user.ID != 0 {
		service.setCachedUser(fmt.Sprintf(UserUUIDCacheKey, uuid), user)
		service.setCachedUser(fmt.Sprintf(UserCacheKey, user.ID), user)
		service.setCachedUser(fmt.Sprintf(UserEmailCacheKey, user.Email), user)
	}

	return err
}

func (service *UserService) GetUserByEmail(user *models.User, email string) error {
	// Try cache first
	if cachedUser, err := service.getCachedUser(fmt.Sprintf(UserEmailCacheKey, email)); err == nil {
		*user = *cachedUser
		return nil
	}

	// Fallback to database
	err := service.DB.
		Where("email = ?", email).
		First(user).Error

	// Cache the result if found
	if err == nil && user.ID != 0 {
		service.setCachedUser(fmt.Sprintf(UserEmailCacheKey, email), user)
		service.setCachedUser(fmt.Sprintf(UserCacheKey, user.ID), user)
		service.setCachedUser(fmt.Sprintf(UserUUIDCacheKey, user.UUID.String()), user)
	}

	return err
}

func (service *UserService) GetUsersInDomain(users *[]*models.User, domain *models.Domain) error {
	// Try cache first
	if service.Redis != nil {
		cacheKey := fmt.Sprintf(UserDomainCacheKey, domain.ID)
		ctx := context.Background()
		if data, err := service.Redis.Redis.Get(ctx, cacheKey).Result(); err == nil {
			if err := json.Unmarshal([]byte(data), users); err == nil {
				return nil
			}
		}
	}

	// Fallback to database with optimized query
	err := service.DB.
		Select("users.*").
		Joins("JOIN domain_users ON domain_users.user_id = users.id").
		Where("domain_users.domain_id = ? AND domain_users.active = ?", domain.ID, true).
		Find(users).Error

	// Cache the result if successful
	if err == nil && service.Redis != nil {
		if data, marshalErr := json.Marshal(users); marshalErr == nil {
			cacheKey := fmt.Sprintf(UserDomainCacheKey, domain.ID)
			ctx := context.Background()
			service.Redis.Redis.Set(ctx, cacheKey, data, CacheTTL)
		}
	}

	return err
}

func (service *UserService) GetUserByUuidInDomain(user *models.User, uuid string, domain *models.Domain) error {
	return service.DB.
		Joins("JOIN domain_users ON domain_users.user_id = users.id").
		Where("domain_users.domain_id = ?", domain.ID).
		Where("uuid = ?", uuid).
		First(user).Error
}

func (service *UserService) DeleteUserByUuidInDomain(user *models.User, uuid string, domain *models.Domain) error {
	// delete the user only if User.uuid == uuid & User.domains contains selected domain
	// get the domains of the user that we want to delete
	// check if the selected domain is present within domains above
	service.DB.
		Joins("JOIN domain_users ON domain_users.user_id = users.id").
		Where("domain_users.domain_id = ?", domain.ID).
		Where("uuid = ?", uuid).
		Preload("Domains").
		First(user)
	if user.ID != 0 {
		for _, d := range user.Domains {
			if d.UUID.String() == domain.UUID.String() {
				service.DB.Delete(user)
				return nil
			}
		}
	}
	return errors.New("user not found")
}
