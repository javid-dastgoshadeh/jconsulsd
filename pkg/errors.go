package pkg

import "errors"

var (
	ErrNotfoundServiceDiscovery      = errors.New("not found target service in service discovery")
	ErrNotRegisterServiceDiscovery   = errors.New("not Register Service Discovery")
	ErrNotDeRegisterServiceDiscovery = errors.New("not DeRegister Service Discovery")
)
