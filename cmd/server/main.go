package main

// Require the packages.
import (
	"context"
	"log"
	"net"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/golang-mongo-grpc/config"

	postProtobufV1 "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/grpc/v1/model/pb"
	userProtobufV1 "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1/model/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	postPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postGrpcPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/grpc/v1"
	postHttpGinPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/http/gin"
	postUseCasePackage "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/usecase"

	userPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	repository "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/data/repository"

	userGrpcPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1"
	userHttpGinPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin"
	userUseCasePackage "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"

	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"

	"go.mongodb.org/mongo-driver/mongo"
)

// Create required variables that we'll re-assign later.
var (
	server      *gin.Engine
	ctx         context.Context
	mongoClient *mongo.Client

	userCollection *mongo.Collection
	userRepository userPackage.UserRepository
	userUseCase    userPackage.UserUseCase
	userController userHttpGinPackage.UserController
	userRouter     userHttpGinPackage.UserRouter

	postRepository postPackage.PostRepository
	postUseCase    postPackage.PostUseCase
	postController postHttpGinPackage.PostHandler
	postRouter     postHttpGinPackage.PostRouter
)

const (
	location = "cmd.server.init."
)

// Init function that will run before the "main" function.
func init() {

	// Load the .env variables.
	loadConfig, loadConfigError := config.LoadConfig(".")
	if loadConfigError != nil {
		loadConfigInternalError := domainError.NewInternalError(location+"init.LoadConfig", loadConfigError.Error())
		logging.Logger(loadConfigInternalError)
	}

	// Create a context.
	ctx = context.TODO()

	repositoryFactory := repository.InjectRepository(loadConfig)

	// Create a DB database instance using the factory.
	db := repositoryFactory.NewRepository(ctx)
	userRepository = repositoryFactory.NewUserRepository(db)
	postRepository = repositoryFactory.NewPostRepository(db)

	// Use Cases.
	userUseCase = userUseCasePackage.NewUserUseCase(userRepository)
	postUseCase = postUseCasePackage.NewPostUseCase(postRepository)

	// Handlers
	userController = userHttpGinPackage.NewUserController(userUseCase)
	postController = postHttpGinPackage.NewPostHandler(postUseCase)

	// Routers.
	userRouter = userHttpGinPackage.NewUserRouter(userController)
	postRouter = postHttpGinPackage.NewPostRouter(postController)

	// Create the Gin Engine instance.
	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoClient.Disconnect(ctx)

	startGinServer(config)
	// startGrpcServer(config)
}

func startGrpcServer(config config.Config) {
	userGrpcServer, err := userGrpcPackage.NewGrpcUserServer(config, userUseCase, userCollection)

	if err != nil {
		log.Fatal("cannot createt gRPC User Server: ", err)
	}

	postGrpcServer, err := postGrpcPackage.NewGrpcPostServer(postUseCase)

	if err != nil {
		log.Fatal("cannot create gRPC Post Server: ", err)
	}

	grpcServer := grpc.NewServer()

	// Register User gRPC server.
	userProtobufV1.RegisterUserUseCaseServer(grpcServer, userGrpcServer)

	// Register Post gRPC server.
	postProtobufV1.RegisterPostUseCaseServer(grpcServer, postGrpcServer)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)

	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}

	log.Printf("start grpc server on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}

}

func startGinServer(config config.Config) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8080", "http://localhost:3000"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")

	userRouter.UserRouter(router, userUseCase)
	postRouter.PostRouter(router, userUseCase)

	log.Fatal(server.Run(":" + config.Port))
}
