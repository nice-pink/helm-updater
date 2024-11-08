package helmupdater

import (
	"encoding/json"
	"os"
	"path"

	helmclient "github.com/mittwald/go-helm-client"
	"github.com/nice-pink/goutil/pkg/log"
	"github.com/nice-pink/helm-auto-updater/pkg/models"
	"github.com/nice-pink/repo-services/pkg/manifest"
)

func Run(configFile string) error {
	c := LoadConfig(configFile)

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

	for _, app := range c.Apps {
		version := GetRemoteVersion(app, helmClient)
		if err := UpdateVersion(app, version, c.BaseFolder); err != nil {
			log.Err(err, "update version error")
		}
	}

	return nil
}

func UpdateVersion(app models.App, version, baseFolder string) error {
	filepath := path.Join(baseFolder, app.Path)
	log.Info("Update app '"+app.Name+"' with version '"+version+"' file", filepath)
	pattern := GetVersionReplacePattern(app)
	_, err := manifest.SetTagInFileWithPattern(version, "", filepath, pattern)
	return err
}

func GetVersionReplacePattern(app models.App) string {
	if app.System == models.Kustomize {
		return `([ ]+version: )([a-zA-Z0-9_.-].*)`
	}
	if app.System == models.ArgoCd {
		return `([ ]+targetRevision: )([a-zA-Z0-9_.-].*)`
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
