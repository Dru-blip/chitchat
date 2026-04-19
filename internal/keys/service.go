package keys

import (
	"chitchat/internal/db/sqlc"
	"context"

	"github.com/google/uuid"
)

type prekeyUpload struct {
	DeviceID    string
	PrekeyIds   []int32
	Prekeys     []string
	SignedKey   string
	SignedKeyID int32
	Signature   string
}

type Service interface {
	GetKeyBundle(ctx context.Context, user_id string) (*sqlc.GetKeybundleRow, error)
	UploadPrekeys(ctx context.Context, data prekeyUpload) error
	// DeletePrekey(ctx context.Context, device_id, user_id string) (*sqlc.DevicePrekey, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) UploadPrekeys(ctx context.Context, data prekeyUpload) error {

	//TODO: switch to parse
	did := uuid.MustParse(data.DeviceID)

	_, err := s.repo.InsertPreKeys(ctx, sqlc.InsertPreKeysParams{
		Deviceid:    did,
		Prekeys:     data.Prekeys,
		Prekeyids:   data.PrekeyIds,
		Signedkey:   data.SignedKey,
		Signedkeyid: data.SignedKeyID,
		Signature:   data.Signature,
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetKeyBundle(ctx context.Context, user_id string) (*sqlc.GetKeybundleRow, error) {
	uid := uuid.MustParse(user_id)
	key_bundle, err := s.repo.GetKeybundle(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &key_bundle, nil
}

// func (s *service) DeletePrekey(ctx context.Context, device_id, user_id string) (*sqlc.DevicePrekey, error) {
// 	uid, did := uuid.MustParse(user_id), uuid.MustParse(device_id)
// 	prekey, err := s.repo.DeletePrekey(ctx, sqlc.DeletePrekeyParams{
// 		UserID:   uid,
// 		DeviceID: did,
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &prekey, nil
// }
