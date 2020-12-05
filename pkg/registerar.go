package pkg

import (
	Client "github.com/go-kit/kit/sd/consul"
	stdConsul "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

// Registrar registers service instance aliveness information to Consul.
type Registrar struct {
	client       Client.Client
	registration *stdConsul.AgentServiceRegistration
	logger       *logrus.Logger
	registerErr  error
}

// NewRegistrar returns a Consul Registrar acting on the provided catalog
// registration.
func NewRegistrar(client Client.Client, r *stdConsul.AgentServiceRegistration, logger *logrus.Logger) *Registrar {
	return &Registrar{
		client:       client,
		registration: r,
		logger:       logger,
		registerErr:  nil,
	}
}

// Register implements sd.Registrar interface.
func (p *Registrar) Register() {
	if err := p.client.Register(p.registration); err != nil {
		p.logger.Info("service Discovery status :: ", "false")
		//panic("service not find serviceDiscovery to register it")
		p.registerErr = ErrNotRegisterServiceDiscovery
	} else {
		p.logger.Info("service Discovery status :: ", " true")
	}
}

// Deregister implements sd.Registrar interface.
func (p *Registrar) Deregister() {
	if err := p.client.Deregister(p.registration); err != nil {
		p.logger.Info("err", err)
		p.registerErr = ErrNotDeRegisterServiceDiscovery
	} else {
		p.logger.Info("action", "deregister")
	}
}
