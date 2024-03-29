package v1

import (
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1/model/pb"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserGrpcServer struct {
	pb.UnimplementedUserUseCaseServer
	applicationConfig config.ApplicationConfig
	userUseCase       user.UserUseCase
	userCollection    *mongo.Collection
}

func NewGrpcUserServer(userUseCase user.UserUseCase, userCollection *mongo.Collection) (*UserGrpcServer, error) {
	applicationConfig := config.AppConfig
	userGrpcServer := &UserGrpcServer{
		applicationConfig: applicationConfig,
		userUseCase:       userUseCase,
		userCollection:    userCollection,
	}
	return userGrpcServer, nil
}
