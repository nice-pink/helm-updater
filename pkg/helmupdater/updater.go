package helmupdater

import (
	"encoding/json"
	"os"
	"path"

	helmclient "github.com/mittwald/go-helm-client"
	"github.com/nice-pink/goutil/pkg/log"
	"github.com/nice-pink/helm-updater/pkg/models"
	"github.com/nice-pink/helm-updater/pkg/notify"
	"github.com/nice-pink/repo-services/pkg/manifest"
	"github.com/nice-pink/repo-services/pkg/util"
)

func Run(configFile string, gitFlags util.GitFlags) error {
	c := LoadConfig(configFile)

	// checkout repo?
	if gitFlags.Url != nil {
		err := util.GitClone(*gitFlags.Url, c.BaseFolder, gitFlags)
		if err != nil {
			return err
		}
	}

	// the client is only used to pull repos so most options don't really matter
	helmClient, err := helmclient.New(&helmclient.Options{
		Namespace:        "default",
		RepositoryCache:  "bin/.helmcache",
		RepositoryConfig: "bin/.helmrepo",
		Debug:            true,
		Linting:          false,
	})
	if err != nil {
		log.Err(err, "create helm client error")
		return err
	}

	failedUpdate := []string{}

	for _, app := range c.Apps {
		version := app.ContainerVersionPrefix + GetRemoteVersion(app, helmClient)
		if version == "" {
			log.Warn("No valid version '"+version+"' for", app.Name)
		}
		replaced, err := UpdateVersion(app, version, c.BaseFolder)
		if err != nil {
			log.Err(err, "update version error")
			failedUpdate = append(failedUpdate, app.Name)
			continue
		}
		if replaced {
			err = GitPush(app, version, c.BaseFolder, gitFlags)
			if err != nil {
				log.Err(err, "git push error")
				failedUpdate = append(failedUpdate, app.Name)
			} else {
				notify.SendNotification(c.Notify, app, version)
			}
		} else {
			log.Info("Already up to date.")
		}
	}

	if len(failedUpdate) > 0 {
		log.Error("Failed updates:")
		for _, item := range failedUpdate {
			log.Info("-", item)
		}
	}

	return nil
}

func UpdateVersion(app models.App, version, baseFolder string) (replaced bool, err error) {
	filepath := path.Join(baseFolder, app.Path)
	log.Info("Update app '"+app.Name+"' with version '"+version+"' file", filepath)
	pattern := GetVersionReplacePattern(app)
	return manifest.SetTagInFileWithPattern(version, "", filepath, pattern)
}

func GitPush(app models.App, version, baseFolder string, gitFlags util.GitFlags) error {
	msg := "Deploy " + app.Name + " version: " + version
	return util.GitPush(baseFolder, msg, gitFlags)
}

func GetVersionReplacePattern(app models.App) string {
	if app.System == models.Kustomize {
		return `([ ]+version: )([a-zA-Z0-9_.-].*)`
	}
	if app.System == models.ArgoCd {
		return `([ ]+targetRevision: )([a-zA-Z0-9_.-].*)`
	}
	if app.System == models.Deployment {
		return `([ ]+image: ` + app.ContainerImage + `:)([a-zA-Z0-9_.-].*)`
	}
	return ""
}

// config

func LoadConfig(filepath string) *models.Config {
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Err(err, "load config error.")
		return nil
	}

	var config models.Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Err(err, "load config error.")
		return nil
	}
	return &config
}
