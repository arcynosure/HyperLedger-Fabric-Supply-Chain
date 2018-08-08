package main

import (
	"fmt"
	"os"

	"github.com/servntire/car-ownership/blockchain"
	"github.com/servntire/car-ownership/web"
	"github.com/servntire/car-ownership/web/controllers"
)

func main() {
	// Definition of the Fabric SDK properties
	fSetup := blockchain.FabricSetup{
		// Channel parameters
		ChannelID:     "car-ownership",
		ChannelConfig: "" + os.Getenv("GOPATH") + "/src/github.com/servntire/car-ownership/fixtures/config/",

		// Chaincode parameters
		ChainCodeID:      "carownership-service",
		ChaincodeGoPath:  os.Getenv("GOPATH"),
		ChaincodePath:    "github.com/servntire/car-ownership/chaincode/",
		ChaincodeVersion: "1.17",
		OrgAdmin:         "Admin",
		OrgName:          "Org1",
		ConfigFile:       "config.yaml",
	}

	// Initialization of the Fabric SDK from the previously set properties
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v\n", err)
	}

	// Install and instantiate the chaincode
	err = fSetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
	}

	// Launch the web application listening
	app := &controllers.Application{
		Fabric: &fSetup,
	}
	web.Serve(app)

}
