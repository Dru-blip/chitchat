package auth

import (
	"chitchat/internal/db/sqlc"
	"chitchat/internal/utils"
	"context"

	"github.com/redis/go-redis/v9"
)

func SetOnboardingToken(ctx context.Context, client *redis.Client, user *sqlc.User) (string, error) {
	token, err := utils.GenerateMagicLinkToken()
	if err != nil {
		return "", err
	}

	options := new(redis.HSetEXOptions)
	options.ExpirationType = redis.HSetEXExpirationEX
	options.ExpirationVal = 2592000

	_, err = client.HSetEXWithArgs(ctx, "onboarding:"+user.ID.String(), options, "email", user.Email, "id", user.ID.String(), "token", token).Result()

	if err != nil {
		return "", err
	}
	return token, nil
}

func RemoveOnboardingToken(ctx context.Context, client *redis.Client, userID string) error {
	_, err := client.HDel(ctx, "onboarding:"+userID, "email", "token", "id").Result()
	return err
}
