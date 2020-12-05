# easy to use consul service discovery
This package is the implementation of consul service discovery register, 

## Dependencies

- [gokit](github.com/go-kit/kit v0.10.0) The library of micro service
- [hashicorp](github.com/hashicorp/consul/api) As the transport layer
- [lagrus](github.com/sirupsen/logrus) For serialization


## Quickstart
First get the project into your GOPATH using the following command:
```bash
$ go get "github.com/javid-dastgoshadeh/jconsulsd"
```

After doing above, you'll be able use this package
To use this package first import pakage

`import "github.com/javid-dastgoshadeh/jconsulsd"`

- ### define ServiceRegister struct

```
  r := ServiceRegister{\n
		ID:            "ID of service",
		Name:          "Name of service",
		Tags:          []string,
		ConsulAddress: "http://127.0.0.1:8500",
		ClientAddress: "localhost",
		HttpPort:      8080,
		GrpcPort:      8000,
		Logger:        &logrus.Logger{},
		TTL:           30,
		ConsulAgent:   &api.Agent{},
		Check:         &api.AgentServiceCheck{},	
	}
```

- ### register service discovery
```srv,err:=r.Register()```

this function return service interface and error value 

- ### now we can use methods to get service

 ```
 srv.GetServerAddressByNameAndTag(serviceName, tag string)

 srv.GetAllRegisteredService()

 srv.GetFirstServiceByTag(filter string)

 srv.GetFirstServiceByID(filter string)

 srv.GetServicesByTag(filter string)

 srv.GetServicesByName(filter string)

 srv.GetFirstServiceContainTags(tags []string)

 ```
 