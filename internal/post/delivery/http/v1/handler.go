package v1

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	"github.com/yachnytskyi/golang-mongo-grpc/models"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utils"
)

type PostHandler struct {
	postService post.Service
}

func NewPostHandler(postService post.Service) PostHandler {
	return PostHandler{postService: postService}
}

func (postHandler *PostHandler) GetAllPosts(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	intPage, err := strconv.Atoi(page)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	intLimit, err := strconv.Atoi(limit)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	fetchedPosts, err := postHandler.postService.GetAllPosts(ctx, intPage, intLimit)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": fetchedPosts})
}

func (postHandler *PostHandler) GetPostById(ctx *gin.Context) {
	postID := ctx.Param("postID")

	fetchedPost, err := postHandler.postService.GetPostById(ctx, postID)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": fetchedPost})
}

func (postHandler *PostHandler) CreatePost(ctx *gin.Context) {
	var post *models.PostCreate = new(models.PostCreate)
	currentUser := ctx.MustGet("currentUser").(*models.UserDB)
	post.User = currentUser.Name
	post.UserID = currentUser.UserID

	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
	}

	createdPost, err := postHandler.postService.CreatePost(ctx, post)

	if err != nil {
		if strings.Contains(err.Error(), "sorry, but this title already exists. Please choose another one") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": createdPost})
}

func (postHandler *PostHandler) UpdatePostById(ctx *gin.Context) {
	postID := ctx.Param("postID")
	currentUserID := utils.GetCurrentUserID(ctx)

	var updatedPostData *models.PostUpdate = new(models.PostUpdate)

	if err := ctx.ShouldBindJSON(&updatedPostData); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	updatedPost, err := postHandler.postService.UpdatePostById(ctx, postID, updatedPostData, currentUserID)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "sorry, but you do not have permissions to do that") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedPost})
}

func (postHandler *PostHandler) DeletePostByID(ctx *gin.Context) {
	postID := ctx.Param("postID")
	currentUserID := utils.GetCurrentUserID(ctx)

	err := postHandler.postService.DeletePostByID(ctx, postID, currentUserID)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "sorry, but you do not have permissions to do that") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
