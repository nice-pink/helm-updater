package helmupdater

import (
	"os"
	"strings"

	"github.com/nice-pink/goutil/pkg/log"
	"github.com/nice-pink/helm-updater/pkg/models"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/repo"
)

const (
	ENV_PREFIX          string = "HELM_UPATER_"
	ENV_DEFAULT         string = ENV_PREFIX + "PRIVATE_REPO"
	ENV_PASSWORD_SUFFIX string = "_PASSWORD"
	ENV_USERNAME_SUFFIX string = "_USERNAME"
)

func (u *Updater) getRemoteVersion(app models.App) string {
	entry := GetRepoEntry(app)
	err := u.helmClient.AddOrUpdateChartRepo(entry)
	if err != nil {
		log.Err(err, "add or update chart repo error")
		return ""
	}

	return u.getChartVersion(entry, app.System)
}

func (u *Updater) getChartVersion(entry repo.Entry, system models.SystemType) string {
	log.Info("Get chart for:", entry.Name)

	chart, _, err := u.helmClient.GetChart(entry.Name, &action.ChartPathOptions{
		Password:           entry.Password,
		PassCredentialsAll: entry.PassCredentialsAll,
		RepoURL:            entry.URL,
		Username:           entry.Username,
	})
	if err != nil {
		log.Err(err, "get chart version error", entry.Name)
		return ""
	}
	log.Info(chart.Metadata.Name, "found. App version:", chart.Metadata.AppVersion, ", Chart version:", chart.Metadata.Version)

	if system == models.K8s {
		// in this case use app version
		return chart.Metadata.AppVersion
	}
	return chart.Metadata.Version
}

func (u *Updater) getReleaseVersion(releaseName string) string {
	log.Info("get release for:", releaseName)

	release, err := u.helmClient.GetRelease(releaseName)
	if err != nil {
		log.Err(err, "get release error", releaseName)
		return ""
	}

	log.Info("release version:", release.Version, "chart version:", release.Chart.Metadata.Version)
	return release.Chart.Metadata.Version
}

func GetRepoEntry(app models.App) repo.Entry {
	// type Entry struct {
	// 	Name                  string `json:"name"`
	// 	URL                   string `json:"url"`
	// 	Username              string `json:"username"`
	// 	Password              string `json:"password"`
	// 	CertFile              string `json:"certFile"`
	// 	KeyFile               string `json:"keyFile"`
	// 	CAFile                string `json:"caFile"`
	// 	InsecureSkipTLSverify bool   `json:"insecure_skip_tls_verify"`
	// 	PassCredentialsAll    bool   `json:"pass_credentials_all"`
	// }
	username, password := GetRepoCredentials(app)
	return repo.Entry{
		Name:               app.Name,
		URL:                app.Repo,
		Username:           username,
		Password:           password,
		PassCredentialsAll: username != "" || password != "",
	}
}

func GetRepoCredentials(app models.App) (username string, password string) {
	if !app.Private {
		return "", ""
	}

	username = app.RepoUsername
	if username == "" {
		username = GetRepoCredentialsEnv(app, ENV_USERNAME_SUFFIX)
	}

	password = app.RepoPassword
	if password == "" {
		password = GetRepoCredentialsEnv(app, ENV_PASSWORD_SUFFIX)
	}

	return username, password
}

func GetRepoCredentialsEnv(app models.App, suffix string) string {
	varNamePrefix := ENV_PREFIX + strings.ToUpper(app.Name)
	envVal := os.Getenv(varNamePrefix + suffix)
	if envVal != "" {
		return envVal
	}
	return os.Getenv(ENV_DEFAULT + suffix)
}

// clean up

func ClearHelmCache(cachePath, repoFilePath string) error {
	return nil

	// this will lead to missing index files error!
	// err := os.RemoveAll(cachePath + "/")
	// if err != nil {
	// 	log.Err(err, "cannot clear cache")
	// }
	// err = os.Remove(repoFilePath)
	// if err != nil {
	// 	log.Err(err, "cannot delete repo file")
	// }
	// return err
}
