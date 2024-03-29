package gin

import (
	"context"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postDelivery "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/http/gin"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userDelivery "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location = "pkg.dependency.delivery.gin."
)

type GinFactory struct {
	Gin    config.Gin
	Server *http.Server
	Router *gin.Engine
}

const (
	shutDownCompleted = "Server connection has been successfully closed..."
)

func (ginFactory *GinFactory) InitializeServer(serverConfig applicationModel.ServerRouters) {
	applicationConfig := config.AppConfig
	ginFactory.Router = gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{applicationConfig.Gin.AllowOrigins}
	corsConfig.AllowCredentials = applicationConfig.Gin.AllowCredentials
	router := ginFactory.Router.Group(applicationConfig.Gin.ServerGroup)
	ginFactory.Router.Use(cors.New(corsConfig))

	// Routers
	serverConfig.UserRouter.UserRouter(router, serverConfig.UserUseCase)
	serverConfig.PostRouter.PostRouter(router, serverConfig.UserUseCase)

	ginFactory.Server = &http.Server{
		Addr:    ":" + ginFactory.Gin.Port,
		Handler: ginFactory.Router,
	}
}

func (ginFactory *GinFactory) LaunchServer(ctx context.Context, container *applicationModel.Container) {
	applicationConfig := config.AppConfig
	runError := ginFactory.Router.Run(":" + applicationConfig.Gin.Port)
	if validator.IsErrorNotNil(runError) {
		container.RepositoryFactory.CloseRepository(ctx)
		runInternalError := domainError.NewInternalError(location+"LaunchServer.Router.Run", runError.Error())
		logging.Logger(runInternalError)
	}
}

func (ginFactory *GinFactory) CloseServer(ctx context.Context) {
	shutDownError := ginFactory.Server.Shutdown(ctx)
	if validator.IsErrorNotNil(shutDownError) {
		shutDownInternalError := domainError.NewInternalError(location+"CloseServer.Server.Shutdown", shutDownError.Error())
		logging.Logger(shutDownInternalError)
	}
	logging.Logger(shutDownCompleted)
}

func (ginFactory *GinFactory) NewUserController(domain interface{}) user.UserController {
	userUseCase := domain.(user.UserUseCase)
	return userDelivery.NewUserController(userUseCase)
}

func (ginFactory *GinFactory) NewUserRouter(controller interface{}) user.UserRouter {
	userController := controller.(user.UserController)
	return userDelivery.NewUserRouter(userController)
}

func (ginFactory *GinFactory) NewPostController(domain interface{}) post.PostController {
	postUseCase := domain.(post.PostUseCase)
	return postDelivery.NewPostController(postUseCase)
}

func (ginFactory *GinFactory) NewPostRouter(controller interface{}) post.PostRouter {
	postController := controller.(post.PostController)
	return postDelivery.NewPostRouter(postController)
}
