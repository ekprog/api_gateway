package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// ProtoRegistry is a list of all instances with info and calling
type ProtoRegistry struct {
	instances map[string]*ProtoInstance
}

func NewProtoRegistry() *ProtoRegistry {
	return &ProtoRegistry{
		instances: make(map[string]*ProtoInstance),
	}
}

func (pr *ProtoRegistry) Init() error {
	// FOLDER auth_service
	instances, err := pr.loadInstances("./proto")
	if err != nil {
		return errors.Wrapf(err, "cannot loat services for instances")
	}

	for _, instance := range instances {
		pr.instances[instance.Name] = instance
		err := instance.loadServices()
		if err != nil {
			return errors.Wrapf(err, "cannot loat services for instance %s", instance.Name)
		}
	}
	return nil
}

func (pr *ProtoRegistry) Instance(name string) *ProtoInstance {
	return pr.instances[name]
}

func (pr *ProtoRegistry) InstanceExists(name string) bool {
	return pr.Instance(name) != nil
}

func (pr *ProtoRegistry) Instances() []string {
	var keys []string
	for k := range pr.instances {
		keys = append(keys, k)
	}
	return keys
}

// CallJson
// in - instance name
// sn - service name
// mn - method name
func (pr *ProtoRegistry) Ping(conn *grpc.ClientConn, instanceName string) ([]byte, error) {
	if !pr.InstanceExists(instanceName) {
		return nil, errors.Errorf("instance %s not found", instanceName)
	}
	instance := pr.Instance(instanceName)

	serviceName := "StatusService"
	if !instance.ServiceExist("StatusService") {
		return nil, errors.Errorf("instance`s status service not found (%s)", instanceName)
	}
	service := instance.Service(serviceName)
	return service.CallJson(conn, "Ping", []byte("{}"), nil)
}

// CallJsonWithContext
// in - instance name
// sn - service name
// mn - method name
func (pr *ProtoRegistry) CallJsonWithContext(conn *grpc.ClientConn, ctx context.Context, in, sn, mn string, jsonIn []byte, headers map[string]string) ([]byte, error) {
	if !pr.InstanceExists(in) {
		return nil, errors.Errorf("instance %s not found", in)
	}
	instance := pr.Instance(in)

	if !instance.ServiceExist(sn) {
		return nil, errors.Errorf("instance`s service not found (%s, %s)", in, sn)
	}
	service := instance.Service(sn)

	return service.CallJsonWithContext(conn, ctx, mn, jsonIn, headers)
}

func (pr *ProtoRegistry) loadInstances(protoPath string) (items []*ProtoInstance, err error) {
	files, err := os.ReadDir(protoPath)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.IsDir() {
			items = append(items, &ProtoInstance{
				Name:     f.Name(),
				Path:     path.Join(protoPath, f.Name()),
				services: make(map[string]*ProtoService),
			})
		}
	}
	return items, nil
}

type ProtoInstance struct {
	Name     string
	Path     string
	services map[string]*ProtoService
}

func (r *ProtoInstance) loadServices() error {

	// Making files list
	var files []string
	err := filepath.Walk(r.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".proto" {
			return nil
		}
		path = "./" + strings.TrimPrefix(path, r.Path+"/")
		files = append(files, path)
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "cannot list instances files")
	}

	// Load each file in registry
	for _, filePath := range files {
		err = r.parseServices(r.Path, filePath)
		if err != nil {
			return errors.Wrap(err, "cannot load file for proto registry")
		}
	}

	return nil
}

func (r *ProtoInstance) parseServices(srcDir, filename string) error {

	srcDir = strings.TrimPrefix(srcDir, "./")
	filename = strings.TrimPrefix(filename, "./")

	// File InstanceRegistry
	tmpFile := path.Join(srcDir, filename+"-tmp.pb")
	cmd := exec.Command("protoc",
		"-I",
		path.Join(srcDir),
		"--include_imports",
		"--descriptor_set_out",
		tmpFile,
		path.Join(srcDir, filename))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	marshalledDescriptorSet, err := os.ReadFile(tmpFile)
	if err != nil {
		return err
	}

	descriptorSet := descriptorpb.FileDescriptorSet{}
	err = proto.Unmarshal(marshalledDescriptorSet, &descriptorSet)
	if err != nil {
		return err
	}

	registry, err := protodesc.NewFiles(&descriptorSet)
	if err != nil {
		return err
	}

	//_, nameOnly := filepath.Split(r.filename)
	desc, err := registry.FindFileByPath(filename)
	if err != nil {
		return err
	}

	// Parse services
	services := desc.Services()
	for i := 0; i < services.Len(); i++ {
		service := services.Get(i)
		serviceObj, err := MakeProtoService(service)
		if err != nil {
			return err
		}
		r.services[string(service.Name())] = serviceObj
	}
	return nil
}

func (r *ProtoInstance) ServiceExist(name string) bool {
	return r.Service(name) != nil
}

func (r *ProtoInstance) Service(name string) *ProtoService {
	return r.services[name]
}

func (r *ProtoInstance) Services() []string {
	var keys []string
	for k := range r.services {
		keys = append(keys, k)
	}

	return keys
}

func (r *ProtoInstance) Call(conn *grpc.ClientConn, service, method string, in, out interface{}, headers map[string]string) error {
	serviceObj := r.services[service]
	if serviceObj == nil {
		return errors.New("service does not exist")
	}

	return serviceObj.Call(conn, method, in, out, headers)
}

func (r *ProtoInstance) CallJson(conn *grpc.ClientConn, service, method string, jsonIn []byte, headers map[string]string) ([]byte, error) {
	serviceObj := r.services[service]
	if serviceObj == nil {
		return nil, errors.New("service does not exist")
	}
	return serviceObj.CallJson(conn, method, jsonIn, headers)
}

