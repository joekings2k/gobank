package gapi

import (
	"context"
	"database/sql"

	db "github.com/joekings2k/gobank/db/sqlc"
	"github.com/joekings2k/gobank/pb"
	"github.com/joekings2k/gobank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user,err := server.store.GetUser(ctx,req.GetUsername())
	if err != nil {
		if err ==sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "Failed to get user: %v", err)
	}
	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "incorrect password: %v", err)
	}
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create access token: %v", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create refresh token: %v", err)
	}
	mtdt := server.extractMetadata(ctx)

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     sql.NullString{String: user.Username, Valid: true},
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create session: %v", err)
	}

	response := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}
	return response, nil
}