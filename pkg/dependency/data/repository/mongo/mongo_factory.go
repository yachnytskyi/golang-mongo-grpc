package mongo

import (
	"context"

	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postRepository "github.com/yachnytskyi/golang-mongo-grpc/internal/post/data/repository/mongo"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userRepository "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/mongo"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
)

const (
	successfully_connected = "Database is successfully connected..."
	successfully_closed    = "Database connection has been successfully closed..."

	location            = "pkg.dependency.data.repository.mongo.NewRepository."
	unsupportedDatabase = "unsupported database type: %s"
)

type MongoDBFactory struct {
	MongoDBConfig config.MongoDBConfig
	MongoClient   *mongo.Client
}

func (mongoDBFactory *MongoDBFactory) NewRepository(ctx context.Context) interface{} {
	var connectError error
	mongoConnection := options.Client().ApplyURI(mongoDBFactory.MongoDBConfig.URI)
	mongoDBFactory.MongoClient, connectError = mongo.Connect(ctx, mongoConnection)
	db := mongoDBFactory.MongoClient.Database(mongoDBFactory.MongoDBConfig.Name)
	if validator.IsErrorNotNil(connectError) {
		logging.Logger(domainError.NewInternalError(location+"mongoClient.Database", connectError.Error()))
		return nil
	}
	connectError = mongoDBFactory.MongoClient.Ping(ctx, readpref.Primary())
	if validator.IsErrorNotNil(connectError) {
		logging.Logger(domainError.NewInternalError(location+"mongoClient.Ping", connectError.Error()))
		return nil
	}
	logging.Logger(successfully_connected)
	return db
}

func (mongoDBFactory *MongoDBFactory) CloseRepository() {
	if validator.IsValueNotNil(mongoDBFactory.MongoClient) {
		mongoDBFactory.MongoClient.Disconnect(context.Background())
		logging.Logger(successfully_closed)
	}
}

func (mongoDBFactory *MongoDBFactory) NewUserRepository(db interface{}) user.UserRepository {
	mongoDB := db.(*mongo.Database)
	return userRepository.NewUserRepository(mongoDB)
}

func (mongoDBFactory *MongoDBFactory) NewPostRepository(db interface{}) post.PostRepository {
	mongoDB := db.(*mongo.Database)
	return postRepository.NewPostRepository(mongoDB)
}
