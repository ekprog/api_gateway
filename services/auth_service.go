package services

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"microservice/app/core"
	"microservice/domain"
	"microservice/pkg/auth_service/api"
)

// AuthService calls remote proto
type AuthService struct {
	log             core.Logger
	endpointService *EndpointConnectionService
	client          api.AuthServiceClient
}

func NewAuthService(log core.Logger,
	endpointService *EndpointConnectionService) *AuthService {
	return &AuthService{
		log:             log,
		endpointService: endpointService,
	}
}

func (s *AuthService) syncServerClient() error {
	conn, updated, err := s.endpointService.GetConnWithStatus("auth_service")
	if err != nil {
		return errors.Wrap(err, "cannot get endpoint client for auth_service")
	}
	if updated || s.client == nil {
		s.client = api.NewAuthServiceClient(conn)
	}
	return nil
}

func (s *AuthService) VerifyRequest(ctx *gin.Context, needRole domain.AccessRole) (*api.User, error) {

	// Extract token
	authTokens := ctx.Request.Header["Authorization"]
	if authTokens == nil || len(authTokens) <= 0 {
		return nil, nil
	}
	authToken := authTokens[0]
	s.log.Debug("Authorization access with token: %s", authToken)

	// Verify
	user, err := s.Verify(ctx, authToken)
	if err != nil {
		return nil, errors.Wrapf(err, "error while verify (auth) user by token %s", authToken)
	}
	if user == nil {
		s.log.Debug("auth_service return not verified %s", authToken)
		return nil, nil
	}

	// Here should be user role
	if 9999 < needRole {
		s.log.Debug("user role (%d) < need role (%d): token = %s", 9999, needRole, authToken)
		return nil, nil
	}

	return user, nil
}

func (s *AuthService) Verify(ctx context.Context, authToken string) (*api.User, error) {

	// if endpoint didnt change than not updating
	err := s.syncServerClient()
	if err != nil {
		return nil, errors.Wrap(err, "cannot make client for auth_service isntance")
	}

	// verifying role
	verifyReq := &api.VerifyRequest{
		AccessToken: authToken,
	}

	verifyRes, err := s.client.Verify(ctx, verifyReq)
	if err != nil {
		return nil, errors.Wrapf(err, "error while verifying access %s", authToken)
	}

	if verifyRes.Status.Code != "success" || verifyRes.User == nil {
		s.log.Debug(fmt.Sprintf("incorrect token (%s) for auth_service", authToken))
		return nil, nil
	}

	return verifyRes.User, nil
}
