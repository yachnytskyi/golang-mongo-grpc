package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	httpGinUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin/utility"
)

type UserRouter struct {
	userHandler UserHandler
}

func NewUserRouter(userHandler UserHandler) UserRouter {
	return UserRouter{userHandler: userHandler}
}

func (userRouter *UserRouter) UserRouter(routerGroup *gin.RouterGroup, userUseCase user.UseCase) {
	router := routerGroup.Group("/users")

	router.POST("/register", userRouter.userHandler.Register)
	router.POST("/login", userRouter.userHandler.Login)

	router.Use(httpGinUtility.DeserializeUser(userUseCase))
	router.POST("/forgotten-password", userRouter.userHandler.ForgottenPassword)
	router.PATCH("/reset-password/:resetToken", userRouter.userHandler.ResetUserPassword)

	router.GET("/refresh", userRouter.userHandler.RefreshAccessToken)
	router.GET("/logout", userRouter.userHandler.Logout)

	router.GET("/me", userRouter.userHandler.GetMe)
	router.PUT("/update", userRouter.userHandler.UpdateUserById)
	router.DELETE("/delete", userRouter.userHandler.Delete)
}
