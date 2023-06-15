package helpers

import (
	"context"

	"github.com/gotodo/config"
	"github.com/gotodo/models"
)

func GetUserIdFromContext(ctx context.Context) uint {
	email := ""
	if emailVal := ctx.Value(config.ContextUserKey); emailVal != nil {
		if mail, ok := emailVal.(string); ok {
			email = mail
		}
	}
	var user models.User
	models.DB.Where("email = ?", email).First(&user)
	return user.ID
}
