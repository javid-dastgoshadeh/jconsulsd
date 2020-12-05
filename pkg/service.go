package pkg

import (
	"fmt"
	"log"
	"time"

	consuls "github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

//ServiceInfo implement service
type ServiceInfo struct {
	ID            string                 //define ID for service
	Name          string                 //define Name for service
	Tags          []string               //define Tags for service
	ConsulAddress string                 //Consul service discover address
	Logger        *logrus.Logger         //logger that use from logrus
	TTL           time.Duration          // time to live duration that define health check period
	ConsulAgent   *api.Agent             // Agent can be used to query the Agent endpoint
	Check         *api.AgentServiceCheck // AgentServiceCheck is used to define a node or service level check
}

// bookService describes the service.
type Service interface {
	GetServerAddressByNameAndTag(serviceName, tag string) (string, string, error)
	GetAllRegisteredService() (map[string]*api.AgentService, error)
	GetFirstServiceByTag(filter string) (string, error)
	GetFirstServiceByID(filter string) (string, error)
	GetServicesByTag(filter string) ([]string, error)
	GetServicesByName(filter string) (map[string]*api.AgentService, error)
	GetFirstServiceContainTags(tags []string) (string, string, error)
}

// Using the service name and one of the service tags
func (serviceInfo *ServiceInfo) GetServerAddressByNameAndTag(serviceName, tag string) (string, string, error) {

	consulAddr := serviceInfo.ConsulAddress
	config := api.DefaultConfig()
	config.Address = consulAddr          // Assign address
	client, err := api.NewClient(config) // Creating consul client
	if err != nil {
		return "", "", err
	}

	// Getting service info from consul agent
	services, _, err := client.Health().Service(serviceName, tag, false, nil)

	if err != nil || len(services) == 0 {
		return "", "", err
	}
	// Extracting address and port from info
	serviceAddress := services[0].Service.Address
	servicePort := services[0].Service.Port
	serviceID := services[0].Service.ID
	// Return address in address:port format
	return fmt.Sprintf("%s:%d", serviceAddress, servicePort), serviceID, nil
}

//Get all Services that register on service discovery
func (serviceInfo *ServiceInfo) GetAllRegisteredService() (map[string]*api.AgentService, error) {

	var agents map[string]*api.AgentService

	consulAddr := serviceInfo.ConsulAddress

	config := api.DefaultConfig()
	config.Address = consulAddr          // Assign address
	client, err := api.NewClient(config) // Creating consul client
	if err != nil {
		return nil, err
	}
	agents, err = client.Agent().Services()

	if err != nil {
		return nil, err
	}
	return agents, nil
}

//Get first Service that find with tag name From service discovery
func (serviceInfo *ServiceInfo) GetFirstServiceByTag(filter string) (string, error) {
	consulAddr := serviceInfo.ConsulAddress

	var agents map[string]*api.AgentService

	config := api.DefaultConfig()
	config.Address = consulAddr          // Assign address
	client, err := api.NewClient(config) // Creating consul client
	if err != nil {
		return "", err
	}
	//
	filterQuery := fmt.Sprintf("%s in Tags", filter)
	agents, err = client.Agent().ServicesWithFilter(filterQuery)

	//agents, err = client.Agent().Services()

	if err != nil {
		return "", err
	}

	if len(agents) < 1 {
		return "", ErrNotfoundServiceDiscovery
	}
	var (
		serviceAddress string
		servicePort    int
	)

	for _, agent := range agents {
		serviceAddress = agent.Address
		servicePort = agent.Port
	}

	return fmt.Sprintf("%s:%d", serviceAddress, servicePort), nil

}

//GetFirstServiceByID
func (serviceInfo *ServiceInfo) GetFirstServiceByID(ID string) (string, error) {

	if ID == "" {
		return "", ErrNotfoundServiceDiscovery
	}
	consulAddr := serviceInfo.ConsulAddress
	var agents map[string]*api.AgentService
	config := api.DefaultConfig()
	config.Address = consulAddr          // Assign address
	client, err := api.NewClient(config) // Creating consul client
	if err != nil {
		return "", err
	}

	//filterQuery := fmt.Sprintf("ID == \"%s\"", ID)
	//agents, err = client.Agent().ServicesWithFilter(filterQuery)

	agents, err = client.Agent().Services()

	if err != nil {
		return "", err
	}

	serviceAddress := agents[ID].Address
	servicePort := agents[ID].Port

	return fmt.Sprintf("%s:%d", serviceAddress, servicePort), nil

}

//GetServicesByTag
func (serviceInfo *ServiceInfo) GetServicesByTag(filter string) ([]string, error) {

	var (
		allAgents []string
		agents    map[string]*api.AgentService
	)
	consulAddr := serviceInfo.ConsulAddress

	config := api.DefaultConfig()
	config.Address = consulAddr          // Assign address
	client, err := api.NewClient(config) // Creating consul client
	if err != nil {
		return nil, err
	}

	filterQuery := fmt.Sprintf("%s in Tags", filter)
	agents, err = client.Agent().ServicesWithFilter(filterQuery)

	if err != nil {
		return nil, err
	}
	for _, agent := range agents {
		serviceAddress := agent.Address
		servicePort := agent.Port
		allAgents = append(allAgents, fmt.Sprintf("%s:%d", serviceAddress, servicePort))
	}

	return allAgents, nil
}

//GetServicesByName
func (serviceInfo *ServiceInfo) GetServicesByName(filter string) (map[string]*api.AgentService, error) {

	var agents map[string]*api.AgentService

	consulAddr := serviceInfo.ConsulAddress
	config := api.DefaultConfig()
	config.Address = consulAddr          // Assign address
	client, err := api.NewClient(config) // Creating consul client
	if err != nil {
		return nil, err
	}

	filterQuery := fmt.Sprintf("Name == %s ", filter)
	agents, err = client.Agent().ServicesWithFilter(filterQuery)

	if err != nil {
		return nil, err
	}
	return agents, nil
}

//Get Service From Service discovery that contains all tag
func (serviceInfo *ServiceInfo) GetFirstServiceContainTags(tags []string) (string, string, error) {
	services, _ := serviceInfo.GetAllRegisteredService()
	for _, service := range services {
		for _, tag := range service.Tags {
			if IfSliceContainsString(tags, tag) {
				server, ID, err := serviceInfo.GetServerAddressByNameAndTag(service.Service, tags[0])
				return server, ID, err
			}
		}

	}
	return "", "", ErrNotfoundServiceDiscovery
}

// Update time to live
func (serviceInfo *ServiceInfo) updateTTL() {
	ticker := time.NewTicker(serviceInfo.TTL / 2)
	for range ticker.C {
		if agentErr := serviceInfo.ConsulAgent.FailTTL("service:"+serviceInfo.ID, ""); agentErr != nil {
			log.Println(agentErr)
		}

		if agentErr := serviceInfo.ConsulAgent.PassTTL("service:"+serviceInfo.ID, ""); agentErr != nil {
			log.Fatalln(agentErr)
		}
	}
}

//NewService
func NewService(serviceInfo *ServiceInfo, client consuls.Client, asr *api.AgentServiceRegistration, logger *logrus.Logger) Service {
	go serviceInfo.updateTTL()
	registrar := NewRegistrar(client, asr, logger)
	registrar.Register()
	return serviceInfo
}
