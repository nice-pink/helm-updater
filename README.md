# Config

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
