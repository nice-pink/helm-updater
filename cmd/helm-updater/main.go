package main

import (
	"flag"
	"os"

	"github.com/nice-pink/goutil/pkg/log"
	"github.com/nice-pink/helm-updater/pkg/helmupdater"
	"github.com/nice-pink/repo-services/pkg/util"
)

var (
	config = flag.String("config", "", "Path to config file.")
)

func main() {
	testFail := flag.Bool("testFail", false, "Test failing.")
	gitFlags := util.GetGitFlags()
	flag.Parse()

	log.Info("*** Start")
	log.Info(os.Args)

	if *testFail {
		log.Error("TEST FAILING!!!")
		os.Exit(2)
	}

	err := helmupdater.Run(*config, gitFlags)
	if err != nil {
		os.Exit(2)
	}
}
