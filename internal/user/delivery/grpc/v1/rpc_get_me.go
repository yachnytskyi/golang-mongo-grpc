package v1

import (
	"context"

	pb "github.com/yachnytskyi/golang-mongo-grpc/pkg/proto-generated"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (userServer *UserServer) GetMe(ctx context.Context, request *pb.GetMeRequest) (*pb.UserResponse, error) {
	userID := request.GetId()
	user, err := userServer.userService.GetUserById(ctx, userID)

	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, err.Error())
	}

	response := &pb.UserResponse{
		User: &pb.User{
			Id:        user.UserID.Hex(),
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}

	return response, nil
}
