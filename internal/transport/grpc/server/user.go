package grpcserver

import (
	"context"

	"github.com/noolingo/proto/codegen/go/apierrors"
	"github.com/noolingo/proto/codegen/go/common"
	"github.com/noolingo/proto/codegen/go/noolingo"
	"github.com/noolingo/user-service/internal/domain"
	"github.com/noolingo/user-service/internal/service"
	"github.com/sirupsen/logrus"
)

type UserServer struct {
	noolingo.UnimplementedUserServer
	logger  *logrus.Logger
	service *service.Services
}

func newUserServer(logger *logrus.Logger, service *service.Services) UserServer {
	return UserServer{logger: logger, service: service}
}

func newResponse(err error) (*common.Response, error) {
	response := &common.Response{
		Result: err == nil,
	}

	if err != nil {
		response.Error = &common.Error{
			Error: err.Error(),
		}
	}

	return response, err
}

func (u UserServer) SignUp(ctx context.Context, req *noolingo.SignUpRequest) (*common.Response, error) {
	_, err := u.service.Auth.SignUp(ctx, &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	return newResponse(err)
}

func (u UserServer) SignIn(ctx context.Context, requset *noolingo.SignInRequest) (*noolingo.SignInReply, error) {
	accessToken, refreshToken, err := u.service.Auth.SignIn(ctx, requset.Email, requset.Password)
	if err != nil {
		return nil, err
	}
	return &noolingo.SignInReply{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u UserServer) Logout(ctx context.Context, req *noolingo.LogoutRequest) (*common.Response, error) {
	_, err := Auth(ctx)
	if err != nil {
		return nil, err
	}
	panic("not implemented") // TODO: Implement
}

func (u UserServer) GetUser(ctx context.Context, req *noolingo.GetUserRequest) (*noolingo.GetUserResponse, error) {
	r, err := Auth(ctx)
	if err != nil {
		return nil, err
	}
	u.logger.Printf("userID: %v", r.UserID)
	user, err := u.service.User.GetUserByID(ctx, r.UserID)
	if err != nil {
		return nil, err
	}
	return &noolingo.GetUserResponse{Result: &noolingo.UserObject{Id: user.ID, Name: user.Name, Email: user.Email}}, nil
}

func (u UserServer) UpdateUser(ctx context.Context, req *noolingo.UpdateUserRequest) (*common.Response, error) {
	err := u.service.User.UpdateUser(ctx, &domain.User{
		Name:  req.Name,
		Email: req.Email,
	})
	return newResponse(err)
}

func (u UserServer) DeleteUser(ctx context.Context, req *noolingo.DeleteUserRequest) (*common.Response, error) {
	r, err := Auth(ctx)
	if err != nil {
		return newResponse(err)
	}
	u.logger.Printf("userID: %v Deleted", r.UserID)
	if r.UserID != req.Id {
		return newResponse(apierrors.ErrInvalidPayload)
	}
	err = u.service.Auth.DeleteUser(ctx, r.UserID)
	return newResponse(err)
}

func (u UserServer) CreateUser(ctx context.Context, req *noolingo.CreateUserRequest) (*common.Response, error) {
	_, err := u.service.Auth.SignUp(ctx, &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return newResponse(err)
}

func (u UserServer) UpdatePassword(_ context.Context, _ *noolingo.UpdatePasswordRequest) (*common.Response, error) {
	panic("not implemented") // TODO: Implement
}

func (u UserServer) Refresh(_ context.Context, _ *noolingo.RefreshRequest) (*noolingo.RefreshReply, error) {
	panic("not implemented") // TODO: Implement
}
