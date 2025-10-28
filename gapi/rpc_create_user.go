package gapi

import (
	"context"

	db "github.com/joekings2k/gobank/db/sqlc"
	"github.com/joekings2k/gobank/pb"
	"github.com/joekings2k/gobank/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword ,err  := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to hash password: %v", err)
		
	} 
	arg := db.CreateUserParams{
		Username: req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName: req.GetFullName(),
		Email: req.GetEmail(),
	}

	user,err  := server.store.CreateUser(ctx,arg)
	if err != nil{
		if pqErr, ok :=err.(*pq.Error);ok {
			switch pqErr.Code.Name(){
			case  "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %v", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "Failed to create user: %v", err)
	}
	response := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return  response, nil
}