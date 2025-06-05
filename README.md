![Logo of helm-updater showing a nautic stearing wheel and and arrow towards the outlines of the kubernetes logo.](https://www.nice.pink/img/helm-updater.png)

# What

*helm-updater* is a tool to auto update or notify on updates of helm apps defined by manifests in gitops repository.

# Manifests

In the current version the helm-updater can update manifests of type:

- kustomize
- k8s (deployments, statefulsets, daemonsets, argo-rollouts, ...)
- argo-cd
- terraform

# Config

An example config file can be found in [cmd/helm-updater/config.yaml](cmd/helm-updater/config.yaml).

Fields for apps:

- *autoUpdate*: Should auto update and git push if finds a new version.
- *name*: Of of app in helm repo
- *repo*: Url of helm repo
- *private*: Is private helm repo?
- *version*: Version matcher for auto updating. `"*"`: all, `"1.*.*"`: fix major version, ...
- *system*: Type of manifest in which app is defined.
- *paths*: Manifest paths in git repo.
- *containerImage*: If system is `k8s` the full image needs to be set (without version) to match the correct container.
- *containerVersionPrefix*: Container images might contain a prefix like "v"
- *repoUsername*: For private helm repos.
- *repoPassword*: For private helm repos.

# Slack notification

Slack notifications can either be sent via slack webhook OR slack oauth token to channel id.

1. To use webhook the env var *HELM_UPDATER_NOTIFY_WEBHOOK* can be set. Alternatively set *webhook* in config.
2. To use slack app with oauth token, set env var *SLACK_TOKEN* and define channelId in config. *token* can also be set in config.

## Helm

### Private repos

Use credentials:

Either add to config file OR use env vars. If both are specified env vars will be preferrd.

Add a general username and password used as default for private helm repos: `HELM_UPDATER_PRIVATE_REPO_USERNAME`, `HELM_UPDATER_PRIVATE_REPO_PASSWORD`

Use specific helm repo username and password: `HELM_UPDATER_${REPO_NAME}_USERNAME`, `HELM_UPDATER_${REPO_NAME}_PASSWORD`

# Build command line executables

## Build single executable

1. Add executables as sub-folder into `cmd` folder. E.g. `cmd/exec`
2. Open terminal and `cd` to base folder of this repo.
3. Type `./build NAME_OF_EXECUTABLE`. E.g. `./build exec`
4. Executable will be created in `bin/NAME_OF_EXECUTABLE`. E.g. `bin/exec`
5. Run executable. E.g. `bin/exec`

## Build all

1. Add executables as sub-folder into `cmd` folder. E.g. `cmd/exec`
2. Open terminal and `cd` to base folder of this repo.
3. Type `./build`.
4. All executables will be created in `bin` folder.
