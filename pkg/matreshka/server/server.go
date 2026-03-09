package server

import (
	"fmt"
	"strings"

	"go.redsock.ru/evon"
	"go.redsock.ru/rerrors"
	"gopkg.in/yaml.v3"
)

const (
	EnvServerName = "server_name"

	grpcPath       = "/{GRPC}"
	fileServerPath = "/{FS}"

	fieldName = "name"
	portField = "port"
)

type Server struct {
	Name string `yaml:"name,omitempty"`
	Port string `yaml:"-"`

	GRPC map[string]*GRPC
	FS   map[string]*FS

	// HTTP - default server handler
	HTTP map[string]*HTTP
}

func (s *Server) UnmarshalYAML(unmarshal func(interface{}) error) error {
	s.GRPC = make(map[string]*GRPC)
	s.HTTP = make(map[string]*HTTP)
	s.FS = make(map[string]*FS)

	m := map[string]yaml.Node{}
	err := unmarshal(m)
	if err != nil {
		return rerrors.Wrap(err, "error unmarshalling YAML")
	}

	for key, value := range m {
		switch key {
		case fieldName:
			s.Name = value.Value
		default:
			// TODO пока что думаем что конфиг будет по корню,
			// TODO но потом надо задуматься над не корневыми конфигами
			vPtr := s.getPtrToImplByWebPath(key)
			if vPtr == nil {
				continue
			}
			err = value.Decode(vPtr)
			if err != nil {
				return rerrors.Wrapf(err, "error decoding server description of type %T", vPtr)
			}
		}

	}

	return nil
}

func (s *Server) MarshalYAML() (any, error) {
	m := map[string]any{}

	if s.Name != "" {
		m[fieldName] = s.Name
	}

	for path, srv := range s.GRPC {
		m[path] = srv
	}

	for path, srv := range s.FS {
		m[path] = srv
	}

	for path, srv := range s.HTTP {
		m[path] = srv
	}

	return m, nil
}

func (s *Server) MarshalEnv(name string) ([]*evon.Node, error) {
	root := &evon.Node{
		Name: name,
	}
	if s.Port == "" {
		return nil, rerrors.New("server must have port")
	}
	root.InnerNodes = append(root.InnerNodes, &evon.Node{
		Name:  name + "_PORT",
		Value: s.Port,
	})

	if name != "" {
		name += evon.ObjectSplitter
	}

	for path, srv := range s.GRPC {
		subPrefix := name + path

		nodes, err := evon.MarshalEnvWithPrefix(subPrefix, srv)
		if err != nil {
			return nil, rerrors.Wrap(err, "error marshalling grpc server desc to env")
		}
		root.InnerNodes = append(root.InnerNodes, nodes)
	}

	for path, srv := range s.FS {
		subPrefix := name + path

		nodes, err := evon.MarshalEnvWithPrefix(subPrefix, srv)
		if err != nil {
			return nil, rerrors.Wrap(err, "error marshalling grpc server desc to env")
		}
		root.InnerNodes = append(root.InnerNodes, nodes)
	}

	for path, srv := range s.HTTP {
		subPrefix := name + path

		nodes, err := evon.MarshalEnvWithPrefix(subPrefix, srv)
		if err != nil {
			return nil, rerrors.Wrap(err, "error marshalling grpc server desc to env")
		}
		root.InnerNodes = append(root.InnerNodes, nodes)
	}

	return []*evon.Node{root}, nil
}

func (s *Server) UnmarshalEnv(v *evon.Node) error {
	s.GRPC = make(map[string]*GRPC)
	s.HTTP = make(map[string]*HTTP)
	s.FS = make(map[string]*FS)

	for _, node := range v.InnerNodes {

		suffix := node.Name[len(v.Name)+1:]
		switch strings.ToLower(suffix) {
		case portField:
			s.Port = fmt.Sprint(node.Value)
			continue
		}
		webPath := node.Name[len(v.Name)+1:]
		ptr := s.getPtrToImplByWebPath(webPath)
		if ptr == nil {
			continue
		}
		err := evon.NodeToStruct(node.Name, node, ptr)
		if err != nil {
			return rerrors.Wrap(err, "error unmarshalling value")
		}

	}

	return nil
}

func (s *Server) getPtrToImplByWebPath(webPath string) any {
	var vPtr any
	switch webPath {
	case grpcPath:
		s.GRPC[webPath] = &GRPC{}
		vPtr = s.GRPC[webPath]
	case fileServerPath:
		s.FS[webPath] = &FS{}
		vPtr = s.FS[webPath]
	default:
	}

	return vPtr
}
