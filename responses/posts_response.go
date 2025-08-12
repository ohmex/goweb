package responses

import (
	"goweb/models"
)

type PostResponse struct {
	UUID      string `json:"uuid" example:"uuid"`
	Username  string `json:"username" example:"John Doe"`
	Title     string `json:"title" example:"Echo"`
	Content   string `json:"content" example:"Echo is nice!"`
	CreatedAt string `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

func NewPostResponse(posts []models.Post) *[]PostResponse {
	postResponse := make([]PostResponse, 0)

	for i := range posts {
		postResponse = append(postResponse, PostResponse{
			UUID:      posts[i].UUID.String(),
			Username:  posts[i].User.Name,
			Title:     posts[i].Title,
			Content:   posts[i].Content,
			CreatedAt: posts[i].CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: posts[i].UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &postResponse
}
