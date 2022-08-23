package main

//go:generate sh injectGitVars.sh

import (
	"embed"

	log "github.com/s00500/env_logger"

	ib "github.com/SKAARHOJ/ibeam-corelib-go"
	pb "github.com/SKAARHOJ/ibeam-corelib-go/ibeam-core"
)

//go:embed model_images
var modelsFS embed.FS

/*

In the main file we read the config, initialize the corelib,
register models and parameters and finally start the implementation and manager routines

*/

func main() {
	ib.ReloadHook()
	ib.SetImageFS(&modelsFS) // make sure we register the devicecore images here

	branch := ""
	if gitBranch != "master" && gitBranch != "main" {
		branch = " branch: " + gitBranch
	}

	log.Infof("core-skaarhoj-template started, version %s (%s) %s", gitTag, gitRevision, branch)

	// Setup core info
	coreInfo := &pb.CoreInfo{
		CoreVersion:    gitTag,
		Description:    "ibeam template core implementation",
		Label:          "Core Template",
		DeviceCategory: pb.DeviceCategory_GenericProtocol,
		Name:           "core-skaarhoj-template",
		MaxDevices:     0, // Max devices allows to specify a maximum amount of devices the core can handle. keep 0 if this limit is not needed
		ConnectionType: pb.ConnectionType_Network,
	}

	// get the default config structure for the cores settings and devices
	config := defaultConfig()

	// Create the manager and registry
	manager, registry, toManager, fromManager := ib.CreateServerWithConfig(coreInfo, &config)

	// Alternatively do only create the components without a config
	// manager, registry, toManager, fromManager := ib.CreateServer(coreInfo)

	// You can start to register models here. A generic model that inherits ALL registered parameters will be created automatically.
	registry.RegisterModel(&pb.ModelInfo{
		Id:          1,
		Name:        "Model 1",
		Description: "A simple description for Model 1",
	})

	configureParameters(registry)

	go processDevices(registry, config, fromManager, toManager)

	// Finally start the server and manager
	manager.StartWithServer(":8502")
}
