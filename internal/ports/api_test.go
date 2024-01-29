package ports

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/ashermp9/fiskaly-test-task/internal/app"
	"github.com/ashermp9/fiskaly-test-task/internal/domain"
	"github.com/ashermp9/fiskaly-test-task/internal/storage"
	"go.uber.org/zap"
)

func TestServer(t *testing.T) {
	loggerZap, _ := zap.NewDevelopment()
	logger := loggerZap.Sugar()
	realStorage := storage.NewStorage()
	appService := app.NewAPIService(realStorage)
	server := NewServer(logger, appService, 8080)

	t.Run("CreateDevice", func(t *testing.T) { testCreateDevice(t, server) })
	t.Run("SignTransaction", func(t *testing.T) { testSignTransaction(t, server) })
	t.Run("HealthCheck", func(t *testing.T) { testHealthCheckHandler(t, server) })

}

func testCreateDevice(t *testing.T, server *Server) {
	deviceRequest := domain.CreateDeviceRequest{
		ID:        "test-device-id",
		Algorithm: domain.AlgorithmRSA,
		Label:     "Test Device",
	}
	requestBody, _ := json.Marshal(deviceRequest)
	request := httptest.NewRequest(http.MethodPost, "/api/v0/create-device", bytes.NewBuffer(requestBody))
	responseRecorder := httptest.NewRecorder()

	server.CreateSignatureDeviceHandler(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	// Perform checks on the created device
	var createdDevice domain.SignatureDevice
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &createdDevice)
	if err != nil {
		t.Fatalf("error unmarshalling response body: %v", err)
	}

	if createdDevice.ID != deviceRequest.ID {
		t.Errorf("expected ID %v, got %v", deviceRequest.ID, createdDevice.ID)
	}
	if createdDevice.Algorithm != deviceRequest.Algorithm {
		t.Errorf("expected Algorithm %v, got %v", deviceRequest.Algorithm, createdDevice.Algorithm)
	}
	if createdDevice.Label != deviceRequest.Label {
		t.Errorf("expected Label %v, got %v", deviceRequest.Label, createdDevice.Label)
	}
	if len(createdDevice.PublicKey) == 0 {
		t.Error("public key should not be empty")
	}
	if len(createdDevice.PrivateKey) == 0 {
		t.Error("private key should not be empty")
	}
	if createdDevice.SignatureCounter != 0 {
		t.Errorf("expected SignatureCounter 0, got %v", createdDevice.SignatureCounter)
	}
	if createdDevice.LastSignature != "" {
		t.Error("expected LastSignature to be empty for a new device")
	}
}

func testHealthCheckHandler(t *testing.T, server *Server) {
	request := httptest.NewRequest(http.MethodGet, "/api/v0/health", nil)
	responseRecorder := httptest.NewRecorder()

	server.HealthCheckHandler(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var healthStatus map[string]string
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &healthStatus)
	if err != nil {
		t.Fatalf("error unmarshalling response body: %v", err)
	}
	if healthStatus["status"] != "pass" || healthStatus["version"] != "v0" {
		t.Errorf("unexpected health status: got %v", healthStatus)
	}
}

func testSignTransaction(t *testing.T, server *Server) {
	// Ensure that a device with 'test-device-id' exists before testing this
	// ...

	// --- First Signature Request ---
	firstSignRequest := domain.SignTransactionRequest{
		DeviceID: "test-device-id",
		Data:     "data to be signed",
	}
	firstSignature := sendSignRequest(t, server, firstSignRequest)

	// --- Second Signature Request ---
	secondSignRequest := domain.SignTransactionRequest{
		DeviceID: "test-device-id",
		Data:     "additional data to be signed",
	}
	sendSignRequest(t, server, secondSignRequest, firstSignature)
}

func sendSignRequest(t *testing.T, server *Server, request domain.SignTransactionRequest, lastSignature ...string) string {
	requestBody, _ := json.Marshal(request)
	httpRequest := httptest.NewRequest(http.MethodPost, "/api/v0/sign-transaction", bytes.NewBuffer(requestBody))
	responseRecorder := httptest.NewRecorder()

	server.SignTransactionHandler(responseRecorder, httpRequest)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var signatureResponse domain.SignatureResponse
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &signatureResponse)
	if err != nil {
		t.Fatalf("error unmarshalling response body: %v", err)
	}

	if signatureResponse.Signature == "" {
		t.Error("expected non-empty signature")
	}

	// Construct expected SignedData format
	expectedLastSignature := ""
	if len(lastSignature) > 0 {
		expectedLastSignature = lastSignature[0]
	} else {
		expectedLastSignature = base64.StdEncoding.EncodeToString([]byte(request.DeviceID))
	}
	expectedSignedDataFormat := fmt.Sprintf("%d_%s_%s", len(lastSignature), request.Data, expectedLastSignature)
	if signatureResponse.SignedData != expectedSignedDataFormat {
		t.Errorf("expected SignedData to be '%s', got '%s'", expectedSignedDataFormat, signatureResponse.SignedData)
	}

	return signatureResponse.Signature
}

func TestConcurrentSignTransaction(t *testing.T) {
	// Shared setup
	loggerZap, _ := zap.NewDevelopment()
	logger := loggerZap.Sugar()
	realStorage := storage.NewStorage()
	appService := app.NewAPIService(realStorage)
	server := NewServer(logger, appService, 8080)

	deviceRequest := domain.CreateDeviceRequest{
		ID:        "test-device-id",
		Algorithm: domain.AlgorithmRSA,
		Label:     "Test Device",
	}
	requestBody, _ := json.Marshal(deviceRequest)
	request := httptest.NewRequest(http.MethodPost, "/api/v0/create-device", bytes.NewBuffer(requestBody))
	responseRecorder := httptest.NewRecorder()

	server.CreateSignatureDeviceHandler(responseRecorder, request)

	time.Sleep(time.Millisecond * 100)

	// Number of concurrent requests
	concurrentRequests := 10
	var wg sync.WaitGroup
	wg.Add(concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func(i int) {
			defer wg.Done()
			signRequest := domain.SignTransactionRequest{
				DeviceID: "test-device-id",
				Data:     fmt.Sprintf("data to be signed %d", i),
			}
			requestBody, _ := json.Marshal(signRequest)
			httpRequest := httptest.NewRequest(http.MethodPost, "/api/v0/sign-transaction", bytes.NewBuffer(requestBody))
			responseRecorder := httptest.NewRecorder()

			server.SignTransactionHandler(responseRecorder, httpRequest)

			if status := responseRecorder.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

		}(i)
	}

	wg.Wait() // Wait for all goroutines to finish

	// Retrieve the device and check its counter
	device, err := realStorage.GetDevice(context.Background(), "test-device-id")
	if err != nil {
		t.Fatalf("failed to retrieve device: %v", err)
	}

	if device.SignatureCounter != concurrentRequests {
		t.Errorf("expected signature counter to be %d, got %d", concurrentRequests, device.SignatureCounter)
	}
}
