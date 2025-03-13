package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Polygon RPC endpoint
const rpcURL = "https://polygon-rpc.com/"

// JSON-RPC Request structure
type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params,omitempty"`
	ID      int           `json:"id"`
}

// JSON-RPC Response structure (Updated to handle objects)
type RPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result"`
}

// Function to send JSON-RPC requests
func sendRPCRequest(method string, params []interface{}) (interface{}, error) {
	requestBody := RPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      2,
	}

	reqBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(rpcURL, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Log response only for errors or specific issues
	if resp.StatusCode != http.StatusOK {
		log.Printf("🛠️ RPC Error Response: %s", string(body))
	}

	var rpcResponse RPCResponse
	err = json.Unmarshal(body, &rpcResponse)
	if err != nil {
		return nil, err
	}

	if rpcResponse.Result == nil {
		return nil, fmt.Errorf("RPC response returned null")
	}

	return rpcResponse.Result, nil
}
// Handler to get the latest block number
func getBlockNumberHandler(w http.ResponseWriter, r *http.Request) {
	blockNumber, err := sendRPCRequest("eth_blockNumber", nil)
	if err != nil {
		http.Error(w, `{"error": "Failed to get block number"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"blockNumber": blockNumber})
}

// Handler to get the latest block number
func getBlockByNumberHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blockNumber := vars["blockNumber"]

	log.Printf("📡 Requesting block: %s", blockNumber)

	params := []interface{}{blockNumber, true}
	blockDetails, err := sendRPCRequest("eth_getBlockByNumber", params)
	if err != nil {
		log.Printf("❌ Error fetching block: %v", err)
		http.Error(w, `{"error": "Failed to fetch block details"}`, http.StatusInternalServerError)
		return
	}

	if blockDetails == nil {
		log.Printf("⚠️ Empty block response for %s", blockNumber)
		http.Error(w, `{"error": "No block found"}`, http.StatusNotFound)
		return
	}

	// Convert blockDetails to a JSON object
	blockData, ok := blockDetails.(map[string]interface{})
	if !ok {
		http.Error(w, `{"error": "Invalid block data"}`, http.StatusInternalServerError)
		return
	}

	// Remove redundant blockHash from transactions
	if transactions, exists := blockData["transactions"].([]interface{}); exists {
		for _, tx := range transactions {
			if txMap, ok := tx.(map[string]interface{}); ok {
				delete(txMap, "blockHash") // Remove blockHash from each transaction
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"block": blockData})
}

// Main function to start the HTTP server
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/blockNumber", getBlockNumberHandler).Methods("GET")
	router.HandleFunc("/block/{blockNumber}", getBlockByNumberHandler).Methods("GET")

	fmt.Println("🚀 Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
