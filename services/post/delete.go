package post

import "goweb/models"

func (postService *Service) Delete(post *models.Post) {
	postService.DB.Delete(post)
}
