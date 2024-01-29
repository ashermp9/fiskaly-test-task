package types

import "fmt"

// Validator interface for request validation
type Validator interface {
	Validate() error
}

type CreateDeviceRequest struct {
	ID        string `json:"id"`
	Algorithm string `json:"algorithm"`
	Label     string `json:"label,omitempty"`
}

func (r CreateDeviceRequest) Validate() error {
	if r.ID == "" {
		return fmt.Errorf("ID is required")
	}
	if r.Algorithm != "RSA" && r.Algorithm != "ECC" {
		return fmt.Errorf("invalid algorithm: must be 'RSA' or 'ECC'")
	}
	return nil
}

type SignTransactionRequest struct {
	DeviceID string `json:"deviceId"`
	Data     string `json:"data"`
}

// Validate performs input validation on a SignTransactionRequest.
func (r SignTransactionRequest) Validate() error {
	if r.DeviceID == "" {
		return fmt.Errorf("DeviceID is required")
	}
	if r.Data == "" {
		return fmt.Errorf("data is required")
	}
	return nil
}
