package notify

import (
	"github.com/nice-pink/goutil/pkg/log"
	"github.com/nice-pink/goutil/pkg/notify"
	"github.com/nice-pink/helm-updater/pkg/models"
)

func ShouldNotify(config models.Notication) bool {
	return config.Enable && config.Webhook != ""
}

func SendNotification(config models.Notication, app models.App, version string) error {
	if !ShouldNotify(config) {
		return nil
	}

	msg := notify.SlackMessage{
		Text:  "🚀 Updated " + app.Name + " to version " + version,
		Info:  "Updated using helm-updater.",
		Color: "#34eb8c",
		Url:   config.Webhook,
	}

	log.Info("Send notification:", msg.Text)
	return notify.Send(msg)
}
