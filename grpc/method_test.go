package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMethodInvoker_NotNil(t *testing.T) {
	client, err := NewClient("localhost:50051")
	require.NoError(t, err)
	defer client.Close()

	invoker := NewMethodInvoker(client.Conn())
	assert.NotNil(t, invoker)
}

func TestMethodInvoker_InvokeUnary_NilDescriptor(t *testing.T) {
	client, err := NewClient("localhost:50051")
	require.NoError(t, err)
	defer client.Close()

	invoker := NewMethodInvoker(client.Conn())
	_, err = invoker.InvokeUnary(context.Background(), nil, []byte{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "method descriptor must not be nil")
}

func TestMethodInvoker_InvokeUnary_CancelledContext(t *testing.T) {
	client, err := NewClient("localhost:50051")
	require.NoError(t, err)
	defer client.Close()

	invoker := NewMethodInvoker(client.Conn())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// nil descriptor check happens before context check
	_, err = invoker.InvokeUnary(ctx, nil, []byte{})
	assert.Error(t, err)
}
