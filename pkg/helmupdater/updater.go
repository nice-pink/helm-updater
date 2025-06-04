package helmupdater

import (
	helmclient "github.com/mittwald/go-helm-client"
	"github.com/nice-pink/goutil/pkg/data"
	"github.com/nice-pink/goutil/pkg/log"
	"github.com/nice-pink/helm-updater/pkg/models"
	"github.com/nice-pink/helm-updater/pkg/notify"
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
		replaced, newAvailable, err := UpdateVersion(app, version, c.BaseFolder)
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
				notify.SendNotification(c.Notify, app, version, true)
			}
		} else if newAvailable {
			notify.SendNotification(c.Notify, app, version, false)
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

func UpdateVersion(app models.App, version, baseFolder string) (replaced, newAvailable bool, err error) {
	// get manifest data and path
	manifest, path, err := getManifest(app, baseFolder)
	if err != nil {
		log.Err(err, "open manifest")
		return false, false, err
	}

	// new version available
	current := getCurrentVersion(app, manifest)
	newAvailable = current != version

	// if should not auto update return
	if !app.AutoUpdate {
		return false, newAvailable, nil
	}

	// update manifest
	_, replaced, err = Update(app, version, path, manifest)
	return replaced, newAvailable, err
}

func GitPush(app models.App, version, baseFolder string, gitFlags util.GitFlags) error {
	msg := "Deploy " + app.Name + " version: " + version
	return util.GitPush(baseFolder, msg, gitFlags)
}

// config

func LoadConfig(filepath string) *models.Config {
	log.Info("Load config from", filepath)

	var config models.Config
	err := data.ReadJsonOrYaml(filepath, &config)
	if err != nil {
		log.Err(err, "load config error.", filepath)
		return nil
	}
	return &config
}
