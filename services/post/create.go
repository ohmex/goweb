package post

import "goweb/models"

func (postService *Service) Create(post *models.Post) {
	postService.DB.Create(post)
}
