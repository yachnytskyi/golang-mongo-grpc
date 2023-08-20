package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	httpGinUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin/utility"
)

type UserRouter struct {
	userController UserController
}

func NewUserRouter(userController UserController) UserRouter {
	return UserRouter{userController: userController}
}

func (userRouter *UserRouter) UserRouter(routerGroup *gin.RouterGroup, userUseCase user.UseCase) {
	router := routerGroup.Group("/users")

	router.POST("/register", userRouter.userController.Register)
	router.POST("/login", userRouter.userController.Login)
	router.GET("/", userRouter.userController.GetAllUsers)
	router.GET("/:userID", userRouter.userController.GetUserById)

	router.Use(httpGinUtility.DeserializeUser(userUseCase))
	router.GET("/me", userRouter.userController.GetMe)
	router.PUT("/update", userRouter.userController.UpdateUserById)
	router.DELETE("/delete", userRouter.userController.Delete)
	router.GET("/refresh", userRouter.userController.RefreshAccessToken)
	router.GET("/logout", userRouter.userController.Logout)
	router.POST("/forgotten-password", userRouter.userController.ForgottenPassword)
	router.PATCH("/reset-password/:resetToken", userRouter.userController.ResetUserPassword)

}
