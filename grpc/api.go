package grpc

import (
	"auth_service/grpc/pb"
	"auth_service/storage/models"
	"context"
	"log"
)

type Server struct{}

func (s *Server) GetUser(ctx context.Context, requestData *pb.UserRequest) (*pb.UserResponse, error) {

	resUser, err := models.GetUserByUsername(requestData.Username)

	if err != nil {
		log.Fatalf("falied responding to grpc call: %v", err)
	}

	user := &pb.User{UserId: resUser.Id.String(), Username: resUser.Username, Email: resUser.Email}
	return &pb.UserResponse{User: user}, nil

}
