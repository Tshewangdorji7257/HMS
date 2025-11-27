package grpc

import (
	authpb "auth-service/proto/auth"
	buildingpb "building-service/proto/building"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	authClient     authpb.AuthServiceClient
	buildingClient buildingpb.BuildingServiceClient
	authConn       *grpc.ClientConn
	buildingConn   *grpc.ClientConn
)

// InitClients initializes gRPC clients for the API Gateway
func InitClients() error {
	authServiceURL := os.Getenv("AUTH_GRPC_URL")
	if authServiceURL == "" {
		authServiceURL = "auth-service:9001"
	}

	buildingServiceURL := os.Getenv("BUILDING_GRPC_URL")
	if buildingServiceURL == "" {
		buildingServiceURL = "building-service:9002"
	}

	// Connect to auth service
	var err error
	authConn, err = grpc.Dial(authServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to auth service: %v", err)
	}
	authClient = authpb.NewAuthServiceClient(authConn)
	log.Printf("✅ API Gateway connected to Auth gRPC service at %s", authServiceURL)

	// Connect to building service
	buildingConn, err = grpc.Dial(buildingServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to building service: %v", err)
	}
	buildingClient = buildingpb.NewBuildingServiceClient(buildingConn)
	log.Printf("✅ API Gateway connected to Building gRPC service at %s", buildingServiceURL)

	return nil
}

// CloseClients closes all gRPC connections
func CloseClients() {
	if authConn != nil {
		authConn.Close()
	}
	if buildingConn != nil {
		buildingConn.Close()
	}
}

// ValidateToken validates a JWT token via auth service
func ValidateToken(token string) (*authpb.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := authClient.ValidateToken(ctx, &authpb.ValidateTokenRequest{
		Token: token,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}

	if !resp.Valid {
		return nil, fmt.Errorf("invalid token: %s", resp.Message)
	}

	return resp.User, nil
}

// GetUserByID retrieves user information by ID
func GetUserByID(userID string) (*authpb.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := authClient.GetUserByID(ctx, &authpb.GetUserByIDRequest{
		UserId: userID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("user not found: %s", resp.Message)
	}

	return resp.User, nil
}

// GetBuildingByID retrieves building information by ID
func GetBuildingByID(buildingID string) (*buildingpb.Building, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := buildingClient.GetBuildingByID(ctx, &buildingpb.GetBuildingByIDRequest{
		BuildingId: buildingID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get building: %v", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("building not found: %s", resp.Message)
	}

	return resp.Building, nil
}

// GetRoomByID retrieves room information by ID
func GetRoomByID(buildingID, roomID string) (*buildingpb.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := buildingClient.GetRoomByID(ctx, &buildingpb.GetRoomByIDRequest{
		BuildingId: buildingID,
		RoomId:     roomID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get room: %v", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("room not found: %s", resp.Message)
	}

	return resp.Room, nil
}
