package dependency

import (
	"context"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	factory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
)

// NewApplication initializes the application by setting up the container,
// injecting dependencies, and configuring the server.
func NewApplication(ctx context.Context) model.Container {
	config := factory.NewConfig(constants.Config)
	logger := factory.NewLogger(ctx, config)

	// Create repository factory and repositories
	repositoryFactory := factory.NewRepositoryFactory(ctx, config, logger)
	createRepository := repositoryFactory.CreateRepository(ctx)
	userRepository := repositoryFactory.NewRepository(createRepository, (*interfaces.UserRepository)(nil))
	postRepository := repositoryFactory.NewRepository(createRepository, (*interfaces.PostRepository)(nil))

	// Create use case factory and use cases.
	usecaseFactory := factory.NewUseCaseFactory(ctx, config, logger, repositoryFactory)
	userUseCase := usecaseFactory.NewUseCase(userRepository).(interfaces.UserUseCase)
	postUseCase := usecaseFactory.NewUseCase(postRepository)

	// Create delivery factory and controllers.
	deliveryFactory := factory.NewDeliveryFactory(ctx, config, logger, repositoryFactory)
	userController := deliveryFactory.NewController(userUseCase, nil)
	postController := deliveryFactory.NewController(userUseCase, postUseCase)

	// Set up server routers with the user use case and controllers.
	serverRouters := interfaces.NewServerRouters(
		userUseCase,
		deliveryFactory.NewRouter(userController).(interfaces.UserRouter),
		deliveryFactory.NewRouter(postController).(interfaces.PostRouter),
		// Add other routers as needed.
	)

	deliveryFactory.CreateDelivery(serverRouters)
	container := model.NewContainer(logger, repositoryFactory, deliveryFactory)
	return container
}
