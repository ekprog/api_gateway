package domain

import "github.com/gin-gonic/gin"

type Service struct {
	ProtoServiceName string
	ProtoDir         string
	ProtoFilename    string
	HttpAddress      string
}

type Route struct {
	HttpMethod   string
	HttpAddress  string
	ProtoService string
	ProtoMethod  string
	AccessRole   AccessRole
}

type RouterService interface {
	SetRouter(group *gin.RouterGroup)
	MakeService(Service) error
	SetAuthService(string) error
	Handle(Route) error
}
