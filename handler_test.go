package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock the RPCClient interface
type MockRPCClient struct {
	mock.Mock
}

func (m *MockRPCClient) sendRPCRequest(method string, params []interface{}) (interface{}, error) {
	args := m.Called(method, params)
	return args.Get(0), args.Error(1)
}

// Test getBlockNumberHandler
func TestGetBlockNumberHandler_Success(t *testing.T) {
	// Create a mock RPCClient instance
	mockRPC := new(MockRPCClient)
	// Mock the return value for the sendRPCRequest method
	mockRPC.On("sendRPCRequest", "eth_blockNumber", nil).Return("0x123456", nil)

	// Create a request and a recorder
	req, err := http.NewRequest("GET", "/blockNumber", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Setup the handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getBlockNumberHandler(w, r, mockRPC)
	})

	// Execute the request
	handler.ServeHTTP(rr, req)

	// Validate the response
	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "0x123456", response["blockNumber"])
}

// Test getBlockByNumberHandler
func TestGetBlockNumberHandler_Failure(t *testing.T) {
	// Create a mock RPCClient instance
	mockRPC := new(MockRPCClient)
	// Mock the return value for the sendRPCRequest method
	mockRPC.On("sendRPCRequest", "eth_blockNumber", nil).Return(nil, nil)  // Return nil here to match expected behavior

	// Create a request and a recorder
	req, err := http.NewRequest("GET", "/blockNumber", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Setup the handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		getBlockNumberHandler(w, r, mockRPC)
	})

	// Execute the request
	handler.ServeHTTP(rr, req)

	// Validate the response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure the error response matches
	assert.Equal(t, `{"error": "Failed to get block number"}`, rr.Body.String())
}
