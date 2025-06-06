package helmupdater

import (
	helmclient "github.com/mittwald/go-helm-client"
	"github.com/nice-pink/goutil/pkg/data"
	"github.com/nice-pink/goutil/pkg/log"
	"github.com/nice-pink/helm-updater/pkg/models"
	"github.com/nice-pink/helm-updater/pkg/notify"
	"github.com/nice-pink/repo-services/pkg/util"
)

type Updater struct {
	notifyClient notify.Client
	gitFlags     util.GitFlags
	config       *models.Config
	helmClient   helmclient.Client
}

func NewUpdater(configFile string, gitFlags util.GitFlags) *Updater {
	c := LoadConfig(configFile)
	if c == nil {
		return nil
	}

	// init notify client
	notifyClient := notify.NewClient(c.Notify)

	// the client is only used to pull repos so most options don't really matter
	helmClient, err := helmclient.New(&helmclient.Options{
		Namespace:        "default",
		RepositoryCache:  c.Helm.CachePath,
		RepositoryConfig: c.Helm.RepoFilePath,
		Debug:            false,
		Linting:          false,
	})
	if err != nil {
		log.Err(err, "create helm client error")
		return nil
	}

	return &Updater{
		config:       c,
		notifyClient: notifyClient,
		gitFlags:     gitFlags,
		helmClient:   helmClient,
	}
}

func (u *Updater) Run() error {
	// checkout repo?
	if u.gitFlags.Url != nil {
		err := util.GitClone(*u.gitFlags.Url, u.config.BaseFolder, u.gitFlags)
		if err != nil {
			return err
		}
	}

	failedUpdate := []string{}

	for _, app := range u.config.Apps {
		log.Info("---", app.Name, "-", app.Repo)
		// version := app.ContainerVersionPrefix + GetRemoteVersion(app, helmClient) // container version prefix is added in other part
		version := u.getRemoteVersion(app)
		if version == "" {
			log.Warn("No valid version '"+version+"' for", app.Name)
		}
		err := u.updateVersion(app, version, u.config.BaseFolder)
		if err != nil {
			failedUpdate = append(failedUpdate, app.Name)
		}

		// clean up helm cache
		if u.config.Helm.CleanUp {
			ClearHelmCache(u.config.Helm.CachePath, u.config.Helm.RepoFilePath)
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

func (u *Updater) updateVersion(app models.App, version, baseFolder string) error {
	for i, path := range app.Paths {
		replaced, newAvailable, err := updateVersionInPath(app, path, version, baseFolder)
		if err != nil {
			log.Err(err, "update version error")
			return err
		}
		if replaced {
			err = GitPush(app, version, baseFolder, u.gitFlags)
			if err != nil {
				return err
			} else {
				u.sendNotification(app, i, version, true)
			}
		} else if newAvailable {
			u.sendNotification(app, i, version, false)
		} else {
			log.Info("Already up to date.")
		}
	}
	return nil
}

func (u *Updater) sendNotification(app models.App, index int, version string, updated bool) {
	log.Info("send")
	// if u.config == nil {
	// 	return
	// }
	// u.notifyClient.SendNotification(u.config.Notify, app, index, version, false)
}

//

func updateVersionInPath(app models.App, appPath, version, baseFolder string) (replaced, newAvailable bool, err error) {
	log.Info("Check", app.Name, "in", appPath)
	// get manifest data and path
	manifest, path, err := getManifest(appPath, baseFolder)
	if err != nil {
		log.Err(err, "open manifest")
		return false, false, err
	}

	// new version available
	current := getCurrentVersion(app, manifest)
	newAvailable = current != version
	if newAvailable {
		log.Info("New version available. Current:", current, "New:", version)
	}

	// if should not auto update return
	if !newAvailable || !app.AutoUpdate {
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
