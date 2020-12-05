package jconsulsd

import (
	consuls "github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/javid-dastgoshadeh/jconsulsd/pkg"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

//ServiceRegister
type ServiceRegister struct {
	ID            string                 //define ID for service
	Name          string                 //define Name for service
	Tags          []string               //define Tags for service
	ConsulAddress string                 //Consul service discover address
	ClientAddress string                 //service host address
	HttpPort      int                    // service http port address
	GrpcPort      int                    // service Grpc port address
	Logger        *logrus.Logger         //logger that use from logrus
	TTL           time.Duration          // time to live duration that define health check period
	ConsulAgent   *api.Agent             // Agent can be used to query the Agent endpoint
	Check         *api.AgentServiceCheck // AgentServiceCheck is used to define a node or service level check
}

//Register service on service discovery
func (r *ServiceRegister) Register() (pkg.Service, error) {

	// Service discovery
	var client consuls.Client
	{
		consulConfig := api.DefaultConfig()
		consulConfig.Address = r.ConsulAddress
		consulClient, err := api.NewClient(consulConfig)
		if err != nil {
			r.Logger.Info("err", err)
			os.Exit(1)
			return nil, err
		}
		r.ConsulAgent = consulClient.Agent() //check health
		client = consuls.NewClient(consulClient)
	}

	asr := &api.AgentServiceRegistration{
		ID:      r.ID,
		Name:    r.Name,
		Address: r.ClientAddress,
		Port:    r.HttpPort,
		Tags:    r.Tags,
		Check:   r.Check,
	}

	i := &pkg.ServiceInfo{
		ID:            r.ID,
		Name:          r.Name,
		ConsulAddress: r.ConsulAddress,
		Logger:        r.Logger,
		TTL:           r.TTL,
		ConsulAgent:   r.ConsulAgent,
		Check:         r.Check,
	}
	//create new service to register service an get service instance
	service := pkg.NewService(i, client, asr, i.Logger)
	return service, nil
}
