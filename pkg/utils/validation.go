package utils

import "github.com/VladPetriv/tg_scanner/internal/model"

func ValidateTelegramUser(user *model.User) {
	if len(user.Username) > 0 {
		return
	}

	user.Username = user.FullName
}
