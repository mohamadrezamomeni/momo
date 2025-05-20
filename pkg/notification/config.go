package notification

type Notification struct {
	Telegram TelegramConfig `koanf:"telegram"`
}
