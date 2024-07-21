package factory

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/data/repository"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logger "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logger"
)

const (
	unsupportedDatabase = "Unsupported database type: %s"
)

func NewRepositoryFactory(ctx context.Context) applicationModel.Repository {
	coreConfig := config.GetCoreConfig()

	switch coreConfig.Database {
	case constants.MongoDB:
		return repository.NewMongoDBRepository()
	// Add other repository options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedDatabase, coreConfig.Database)
		internalError := domainError.NewInternalError(location+"NewRepositoryFactory", notification)
		logger.Logger(internalError)
		applicationModel.GracefulShutdown(ctx, nil, nil)
		panic(internalError)
	}
}
