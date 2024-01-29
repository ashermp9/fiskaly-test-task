package types

import (
	"github.com/ashermp9/fiskaly-test-task/internal/domain"
)

func ConvertToDomainCreateDeviceRequest(apiRequest CreateDeviceRequest) domain.CreateDeviceRequest {
	return domain.CreateDeviceRequest{
		ID:        apiRequest.ID,
		Algorithm: domain.Algorithm(apiRequest.Algorithm),
		Label:     apiRequest.Label,
	}
}

func ConvertToDomainSignTransactionRequest(apiRequest SignTransactionRequest) domain.SignTransactionRequest {
	return domain.SignTransactionRequest{
		DeviceID: apiRequest.DeviceID,
		Data:     apiRequest.Data,
	}
}
