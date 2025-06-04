package helmupdater

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/nice-pink/goutil/pkg/log"
	"github.com/nice-pink/helm-updater/pkg/models"
	"github.com/nice-pink/helm-updater/pkg/utils"
)

func Update(app models.App, version, outputFilepath, manifest string) (string, bool, error) {
	replacePattern := getVersionPattern(app)
	log.Info("Version Pattern:", replacePattern)

	// Create the replacement string with the version
	fullVersion := app.ContainerVersionPrefix + version
	replacement := getVersionReplacePattern(app, fullVersion)
	log.Info("Replacement:", replacement)

	// Perform the replacement
	newContent, err := utils.Replace(replacePattern, replacement, manifest)
	if err != nil {
		log.Err(err, "replace error")
		return "", false, err
	}

	// Check if the content actually changed
	if newContent == manifest {
		log.Info("same content")
		return "", false, nil
	}

	// Write the updated content back to the file
	if outputFilepath != "" {
		err = os.WriteFile(outputFilepath, []byte(newContent), 0644)
		if err != nil {
			return "", false, err
		}
	}

	return newContent, true, nil
}

func getVersionPattern(app models.App) string {
	if app.System == models.Kustomize {
		return `([ ]+version: )([a-zA-Z0-9_.-].*)`
	}
	if app.System == models.ArgoCd {
		return `([ ]+targetRevision: )([a-zA-Z0-9_.-].*)`
	}
	if app.System == models.Terraform {
		return `([ ]+version[ ]+=[ ]+` + `")([a-zA-Z0-9_.-].*)(")`
	}
	if app.System == models.Deployment {
		//image: quay.io/oauth2-proxy/oauth2-proxy:v7.9.0
		return `([ ]+image: ` + app.ContainerImage + `:)([a-zA-Z0-9_.-].*)`
	}
	return ""
}

func getVersionReplacePattern(app models.App, version string) string {
	if app.System == models.Kustomize {
		return "${1}" + version
	}
	if app.System == models.ArgoCd {
		return "${1}" + version
	}
	if app.System == models.Terraform {
		return "${1}" + version + "${3}"
	}
	if app.System == models.Deployment {
		return "${1}" + version
	}
	return ""
}

func getVersionIsolatePattern() string {
	return "${2}"
}

func getCurrentVersion(app models.App, manifest string) string {
	versionPattern := getVersionPattern(app)
	isolatePattern := getVersionIsolatePattern()
	v := utils.Find(versionPattern, isolatePattern, manifest)
	return strings.TrimPrefix(v, app.ContainerVersionPrefix)
}

func updateVersionSegment(app models.App, version string, manifest string) (string, bool) {
	versionCurrent := getCurrentVersion(app, manifest)
	versionCurrent = strings.TrimPrefix(versionCurrent, app.ContainerVersionPrefix)

	// is the same version?
	if versionCurrent == version {
		return version, false
	}

	// version matcher
	if app.Version == "" || app.Version == "*" || app.Version == "*.*" || app.Version == "*.*.*" {
		return version, true
	}

	// construct version
	match := strings.Split(app.Version, ".")
	vC := strings.Split(versionCurrent, ".")
	v := strings.Split(version, ".")
	if len(v) <= 1 || len(vC) != len(v) {
		// set the new version,
		// - if version has max len 1 then
		// OR
		// - if the format of the current version and new version is different
		return version, true
	}

	newVersion := ""
	for i, seg := range v {
		if i >= len(match) {
			newVersion += getVersionSegment(seg, i)
			continue
		}

		if match[i] == "*" {
			newVersion += getVersionSegment(seg, i)
		} else {
			newVersion += getVersionSegment(match[i], i)
		}
	}
	return newVersion, true
}

func getVersionSegment(seg string, i int) string {
	if i > 0 {
		return "." + seg
	}
	return seg
}

// data

func getManifest(app models.App, baseFolder string) (content string, path string, err error) {
	// Construct the full path to the manifest file
	fullPath := filepath.Join(baseFolder, app.Path)

	// Read the manifest file
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", "", err
	}
	return string(data), fullPath, err
}
