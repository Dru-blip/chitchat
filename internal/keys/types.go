package keys

import (
	"chitchat/internal/db/sqlc"
	"context"

	"github.com/google/uuid"
)

type SignedPreKey struct {
	ID        int32  `json:"id" validate:"required,gt=0"`
	Key       string `json:"key" validate:"required"`
	Signature string `json:"signature" validate:"required"`
}

type UploadPayload struct {
	Prekeys      []string     `json:"prekeys" validate:"required,min=1"`
	PrekeyIds    []int32      `json:"prekeyIds" validate:"required,min=1"`
	SignedPreKey SignedPreKey `json:"signedPreKey" validate:"required"`
}

type Repository interface {
	GetKeybundle(ctx context.Context, userID uuid.UUID) (sqlc.GetKeybundleRow, error)
	InsertPreKeys(ctx context.Context, arg sqlc.InsertPreKeysParams) error
	// DeletePrekey(ctx context.Context, arg sqlc.DeletePrekeyParams) (sqlc.DevicePrekey, error)
}
