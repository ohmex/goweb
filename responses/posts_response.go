package responses

import (
	"goweb/models"
)

type PostResponse struct {
	UUID     string `json:"uuid" example:"uuid"`
	Username string `json:"username" example:"John Doe"`
	Title    string `json:"title" example:"Echo"`
	Content  string `json:"content" example:"Echo is nice!"`
}

func NewPostResponse(posts []models.Post) *[]PostResponse {
	postResponse := make([]PostResponse, 0)

	for i := range posts {
		postResponse = append(postResponse, PostResponse{
			UUID:     posts[i].UUID.String(),
			Username: posts[i].User.Name,
			Title:    posts[i].Title,
			Content:  posts[i].Content,
		})
	}

	return &postResponse
}
