package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

// RegistryClient 注册客户端接口 方便后续扩展
type RegistryClient interface {
	Register(address string, port int, name string, tags []string, serviceId string) error
	Deregister(serviceId string) error
}

/*
未来可以这样扩展，调用方零改动
type EtcdRegistry struct { ... }
func (e *EtcdRegistry) Register(...) error { ... }  // 实现同一个接口

工厂函数返回不同实现，调用方无感知

	func NewRegistryClient(...) RegistryClient {
	    if useEtcd {
	        return &EtcdRegistry
	        {...}  返回不同的实现
	    }
	    return &Registry{...}
	}
*/
type Registry struct {
	Host string
	Port int
}

func NewRegistryClient(host string, port int) RegistryClient {
	return &Registry{
		Host: host,
		Port: port,
	}
}

func (r *Registry) Register(address string, port int, name string, tags []string, serviceId string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成检查对象
	checker := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", address, port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	registration := &api.AgentServiceRegistration{
		Name:    name,
		ID:      serviceId,
		Address: address,
		Port:    port,
		Tags:    tags,
		Check:   checker,
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}
	return nil
}

func (r *Registry) Deregister(serviceId string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	err = client.Agent().ServiceDeregister(serviceId)
	if err != nil {
		return err
	}
	return nil
}
