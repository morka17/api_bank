package gapi

import (
	"github.com/morka17/shiny_bank/v1/pb"
	db "github.com/morka17/shiny_bank/v1/src/db/sqlc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func UserToProto(user db.User) *pb.User{
	return &pb.User{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangeAt),
		CreatedAt: timestamppb.New(user.CreatedAt),		
	}
}