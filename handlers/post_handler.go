package handlers

import (
	"goweb/api"
	"goweb/server"
	"net/http"

	"goweb/models"
	"goweb/requests"
	"goweb/responses"
	"goweb/services"
	"goweb/util"

	"github.com/labstack/echo/v4"
)

// PostHandler provides endpoints for managing posts, including CRUD operations.
type PostHandler struct {
	server      *server.Server
	postService *services.PostService
}

// NewPostHandler initializes the PostHandler with the provided server.
func NewPostHandler(server *server.Server) *PostHandler {
	return &PostHandler{
		server:      server,
		postService: services.NewPostService(server.DB),
	}
}

// Type returns the string identifier for the PostHandler.
func (u PostHandler) Type() string {
	return "Post"
}

// List godoc
// @Summary List posts
// @Description Returns a list of posts for the specified domain.
// @ID post-list
// @Tags Post Management
// @Accept json
// @Produce json
// @Success 200 {array} responses.PostResponse
// @Failure 400 {object} api.Response
// @Router /posts [get]
func (h *PostHandler) List(e echo.Context) error {
	d, err := util.ExtractDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	domain, _ := d.(*models.Domain)
	var posts []models.Post
	h.postService.GetPostsInDomain(&posts, domain)
	return api.WebResponse(e, http.StatusOK, responses.NewPostResponse(posts))
}

// Create godoc
// @Summary Create post
// @Description Creates a new post in the specified domain.
// @ID post-create
// @Tags Post Management
// @Accept json
// @Produce json
// @Param params body requests.CreatePostRequest true "Post creation data"
// @Success 201 {object} api.Response
// @Failure 400 {object} api.Response
// @Router /posts [post]
func (h *PostHandler) Create(e echo.Context) error {
	createRequest, err := util.BindAndValidate[requests.CreatePostRequest](e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	d, err := util.ExtractDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	domain, _ := d.(*models.Domain)
	userVal := e.Get("user")
	user, ok := userVal.(*models.User)
	if !ok || user == nil {
		return api.WebResponse(e, http.StatusUnauthorized, api.FIELD_VALIDATION_ERROR("User not found in context"))
	}
	post := &models.Post{
		BaseResource: models.BaseResource{Domain: domain.UUID},
		Title:        createRequest.Title,
		Content:      createRequest.Content,
		UserID:       int(user.ID),
	}
	h.postService.Create(post)
	return api.WebResponse(e, http.StatusCreated, api.RESOURCE_CREATED("Post created successfully"))
}

// Read godoc
// @Summary Get post
// @Description Returns the details of a post by UUID within the specified domain.
// @ID post-read
// @Tags Post Management
// @Accept json
// @Produce json
// @Param uuid path string true "Post UUID"
// @Success 200 {object} responses.PostResponse
// @Failure 400 {object} api.Response
// @Failure 404 {object} api.Response
// @Router /posts/{uuid} [get]
func (h *PostHandler) Read(e echo.Context) error {
	d, err := util.ExtractDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	domain, _ := d.(*models.Domain)
	uuid, err := util.GetUUIDParam(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, err)
	}
	var post models.Post
	h.postService.GetPostByUuidInDomain(&post, uuid, domain)
	if post.ID == 0 {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Post not found"))
	}
	return api.WebResponse(e, http.StatusOK, responses.NewPostResponse([]models.Post{post}))
}

// Update godoc
// @Summary Update post
// @Description Modifies the details of a post by UUID within the specified domain.
// @ID post-update
// @Tags Post Management
// @Accept json
// @Produce json
// @Param uuid path string true "Post UUID"
// @Param params body requests.UpdatePostRequest true "Post update data"
// @Success 200 {object} responses.PostResponse
// @Failure 400 {object} api.Response
// @Failure 404 {object} api.Response
// @Router /posts/{uuid} [put]
func (h *PostHandler) Update(e echo.Context) error {
	updateRequest, err := util.BindAndValidate[requests.UpdatePostRequest](e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	d, err := util.ExtractDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	domain, _ := d.(*models.Domain)
	uuid, err := util.GetUUIDParam(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, err)
	}
	var post models.Post
	h.postService.GetPostByUuidInDomain(&post, uuid, domain)
	if post.ID == 0 {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Post not found"))
	}
	h.postService.Update(&post, updateRequest)
	return api.WebResponse(e, http.StatusOK, responses.NewPostResponse([]models.Post{post}))
}

// Delete godoc
// @Summary Delete post
// @Description Removes a post by UUID from the specified domain.
// @ID post-delete
// @Tags Post Management
// @Accept json
// @Produce json
// @Param uuid path string true "Post UUID"
// @Success 200 {object} api.Response
// @Failure 400 {object} api.Response
// @Failure 404 {object} api.Response
// @Router /posts/{uuid} [delete]
func (h *PostHandler) Delete(e echo.Context) error {
	d, err := util.ExtractDomain(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, api.FIELD_VALIDATION_ERROR())
	}
	domain, _ := d.(*models.Domain)
	uuid, err := util.GetUUIDParam(e)
	if err != nil {
		return api.WebResponse(e, http.StatusBadRequest, err)
	}
	var post models.Post
	h.postService.GetPostByUuidInDomain(&post, uuid, domain)
	if post.ID == 0 {
		return api.WebResponse(e, http.StatusNotFound, api.RESOURCE_NOT_FOUND("Post not found"))
	}
	h.postService.Delete(&post)
	return api.WebResponse(e, http.StatusOK, api.RESOURCE_DELETED("Post deleted successfully"))
}
