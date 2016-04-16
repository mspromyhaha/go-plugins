package eureka

import (
	"encoding/json"
	"errors"

	"github.com/hudl/fargo"
	"github.com/micro/go-micro/registry"
)

func appToService(app *fargo.Application) []*registry.Service {
	serviceMap := make(map[string]*registry.Service)

	for _, instance := range app.Instances {
		id := instance.Id()
		addr := instance.IPAddr
		port := instance.Port

		var version string
		var metadata map[string]string
		var endpoints []*registry.Endpoint

		// get version
		k, err := instance.Metadata.GetString("version")
		if err != nil {
			continue
		}

		k, err = instance.Metadata.GetString("endpoints")
		if err == nil {
			json.Unmarshal([]byte(k), &endpoints)
		}

		k, err = instance.Metadata.GetString("metadata")
		if err == nil {
			json.Unmarshal([]byte(k), &metadata)
		}

		// get existing service
		service, ok := serviceMap[version]
		if !ok {
			// create new if doesn't exist
			service = &registry.Service{
				Name:      app.Name,
				Version:   version,
				Endpoints: endpoints,
			}
		}

		// append node
		service.Nodes = append(service.Nodes, &registry.Node{
			Id:       id,
			Address:  addr,
			Port:     port,
			Metadata: metadata,
		})

		// save
		serviceMap[version] = service
	}

	var services []*registry.Service

	for _, service := range serviceMap {
		services = append(services, service)
	}

	return services
}

// only parses first node
func serviceToInstance(service *registry.Service) (*fargo.Instance, error) {
	if len(service.Nodes) == 0 {
		return nil, errors.New("Require nodes")
	}

	node := service.Nodes[0]

	instance := &fargo.Instance{
		App:    service.Name,
		IPAddr: node.Address,
		Port:   node.Port,
		Status: fargo.UP,
		UniqueID: func(i fargo.Instance) string {
			return node.Id
		},
		DataCenterInfo: fargo.DataCenterInfo{Name: "micro"},
	}

	// set version
	instance.SetMetadataString("version", service.Version)

	// set endpoints
	if b, err := json.Marshal(service.Endpoints); err == nil {
		instance.SetMetadataString("endpoints", string(b))
	}

	// set metadata
	if b, err := json.Marshal(node.Metadata); err == nil {
		instance.SetMetadataString("metadata", string(b))
	}

	return instance, nil
}
