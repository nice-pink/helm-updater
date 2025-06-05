package notify

import (
	"os"

	"github.com/nice-pink/goutil/pkg/log"
	"github.com/nice-pink/helm-updater/pkg/models"
	"github.com/nice-pink/slack-app/pkg/send"
)

type SlackClient struct {
	client    *send.Client
	channelId string
}

func NewSlackClient(token string) *SlackClient {
	return &SlackClient{
		client:    send.NewClient(token),
		channelId: os.Getenv("HELM_UPDATER_SLACK_CHANNEL_ID"),
	}
}

func (c *SlackClient) ShouldNotify(config models.Notication) bool {
	return config.Enable && (config.ChannelId != "" || c.channelId != "")
}

func (c *SlackClient) SendNotification(config models.Notication, app models.App, version string, updated bool) error {
	if !c.ShouldNotify(config) {
		return nil
	}

	msg := c.getMessage(config, app, version, updated)
	log.Info("Send notification:", msg.Text)
	return c.client.SendMsg(msg)
}

func (c *SlackClient) getMessage(config models.Notication, app models.App, version string, updated bool) send.Msg {
	// prefer env var for notification webhook
	channelId := os.Getenv("HELM_UPDATER_SLACK_CHANNEL_ID")
	if channelId == "" {
		channelId = config.ChannelId
	}

	if updated {
		// updated version
		return send.Msg{
			Header:    "🚀 Updated " + app.Name + " to version " + version,
			Text:      "Updated using helm-updater.",
			Color:     "#34eb8c",
			ChannelId: channelId,
		}
	}

	// new version available but not updated
	return send.Msg{
		Header:    "New version avaiable for " + app.Name,
		Text:      "Version available " + version,
		Color:     "#349ceb",
		ChannelId: channelId,
	}
}
