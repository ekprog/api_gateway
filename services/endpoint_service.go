package services

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"microservice/app/core"
	"microservice/domain"
)

// EndpointConnectionService get endpoint for service instance
type EndpointConnectionService struct {
	log           core.Logger
	instancesRepo domain.InstancesRepository

	clients   map[string]*grpc.ClientConn
	endpoints map[string]string
}

func NewEndpointConnectionService(log core.Logger, instancesRepo domain.InstancesRepository) *EndpointConnectionService {
	return &EndpointConnectionService{
		clients:       make(map[string]*grpc.ClientConn),
		endpoints:     make(map[string]string),
		log:           log,
		instancesRepo: instancesRepo}
}

func (s *EndpointConnectionService) GetConn(instanceName string) (*grpc.ClientConn, error) {
	conn, _, err := s.GetConnWithStatus(instanceName)
	return conn, err
}

// GetConnWithStatus also return bool status if new value was fetched
func (s *EndpointConnectionService) GetConnWithStatus(instanceName string) (*grpc.ClientConn, bool, error) {
	instance, err := s.instancesRepo.GetByFolder(instanceName)
	if err != nil {
		return nil, false, errors.Wrapf(err, "cannot get %s instance", instanceName)
	}
	if instance == nil {
		return nil, false, errors.Errorf("%s not found", instanceName)
	}

	// If endpoint is the same then not updating
	if c, ok1 := s.clients[instanceName]; ok1 {
		if v, ok2 := s.endpoints[instanceName]; ok2 && v == instance.Endpoint {
			return c, false, nil
		}
	}

	// creating new connection
	conn, err := grpc.Dial(instance.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, false, errors.Wrapf(err, "cannot make Dial for AuthService with address %s", s.endpoints[instanceName])
	}
	s.endpoints[instanceName] = instance.Endpoint
	s.clients[instanceName] = conn

	return conn, true, nil
}
