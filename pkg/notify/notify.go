package notify

import (
	"os"

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

	// prefer env var for notification webhook
	url := os.Getenv("HELM_UPDATER_NOTIFY_WEBHOOK")
	if url == "" {
		url = config.Webhook
	}

	msg := notify.SlackMessage{
		Text:  "🚀 Updated " + app.Name + " to version " + version,
		Info:  "Updated using helm-updater.",
		Color: "#34eb8c",
		Url:   url,
	}

	log.Info("Send notification:", msg.Text)
	return notify.Send(msg)
}
