package services

import (
	"goweb/models"
	"goweb/requests"

	"gorm.io/gorm"
)

type PostService struct {
	DB *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{DB: db}
}

func (postService *PostService) Create(post *models.Post) {
	postService.DB.Create(post)
}

func (postService *PostService) Delete(post *models.Post) {
	postService.DB.Delete(post)
}

func (postService *PostService) Update(post *models.Post, updatePostRequest *requests.UpdatePostRequest) {
	post.Content = updatePostRequest.Content
	post.Title = updatePostRequest.Title
	postService.DB.Save(post)
}

func (postService *PostService) GetPostsInDomain(posts *[]models.Post, domain *models.Domain) {
	postService.DB.Preload("User").Where("domain = ?", domain.UUID).Find(posts)
}

func (postService *PostService) GetPostByUuidInDomain(post *models.Post, uuid string, domain *models.Domain) {
	postService.DB.Preload("User").Where("uuid = ? AND domain = ?", uuid, domain.UUID).First(post)
}
