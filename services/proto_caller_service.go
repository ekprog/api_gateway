package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"microservice/app"
	"microservice/app/core"
)

type ProtoCall struct {
	Instance string
	Service  string
	Method   string
	Data     []byte
	Headers  map[string]string
}

// ProtoCallerService делает вызов к микросервисам и парсит запрос
type ProtoCallerService struct {
	log core.Logger

	// Данные об ip сервисов
	endpointService *EndpointConnectionService

	// Описание api сервисов
	protoRegistry *app.ProtoRegistry
}

func NewProtoCallerService(log core.Logger,
	protoRegistry *app.ProtoRegistry,
	endpointService *EndpointConnectionService) *ProtoCallerService {
	return &ProtoCallerService{
		log:             log,
		protoRegistry:   protoRegistry,
		endpointService: endpointService,
	}
}

func (s *ProtoCallerService) Call(ctx context.Context, call ProtoCall) ([]byte, error) {

	// Find in proto registry
	protoInstance := s.protoRegistry.Instance(call.Instance)
	if protoInstance == nil {
		return nil, errors.Errorf("cannot find instance %s", call.Instance)
	}

	service := protoInstance.Service(call.Service)
	if service == nil {
		return nil, errors.Errorf("cannot find service %s in instance %s", call.Service, call.Instance)
	}

	// Get connection
	conn, err := s.endpointService.GetConn(ctx, call.Instance)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get endpoint client for %s", call.Instance)
	}

	// Call
	if call.Data == nil {
		call.Data = []byte("{}")
	}
	res, err := service.CallJsonWithContext(conn, ctx, call.Method, call.Data, call.Headers)
	if err != nil {
		return nil, errors.Wrapf(err, "in instance error (%s.%s.%s)", call.Instance, call.Service, call.Method)
	}
	s.log.Debug(fmt.Sprintf("response from %s: %s", call.Instance, (string)(res)))

	return res, nil
}

func (s *ProtoCallerService) CallAndParse(ctx context.Context, call ProtoCall, out interface{}) ([]byte, error) {
	response, err := s.Call(ctx, call)
	if err != nil {
		return nil, err
	}

	// Parsing status code (REQUIRED FOR EACH API SERVICE)
	err = json.Unmarshal(response, out)
	if err != nil {
		return nil, errors.Wrapf(err, "error while unmarshal for status code from %s.%s.%s", call.Instance, call.Service, call.Method)
	}
	return response, nil
}
