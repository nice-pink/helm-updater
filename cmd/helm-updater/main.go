package main

import (
	"flag"
	"os"

	"github.com/nice-pink/goutil/pkg/log"
	"github.com/nice-pink/helm-auto-updater/pkg/helmupdater"
)

var (
	config = flag.String("config", "", "Path to config file.")
)

func main() {
	flag.Parse()

	log.Info("*** Start")
	log.Info(os.Args)

	helmupdater.Run(*config)
}
