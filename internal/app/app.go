package app

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/ashermp9/fiskaly-test-task/internal/domain"
)

type APIStorage interface {
	AddDevice(ctx context.Context, device domain.SignatureDevice)
	GetDevice(ctx context.Context, id string) (domain.SignatureDevice, error)
	GenerateKeys(ctx context.Context, algorithm domain.Algorithm) ([]byte, []byte, error)
	SignTransaction(ctx context.Context, deviceID string, data []byte) ([]byte, error)
	LockDevice(ctx context.Context, deviceID string)
	UnlockDevice(ctx context.Context, deviceID string)
}

type APIService struct {
	storage APIStorage
}

func NewAPIService(storage APIStorage) *APIService {
	return &APIService{storage: storage}
}

func (app *APIService) CreateDevice(
	ctx context.Context, request domain.CreateDeviceRequest,
) (domain.SignatureDevice, error) {
	publicKey, privateKey, err := app.storage.GenerateKeys(ctx, request.Algorithm)
	if err != nil {
		return domain.SignatureDevice{}, err
	}

	device := domain.SignatureDevice{
		ID:               request.ID,
		Algorithm:        request.Algorithm,
		PublicKey:        publicKey,
		PrivateKey:       privateKey,
		Label:            request.Label,
		SignatureCounter: 0,
		LastSignature:    "",
	}

	app.storage.AddDevice(ctx, device)
	return device, nil
}

func (app *APIService) SignTransaction(
	ctx context.Context, request domain.SignTransactionRequest,
) (domain.SignatureResponse, error) {

	app.storage.LockDevice(ctx, request.DeviceID)
	defer app.storage.UnlockDevice(ctx, request.DeviceID)

	device, err := app.storage.GetDevice(ctx, request.DeviceID)
	if err != nil {
		return domain.SignatureResponse{}, err
	}

	// Formulate the data to be signed
	var lastSignature string
	if device.SignatureCounter == 0 {
		// Use base64-encoded device ID for the first transaction
		lastSignature = base64.StdEncoding.EncodeToString([]byte(device.ID))
	} else {
		lastSignature = base64.StdEncoding.EncodeToString([]byte(device.LastSignature))
	}
	dataToBeSigned := fmt.Sprintf("%d_%s_%s", device.SignatureCounter, request.Data, lastSignature)

	// Sign the data
	signature, err := app.storage.SignTransaction(ctx, request.DeviceID, []byte(dataToBeSigned))
	if err != nil {
		return domain.SignatureResponse{}, err
	}

	// Update the device's signature counter and last signature
	device.SignatureCounter++
	device.LastSignature = string(signature)
	app.storage.AddDevice(ctx, device)

	return domain.SignatureResponse{
		Signature:  base64.StdEncoding.EncodeToString(signature),
		SignedData: dataToBeSigned,
	}, nil
}
