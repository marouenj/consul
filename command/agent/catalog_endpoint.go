package agent

import (
	"fmt"
	"github.com/marouenj/consul/consul/structs"
	"net/http"
	"strings"
)

func (s *HTTPServer) CatalogRegister(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	var args structs.RegisterRequest
	if err := decodeBody(req, &args, nil); err != nil {
		resp.WriteHeader(400)
		resp.Write([]byte(fmt.Sprintf("Request decode failed: %v", err)))
		return nil, nil
	}

	// Setup the default DC if not provided
	if args.Datacenter == "" {
		args.Datacenter = s.agent.config.Datacenter
	}

	// Forward to the servers
	var out struct{}
	if err := s.agent.RPC("Catalog.Register", &args, &out); err != nil {
		return nil, err
	}
	return true, nil
}

func (s *HTTPServer) CatalogDeregister(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	var args structs.DeregisterRequest
	if err := decodeBody(req, &args, nil); err != nil {
		resp.WriteHeader(400)
		resp.Write([]byte(fmt.Sprintf("Request decode failed: %v", err)))
		return nil, nil
	}

	// Setup the default DC if not provided
	if args.Datacenter == "" {
		args.Datacenter = s.agent.config.Datacenter
	}

	// Forward to the servers
	var out struct{}
	if err := s.agent.RPC("Catalog.Deregister", &args, &out); err != nil {
		return nil, err
	}
	return true, nil
}

func (s *HTTPServer) CatalogDatacenters(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	var out []string
	if err := s.agent.RPC("Catalog.ListDatacenters", struct{}{}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *HTTPServer) CatalogNodes(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	// Setup the request
	args := structs.DCSpecificRequest{}
	if done := s.parse(resp, req, &args.Datacenter, &args.QueryOptions); done {
		return nil, nil
	}

	var out structs.IndexedNodes
	defer setMeta(resp, &out.QueryMeta)
	if err := s.agent.RPC("Catalog.ListNodes", &args, &out); err != nil {
		return nil, err
	}
	return out.Nodes, nil
}

func (s *HTTPServer) CatalogServices(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	// Set default DC
	args := structs.DCSpecificRequest{}
	if done := s.parse(resp, req, &args.Datacenter, &args.QueryOptions); done {
		return nil, nil
	}

	var out structs.IndexedServices
	defer setMeta(resp, &out.QueryMeta)
	if err := s.agent.RPC("Catalog.ListServices", &args, &out); err != nil {
		return nil, err
	}
	return out.Services, nil
}

func (s *HTTPServer) CatalogArchetypes(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	// Set default DC
	args := structs.DCSpecificRequest{}
	if done := s.parse(resp, req, &args.Datacenter, &args.QueryOptions); done {
		return nil, nil
	}

	var out structs.IndexedArchetypes
	defer setMeta(resp, &out.QueryMeta)
	if err := s.agent.RPC("Catalog.ListArchetypes", &args, &out); err != nil {
		return nil, err
	}
	return out.Archetypes, nil
}

func (s *HTTPServer) CatalogServiceNodes(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	// Set default DC
	args := structs.ServiceSpecificRequest{}
	if done := s.parse(resp, req, &args.Datacenter, &args.QueryOptions); done {
		return nil, nil
	}

	// Check for a tag
	params := req.URL.Query()
	if _, ok := params["tag"]; ok {
		args.ServiceTag = params.Get("tag")
		args.TagFilter = true
	}

	// Pull out the service name
	args.ServiceName = strings.TrimPrefix(req.URL.Path, "/v1/catalog/service/")
	if args.ServiceName == "" {
		resp.WriteHeader(400)
		resp.Write([]byte("Missing service name"))
		return nil, nil
	}

	// Make the RPC request
	var out structs.IndexedServiceNodes
	defer setMeta(resp, &out.QueryMeta)
	if err := s.agent.RPC("Catalog.ServiceNodes", &args, &out); err != nil {
		return nil, err
	}
	return out.ServiceNodes, nil
}

func (s *HTTPServer) CatalogArchetypeNodes(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	// Set default DC
	args := structs.ArchetypeSpecificRequest{}
	if done := s.parse(resp, req, &args.Datacenter, &args.QueryOptions); done {
		return nil, nil
	}

	// Check for a tag
	params := req.URL.Query()
	if _, ok := params["tag"]; ok {
		args.ArchetypeTag = params.Get("tag")
		args.TagFilter = true
	}

	// Pull out the service name
	args.ArchetypeName = strings.TrimPrefix(req.URL.Path, "/v1/catalog/archetype/")
	if args.ArchetypeName == "" {
		resp.WriteHeader(400)
		resp.Write([]byte("Missing archetype name"))
		return nil, nil
	}

	// Make the RPC request
	var out structs.IndexedArchetypeNodes
	defer setMeta(resp, &out.QueryMeta)
	if err := s.agent.RPC("Catalog.ArchetypeNodes", &args, &out); err != nil {
		return nil, err
	}
	return out.ArchetypeNodes, nil
}

func (s *HTTPServer) CatalogNodeServices(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	// Set default Datacenter
	args := structs.NodeSpecificRequest{}
	if done := s.parse(resp, req, &args.Datacenter, &args.QueryOptions); done {
		return nil, nil
	}

	// Pull out the node name
	args.Node = strings.TrimPrefix(req.URL.Path, "/v1/catalog/node/")
	if args.Node == "" {
		resp.WriteHeader(400)
		resp.Write([]byte("Missing node name"))
		return nil, nil
	}

	// Make the RPC request
	var out structs.IndexedNodeServices
	defer setMeta(resp, &out.QueryMeta)
	if err := s.agent.RPC("Catalog.NodeServices", &args, &out); err != nil {
		return nil, err
	}
	return out.NodeServices, nil
}

func (s *HTTPServer) CatalogNodeArchetypes(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	// Set default Datacenter
	args := structs.NodeSpecificRequest{}
	if done := s.parse(resp, req, &args.Datacenter, &args.QueryOptions); done {
		return nil, nil
	}

	// Pull out the node name
	args.Node = strings.TrimPrefix(req.URL.Path, "/v1/catalog/node/")
	if args.Node == "" {
		resp.WriteHeader(400)
		resp.Write([]byte("Missing node name"))
		return nil, nil
	}

	// Make the RPC request
	var out structs.IndexedNodeArchetypes
	defer setMeta(resp, &out.QueryMeta)
	if err := s.agent.RPC("Catalog.NodeArchetypes", &args, &out); err != nil {
		return nil, err
	}
	return out.NodeArchetypes, nil
}
