package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"testing"
	"time"
)

type MockClient struct {
	handledRequests int
	mu              sync.Mutex
}

func (mc *MockClient) handleRequest() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.handledRequests++
}

func TestServer_Run(t *testing.T) {
	os.Setenv("MAX_CONNECTIONS", "20")
	defer os.Unsetenv("MAX_CONNECTIONS")

	os.Setenv("BROKER_HOST", "0.0.0.0")
	defer os.Unsetenv("BROKER_HOST")

	os.Setenv("BROKER_PORT", "9092")
	defer os.Unsetenv("BROKER_PORT")

	// Mock server and client
	mockClient := &MockClient{}
	server := &Server{
		host: "0.0.0.0",
		port: "9092",
	}

	// Create a context with cancellation
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the server in a goroutine
	go server.Run()

	// Allow the server some time to start
	time.Sleep(100 * time.Millisecond)

	// Dial the server to simulate a client
	addr := fmt.Sprintf("%s:%s", server.host, server.port)
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	for i := 0; i < 100; i++ {
		go func(i int) {
			conn.Write([]byte("Hello Server!"))
			mockClient.handleRequest()

		}(i)
	}

	// Allow time for the server to process
	time.Sleep(200 * time.Millisecond)

	// Verify that the server handled the request
	mockClient.mu.Lock()
	if mockClient.handledRequests != 100 {
		t.Errorf("Expected 100 handled request, got %d", mockClient.handledRequests)
	}
	mockClient.mu.Unlock()

	// Cancel the context to shut down the server
	cancel()
	time.Sleep(100 * time.Millisecond)
}