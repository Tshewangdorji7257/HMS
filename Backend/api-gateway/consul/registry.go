package consul

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/hashicorp/consul/api"
)

var consulClient *api.Client

// InitConsul initializes the Consul client and registers the service
func InitConsul() error {
	consulHost := os.Getenv("CONSUL_HOST")
	if consulHost == "" {
		consulHost = "localhost"
	}

	consulPort := os.Getenv("CONSUL_PORT")
	if consulPort == "" {
		consulPort = "8500"
	}

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", consulHost, consulPort)

	var err error
	consulClient, err = api.NewClient(config)
	if err != nil {
		return fmt.Errorf("failed to create consul client: %v", err)
	}

	log.Printf("✅ Connected to Consul at %s", config.Address)
	return nil
}

// RegisterService registers the service with Consul
func RegisterService() error {
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "auth-service"
	}

	serviceID := os.Getenv("SERVICE_ID")
	if serviceID == "" {
		serviceID = "auth-service-1"
	}

	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8001"
	}
	port, _ := strconv.Atoi(portStr)

	grpcPortStr := os.Getenv("GRPC_PORT")
	if grpcPortStr == "" {
		grpcPortStr = "9001"
	}
	grpcPort, _ := strconv.Atoi(grpcPortStr)

	// Get container IP
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	registration := &api.AgentServiceRegistration{
		ID:   serviceID,
		Name: serviceName,
		Port: port,
		Tags: []string{"api", "v1", "go"},
		Meta: map[string]string{
			"grpc_port": grpcPortStr,
			"version":   "1.0.0",
		},
		Check: &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/health", hostname, port),
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}

	err = consulClient.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("failed to register service: %v", err)
	}

	log.Printf("✅ Service registered: %s (ID: %s) on port %d", serviceName, serviceID, port)
	log.Printf("   gRPC port: %d", grpcPort)
	return nil
}

// DeregisterService deregisters the service from Consul
func DeregisterService() error {
	serviceID := os.Getenv("SERVICE_ID")
	if serviceID == "" {
		serviceID = "auth-service-1"
	}

	err := consulClient.Agent().ServiceDeregister(serviceID)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %v", err)
	}

	log.Printf("✅ Service deregistered: %s", serviceID)
	return nil
}

// DiscoverService discovers a service by name and returns its address
func DiscoverService(serviceName string) (string, error) {
	services, _, err := consulClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", fmt.Errorf("failed to discover service: %v", err)
	}

	if len(services) == 0 {
		return "", fmt.Errorf("no healthy instances found for service: %s", serviceName)
	}

	// Return the first healthy instance
	service := services[0]
	address := fmt.Sprintf("http://%s:%d", service.Service.Address, service.Service.Port)
	
	if service.Service.Address == "" {
		address = fmt.Sprintf("http://%s:%d", service.Node.Address, service.Service.Port)
	}

	return address, nil
}

// GetGRPCAddress gets the gRPC address for a service
func GetGRPCAddress(serviceName string) (string, error) {
	services, _, err := consulClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", fmt.Errorf("failed to discover service: %v", err)
	}

	if len(services) == 0 {
		return "", fmt.Errorf("no healthy instances found for service: %s", serviceName)
	}

	service := services[0]
	grpcPort := service.Service.Meta["grpc_port"]
	if grpcPort == "" {
		return "", fmt.Errorf("grpc_port not found in service metadata")
	}

	address := fmt.Sprintf("%s:%s", service.Service.Address, grpcPort)
	if service.Service.Address == "" {
		address = fmt.Sprintf("%s:%s", service.Node.Address, grpcPort)
	}

	return address, nil
}
