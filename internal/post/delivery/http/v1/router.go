package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/middleware"
	// "github.com/yachnytskyi/golang-mongo-grpc/pkg/middleware"
)

type PostRouter struct {
	postHandler PostHandler
}

func NewPostRouter(postHandler PostHandler) PostRouter {
	return PostRouter{postHandler: postHandler}
}

func (postRouter *PostRouter) PostRouter(routerGroup *gin.RouterGroup, userUseCase user.UseCase) {
	router := routerGroup.Group("/posts")

	router.GET("/", postRouter.postHandler.GetAllPosts)
	router.GET("/:postID", postRouter.postHandler.GetPostById)

	router.Use(middleware.DeserializeUser(userUseCase))

	router.POST("/", postRouter.postHandler.CreatePost)
	router.PUT("/:postID", postRouter.postHandler.UpdatePostById)
	router.DELETE("/:postID", postRouter.postHandler.DeletePostByID)
}