func (r *ProtoInstance) CallWithContext(conn *grpc.ClientConn, ctx context.Context, service, method string, in, out interface{}, headers map[string]string) error {
	serviceObj := r.services[service]
	if serviceObj == nil {
		return errors.New("service does not exist")
	}
	return serviceObj.CallWithContext(conn, ctx, method, in, out, headers)
}

func (r *ProtoInstance) CallJsonWithContext(conn *grpc.ClientConn, ctx context.Context, service, method string, jsonIn []byte, headers map[string]string) ([]byte, error) {
	serviceObj := r.services[service]
	if serviceObj == nil {
		return nil, errors.New("service does not exist")
	}
	return serviceObj.CallJsonWithContext(conn, ctx, method, jsonIn, headers)
}

type ProtoService struct {
	service protoreflect.ServiceDescriptor
	methods map[string]*ProtoMethod
}

func MakeProtoService(descriptor protoreflect.ServiceDescriptor) (*ProtoService, error) {
	serviceObj := &ProtoService{
		service: descriptor,
		methods: make(map[string]*ProtoMethod),
	}

	// Methods
	methods := descriptor.Methods()
	for i := 0; i < methods.Len(); i++ {
		method := methods.Get(i)
		methodObj, err := MakeProtoMethod(serviceObj, method)
		if err != nil {
			return nil, err
		}
		serviceObj.methods[string(method.Name())] = methodObj
	}
	return serviceObj, nil
}

func (s *ProtoService) Name() string {
	return string(s.service.Name())
}

func (s *ProtoService) Methods() []string {
	var list []string

	for _, method := range s.methods {
		list = append(list, method.Name())
	}
	return list
}

func (s *ProtoService) Call(conn *grpc.ClientConn, method string, in, out interface{}, headers map[string]string) error {
	methodObj := s.methods[method]
	if methodObj == nil {
		return errors.New("method does not exist")
	}
	return methodObj.Call(conn, in, out, headers)
}

func (s *ProtoService) CallJson(conn *grpc.ClientConn, method string, jsonIn []byte, headers map[string]string) ([]byte, error) {
	methodObj := s.methods[method]
	if methodObj == nil {
		return nil, errors.New("method does not exist")
	}
	return methodObj.CallJson(conn, jsonIn, headers)
}

func (s *ProtoService) CallWithContext(conn *grpc.ClientConn, ctx context.Context, method string, in, out interface{}, headers map[string]string) error {
	methodObj := s.methods[method]
	if methodObj == nil {
		return errors.New("method does not exist")
	}
	return methodObj.CallWithContext(ctx, conn, in, out, headers)
}

func (s *ProtoService) CallJsonWithContext(conn *grpc.ClientConn, ctx context.Context, method string, jsonIn []byte, headers map[string]string) ([]byte, error) {
	methodObj := s.methods[method]
	if methodObj == nil {
		return nil, errors.New(fmt.Sprintf("method %s does not exist", method))
	}
	return methodObj.CallJsonWithContext(ctx, conn, jsonIn, headers)
}

type ProtoMethod struct {
	parent   *ProtoService
	method   protoreflect.MethodDescriptor
	request  protoreflect.MessageDescriptor
	response protoreflect.MessageDescriptor
}

func MakeProtoMethod(parent *ProtoService, descriptor protoreflect.MethodDescriptor) (*ProtoMethod, error) {

	methodObj := &ProtoMethod{
		parent:   parent,
		method:   descriptor,
		request:  descriptor.Input(),
		response: descriptor.Output(),
	}

	return methodObj, nil
}

func (m *ProtoMethod) Name() string {
	return string(m.method.Name())
}

func (m *ProtoMethod) Call(conn *grpc.ClientConn, in interface{}, out interface{}, headers map[string]string) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return m.CallWithContext(ctx, conn, in, out, headers)
}

func (m *ProtoMethod) CallJson(conn *grpc.ClientConn, jsonInput []byte, headers map[string]string) ([]byte, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return m.CallJsonWithContext(ctx, conn, jsonInput, headers)
}

func (m *ProtoMethod) CallWithContext(ctx context.Context, conn *grpc.ClientConn, in interface{}, out interface{}, headers map[string]string) error {

	jsonInput, err := json.Marshal(in)
	if err != nil {
		return err
	}

	jsonOutput, err := m.CallJsonWithContext(ctx, conn, jsonInput, headers)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonOutput, out)
	if err != nil {
		return err
	}

	return nil
}

func (m *ProtoMethod) CallJsonWithContext(ctx context.Context, conn *grpc.ClientConn, jsonInput []byte, headers map[string]string) ([]byte, error) {

	requestObj := dynamicpb.NewMessage(m.request)
	responseObj := dynamicpb.NewMessage(m.response)

	err := protojson.Unmarshal(jsonInput, requestObj)
	if err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal json to proto descriptor")
	}

	if headers != nil && len(headers) != 0 {
		headersMeta := metadata.New(headers)
		ctx = metadata.NewOutgoingContext(ctx, headersMeta)
	}

	caller := "/" + string(m.parent.service.FullName()) + "/" + string(m.method.Name())
	err = conn.Invoke(ctx, caller, requestObj, responseObj)
	if err != nil {
		return nil, errors.Wrap(err, "error while invoke proto method")
	}

	jsonOutput, err := protojson.Marshal(responseObj)
	if err != nil {
		return nil, errors.Wrap(err, "cannot marshal proto response to json")
	}
	return jsonOutput, nil
}
