package services

import (
	"context"
	"encoding/json"
	"fmt"
	"goweb/models"
	"goweb/requests"
	"goweb/server"
	"time"

	"gorm.io/gorm"
)

type PostService struct {
	DB    *gorm.DB
	Redis *server.Server
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{DB: db}
}

func (postService *PostService) SetRedis(server *server.Server) {
	postService.Redis = server
}

const (
	PostDomainCacheKey = "posts:domain:%s"
	PostCacheTTL       = 15 * time.Minute
)

func (postService *PostService) invalidatePostsCache(domainUUID string) {
	if postService.Redis == nil {
		return
	}
	ctx := context.Background()
	key := fmt.Sprintf(PostDomainCacheKey, domainUUID)
	postService.Redis.Redis.Del(ctx, key)
}

func (postService *PostService) Create(post *models.Post) {
	postService.DB.Create(post)
	postService.invalidatePostsCache(post.Domain.String())
}

func (postService *PostService) Delete(post *models.Post) {
	postService.DB.Delete(post)
	postService.invalidatePostsCache(post.Domain.String())
}

func (postService *PostService) Update(post *models.Post, updatePostRequest *requests.UpdatePostRequest) {
	post.Content = updatePostRequest.Content
	post.Title = updatePostRequest.Title
	postService.DB.Save(post)
	postService.invalidatePostsCache(post.Domain.String())
}

func (postService *PostService) GetPostsInDomain(posts *[]models.Post, domain *models.Domain) {
	if postService.Redis != nil {
		cacheKey := fmt.Sprintf(PostDomainCacheKey, domain.UUID.String())
		ctx := context.Background()
		if data, err := postService.Redis.Redis.Get(ctx, cacheKey).Result(); err == nil {
			if err := json.Unmarshal([]byte(data), posts); err == nil {
				return
			}
		}
	}
	postService.DB.Preload("User").Where("domain = ?", domain.UUID).Find(posts)
	if postService.Redis != nil {
		if data, err := json.Marshal(posts); err == nil {
			cacheKey := fmt.Sprintf(PostDomainCacheKey, domain.UUID.String())
			ctx := context.Background()
			postService.Redis.Redis.Set(ctx, cacheKey, data, PostCacheTTL)
		}
	}
}

func (postService *PostService) GetPostByUuidInDomain(post *models.Post, uuid string, domain *models.Domain) {
	postService.DB.Preload("User").Where("uuid = ? AND domain = ?", uuid, domain.UUID).First(post)
}
