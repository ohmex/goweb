package services

import (
	"goweb/models"
	"goweb/requests"

	"gorm.io/gorm"
)

func NewPostService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

func (postService *Service) Create(post *models.Post) {
	postService.DB.Create(post)
}

func (postService *Service) Delete(post *models.Post) {
	postService.DB.Delete(post)
}

func (postService *Service) Update(post *models.Post, updatePostRequest *requests.UpdatePostRequest) {
	post.Content = updatePostRequest.Content
	post.Title = updatePostRequest.Title
	postService.DB.Save(post)
}
