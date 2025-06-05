package notify

import (
	"os"

	"github.com/nice-pink/helm-updater/pkg/models"
)

type Client interface {
	ShouldNotify(config models.Notication) bool
	SendNotification(config models.Notication, app models.App, version string, updated bool) error
}

func NewClient(config models.Notication) Client {
	if config.Enable && (config.ChannelId != "" || os.Getenv("HELM_UPDATER_SLACK_CHANNEL_ID") != "") {
		return NewSlackClient(config.Token)
	}
	if config.Enable && (config.Webhook != "" || os.Getenv("HELM_UPDATER_NOTIFY_WEBHOOK") != "") {
		return NewNotifierClient()
	}
	return &Dummy{}
}

// dummy

type Dummy struct{}

func (c *Dummy) ShouldNotify(config models.Notication) bool {
	return false
}

func (c *Dummy) SendNotification(config models.Notication, app models.App, version string, updated bool) error {
	return nil
}
