package entities

// PutFileRequest is a struct that represents a put file request parameters.
type PutFileRequest struct {
	ClientID string `json:"client_id" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Data     []byte `json:"data" validate:"required"`
}

// GetFileRequest is a struct that represents a get file request parameters.
type GetFileRequest struct {
	ClientID string `json:"client_id" validate:"required"`
	Name     string `json:"name" validate:"required"`
}

// GetFileResponse is a struct that represents a get file response parameters.
type GetFileResponse struct {
	Name string `json:"name"`
	Data []byte `json:"data,omitempty"`
}
