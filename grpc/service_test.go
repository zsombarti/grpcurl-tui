package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServiceExplorer_NotNil(t *testing.T) {
	client, err := NewClient("localhost:50051")
	require.NoError(t, err)
	defer client.Close()

	reflector := NewReflector(client.Conn())
	explorer := NewServiceExplorer(reflector)
	assert.NotNil(t, explorer)
}

func TestServiceExplorer_ListServices_Timeout(t *testing.T) {
	client, err := NewClient("localhost:50051")
	require.NoError(t, err)
	defer client.Close()

	reflector := NewReflector(client.Conn())
	explorer := NewServiceExplorer(reflector)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	time.Sleep(5 * time.Millisecond)

	_, err = explorer.ListServices(ctx)
	assert.Error(t, err)
}

func TestServiceExplorer_ListServices_CancelledContext(t *testing.T) {
	client, err := NewClient("localhost:50051")
	require.NoError(t, err)
	defer client.Close()

	reflector := NewReflector(client.Conn())
	explorer := NewServiceExplorer(reflector)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = explorer.ListServices(ctx)
	assert.Error(t, err)
}
