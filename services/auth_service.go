package services

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"microservice/app/core"
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

func (s *AuthService) syncServerClient(ctx context.Context) error {
	conn, updated, err := s.endpointService.GetConnWithStatus(ctx, "auth_service")
	if err != nil {
		return errors.Wrap(err, "cannot get endpoint client for auth_service")
	}
	if updated || s.client == nil {
		s.client = api.NewAuthServiceClient(conn)
	}
	return nil
}

func (s *AuthService) Verify(ctx context.Context, authToken string, needRole core.AccessRole) (*api.User, error) {

	// if endpoint didnt change than not updating
	err := s.syncServerClient(ctx)
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

	// Here should be user role
	if 9999 < needRole {
		s.log.Debug("user role (%d) < need role (%d): token = %s", 9999, needRole, authToken)
		return nil, nil
	}

	return verifyRes.User, nil
}
