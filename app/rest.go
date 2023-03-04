package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"
)

var (
	restServer *gin.Engine
)

func InitRest() error {
	restServer = gin.Default()

	// CORS
	restServer.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{},
		AllowCredentials: true,
		AllowOriginFunc:  nil,
		MaxAge:           12 * time.Hour,
	}))

	return nil
}

func RunRestServer() {
	restServer.Run(os.Getenv("REST_ADDRESS"))
}

type RestDeliveryService interface {
	Route(r *gin.RouterGroup) error
}

func InitRestDelivery(path string, d RestDeliveryService) error {
	g := restServer.Group(path)
	err := d.Route(g)
	if err != nil {
		return err
	}
	return nil
}

func RegisterDelivery() {

}

func createProtoRegistry(srcDir string, filename string) (*protoregistry.Files, error) {
	// Create descriptors using the protoc binary.
	// Imported dependencies are included so that the descriptors are self-contained.
	srcDir = path.Join("./proto", srcDir)
	tmpFile := path.Join(srcDir, filename+"-tmp.pb")
	cmd := exec.Command("protoc",
		"--include_imports",
		"--descriptor_set_out="+tmpFile,
		"-I "+srcDir,
		path.Join(srcDir, filename))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile)

	marshalledDescriptorSet, err := ioutil.ReadFile(tmpFile)
	if err != nil {
		return nil, err
	}
	descriptorSet := descriptorpb.FileDescriptorSet{}
	err = proto.Unmarshal(marshalledDescriptorSet, &descriptorSet)
	if err != nil {
		return nil, err
	}

	files, err := protodesc.NewFiles(&descriptorSet)
	if err != nil {
		return nil, err
	}

	return files, nil
}

type ReflectionService struct {
	conn       *grpc.ClientConn
	request    protoreflect.MessageDescriptor
	response   protoreflect.MessageDescriptor
	rpcService string
	rpcMethod  string
}

func RegisterReflectionService(serviceName, protoPath, addr, rpcService, rpcMethod, restMethod string) (*ReflectionService, error) {

	// Creating registry
	registry, err := createProtoRegistry(serviceName, protoPath)
	if err != nil {
		return nil, err
	}

	// Parse file
	desc, err := registry.FindFileByPath(protoPath)
	if err != nil {
		return nil, err
	}

	services := desc.Services()
	//for i := 0; i < services.Len(); i++ {
	//	service := services.Get(i)
	//	log.Info(service.Name())
	//	for i := 0; i < services.Len(); i++ {
	//
	//	}
	//}
	//log.Info("")

	// ByName
	service := services.ByName(protoreflect.Name(rpcService))
	if service == nil {
		return nil, err
	}

	methods := service.Methods()
	method := methods.ByName(protoreflect.Name(rpcMethod))
	if method == nil {
		return nil, err
	}

	request := method.Input()
	response := method.Input()

	conn, err := grpc.Dial("localhost:8086", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	s := &ReflectionService{
		conn:       conn,
		request:    request,
		response:   response,
		rpcService: rpcService,
		rpcMethod:  rpcMethod,
	}

	restServer.Handle(restMethod, addr, s.httpHandler)

	return s, nil
}

func (s *ReflectionService) httpHandler(c *gin.Context) {

	requestObj := dynamicpb.NewMessage(s.request)
	responseObj := dynamicpb.NewMessage(s.response)

	// Parse
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, "validation error")
		return
	}
	err = protojson.Unmarshal(jsonData, requestObj)
	if err != nil {
		c.JSON(500, "validation error")
		return
	}

	// Invoke

	err = s.conn.Invoke(c, "/pb."+s.rpcService+"/"+s.rpcMethod, requestObj, responseObj)
	if err != nil {
		c.JSON(500, "validation error")
		return
	}
	marshal, err := protojson.Marshal(responseObj)
	if err != nil {
		c.JSON(500, "validation error")
		return
	}

	c.Status(200)
	c.Header("content-type", "application/json")
	c.Writer.Write(marshal)
}
