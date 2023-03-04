package services

import (
	"api_gateway/app"
	"api_gateway/domain"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
)

type RouterService struct {
	log         app.Logger
	g           *gin.RouterGroup
	services    map[string]*app.Service
	authService *app.Service
}

func NewRouterService(log app.Logger) *RouterService {
	return &RouterService{
		log:      log,
		services: map[string]*app.Service{},
	}
}

func (s *RouterService) SetRouter(g *gin.RouterGroup) {
	s.g = g
}

func (s *RouterService) SetAuthService(serviceName string) error {
	service, ok := s.services[serviceName]
	if !ok {
		return errors.New("cannot find auth service with name " + serviceName)
	}
	s.authService = service
	return nil
}

func (s *RouterService) MakeService(sObj domain.Service) error {
	registry := app.NewRegistry()
	err := registry.Reload(sObj.ProtoDir, sObj.ProtoFilename)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("cannot create %s registry", sObj.ProtoServiceName))
	}

	service := registry.Service(sObj.ProtoServiceName)
	if service == nil {
		return errors.New(fmt.Sprintf("proto service %s was not found", sObj.ProtoServiceName))
	}

	client, err := s.makeClient(sObj.HttpAddress, "")
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("cannot create http client for %s", sObj.ProtoServiceName))
	}
	service.SetClient(client)

	s.services[sObj.ProtoServiceName] = service

	return nil
}

func (s *RouterService) Handle(r domain.Route) error {

	service, ok := s.services[r.ProtoService]
	if !ok {
		return errors.New("service " + r.ProtoService + " was not found")
	}

	s.g.Handle(r.HttpMethod, r.HttpAddress, func(c *gin.Context) {
		c.Header("content-type", "application/json")

		// DATA
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			s.log.Error(fmt.Sprintf("error while parsing data from %s: %s", r.HttpAddress, err.Error()))
			c.JSON(500, domain.ServerError())
			return
		}

		// HEADERS
		headers := map[string]string{}

		// need auth
		if r.AccessRole > domain.RoleGuest {
			authTokens := c.Request.Header["Authorization"]
			if authTokens == nil || len(authTokens) <= 0 {
				c.JSON(500, domain.UnauthorizedError())
				return
			}

			authToken := authTokens[0]
			if len(authToken) == 0 {
				c.JSON(500, domain.UnauthorizedError())
				return
			}
			s.log.Debug("Authorization access with token: %s", authToken)

			// verifying role
			verifyReq := domain.VerifyRequest{
				AccessToken: authToken,
			}
			var verifyRes domain.VerifyResponse
			err := s.authService.CallWithContext(c, "Verify", verifyReq, &verifyRes, nil)
			if err != nil {
				s.log.Error(fmt.Sprintf("error while verifying access %s: %s", r.HttpAddress, err.Error()))
				c.JSON(500, domain.ServerError())
				return
			}

			if verifyRes.Status.Code != "success" || verifyRes.User == nil {
				s.log.Debug(fmt.Sprintf("incorrect token (%s) from %s", authToken, r.HttpAddress))
				c.JSON(500, domain.UnauthorizedError())
				return
			}

			headers["user_id"] = verifyRes.User.Id
		}

		// Continue call
		res, err := service.CallJsonWithContext(c, r.ProtoMethod, jsonData, headers)
		if err != nil {
			s.log.Error(fmt.Sprintf("error while redirect from %s: %s", r.HttpAddress, err.Error()))
			c.JSON(500, domain.ServerError())
			return
		}
		s.log.Debug(fmt.Sprintf("response from %s: %s", r.HttpAddress, (string)(res)))

		// Response status code
		type Status struct {
			Status domain.Status `json:"status"`
		}
		stat := &Status{}
		err = json.Unmarshal(res, stat)
		if err != nil {
			s.log.Error(fmt.Sprintf("error while unmarshal for status code from %s: %s", r.HttpAddress, err.Error()))
			c.JSON(500, domain.ServerError())
			return
		}
		if stat.Status.Code != "success" {
			c.Status(500)
		} else {
			c.Status(200)
		}

		c.Writer.Write(res)
	})
	return nil
}

func (s *RouterService) makeClient(addr string, certRoot string) (*grpc.ClientConn, error) {

	tslEnable := certRoot != ""
	if tslEnable {

		crt := certRoot + "/ca.cert"
		key := certRoot + "/ca.key"
		caN := certRoot + "/ca.cert"

		// Load the client certificates from disk
		certificate, err := tls.LoadX509KeyPair(crt, key)
		if err != nil {
			return nil, errors.Wrap(err, "could not load client key pair")
		}

		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(caN)
		if err != nil {
			return nil, errors.Wrap(err, "could not read ca certificate")
		}

		// Append the certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return nil, errors.Wrap(err, "failed to append ca certs")
		}

		creds := credentials.NewTLS(&tls.Config{
			ServerName:   addr, // NOTE: this is required!
			Certificates: []tls.Certificate{certificate},
			RootCAs:      certPool,
		})

		// Create a connection with the TLS credentials
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
		if err != nil {
			return nil, errors.Wrap(err, "could not dial")
		}
		return conn, nil
	} else {
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, errors.Wrap(err, "did not connect")
		}
		return conn, nil
	}
}
