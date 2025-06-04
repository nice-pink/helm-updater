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

func SendNotification(config models.Notication, app models.App, version string, updated bool) error {
	if !ShouldNotify(config) {
		return nil
	}

	msg := getMessage(config, app, version, updated)
	log.Info("Send notification:", msg.Text)
	return notify.Send(msg)
}

func getMessage(config models.Notication, app models.App, version string, updated bool) notify.SlackMessage {
	// prefer env var for notification webhook
	url := os.Getenv("HELM_UPDATER_NOTIFY_WEBHOOK")
	if url == "" {
		url = config.Webhook
	}

	if updated {
		// updated version
		return notify.SlackMessage{
			Text:  "🚀 Updated " + app.Name + " to version " + version,
			Info:  "Updated using helm-updater.",
			Color: "#34eb8c",
			Url:   url,
		}
	}

	// new version available but not updated
	return notify.SlackMessage{
		Text:  "New version avaiable for " + app.Name,
		Info:  "Version available " + version,
		Color: "#349ceb",
		Url:   url,
	}
}
