package v1

import (
	"github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/grpc/v1/model/pb"
)

type PostGrpcServer struct {
	pb.UnimplementedPostUseCaseServer
	postUseCase post.UseCase
}

func NewGrpcPostServer(postUseCase post.UseCase) (*PostGrpcServer, error) {
	postGrpcServer := &PostGrpcServer{
		postUseCase: postUseCase,
	}

	return postGrpcServer, nil
}
