package domain

type Algorithm string

const (
	AlgorithmRSA Algorithm = "RSA"
	AlgorithmECC Algorithm = "ECC"
)

type SignatureDevice struct {
	ID               string    // Unique identifier, e.g., UUID
	Algorithm        Algorithm // 'RSA' or 'ECC'
	PublicKey        []byte    // Encoded public key
	PrivateKey       []byte    // Encoded private key, should be securely stored
	Label            string    // User-provided label for the device
	SignatureCounter int       // Counts the number of signatures made
	LastSignature    string    // Last signed message
}

type CreateDeviceRequest struct {
	ID        string
	Algorithm Algorithm // 'RSA' or 'ECC'
	Label     string    // Optional label for the device
}

type CreateDeviceResponse struct {
	Device SignatureDevice // The newly created device
}

type SignTransactionRequest struct {
	DeviceID string // The ID of the signature device to use
	Data     string // The data to be signed
}

type SignatureResponse struct {
	Signature  string // The base64 encoded signature
	SignedData string // The original data with signature counter and last signature
}
