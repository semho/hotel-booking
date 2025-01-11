package mapper

import (
	"github.com/semho/hotel-booking/auth-service/internal/domain/model"
	pb "github.com/semho/hotel-booking/pkg/proto/auth_v1/auth"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProtoAuthResponse(resp *model.AuthResponse) *pb.AuthResponse {
	return &pb.AuthResponse{
		AccessToken:           resp.AccessToken,
		RefreshToken:          resp.RefreshToken,
		AccessTokenExpiresAt:  timestamppb.New(resp.AccessTokenExpiresAt),
		RefreshTokenExpiresAt: timestamppb.New(resp.RefreshTokenExpiresAt),
		User:                  ToProtoUser(resp.User),
	}
}

func ToProtoUser(user *model.User) *pb.UserInfo {
	return &pb.UserInfo{
		Id:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Role:      ToProtoUserRole(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func ToProtoUserRole(role model.UserRole) pb.UserRole {
	switch role {
	case model.UserRoleAdmin:
		return pb.UserRole_USER_ROLE_ADMIN
	case model.UserRoleUser:
		return pb.UserRole_USER_ROLE_USER
	default:
		return pb.UserRole_USER_ROLE_UNSPECIFIED
	}
}

func ToProtoUserInfo(user *model.User) *pb.UserInfo {
	return &pb.UserInfo{
		Id:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Role:      ToProtoUserRole(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
