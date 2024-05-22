package main

import (
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"log"
	"plugin"
)

func main() {
	pluginName := "rabbitmq_cdc_connector"
	// Load the plugin
	// 1. Search the plugins directory for a file with the same name as the pluginName
	// that was passed in as an argument and attempt to load the shared object file.
	plug, err := plugin.Open(fmt.Sprintf("/Users/antonioradesca/Code/cdc-microservices/microservices/plugins/bins/%s.so", pluginName))
	if err != nil {
		log.Fatal(err)
	}
	// 2. Look for an exported symbol such as a function or variable
	// in our case we expect that every plugin will have exported a single struct
	// that implements the Shipper interface with the name "Shipper"
	connectorSample, err := plug.Lookup("Connector")
	if err != nil {
		log.Fatal(err)
	}
	// 3. Attempt to cast the symbol to the Shipper
	// this will allow us to call the methods on the plugins if the plugin
	// implemented the required methods or fail if it does not implement it.
	var connector cdc_shared.ConnectorProvider
	connector, ok := connectorSample.(cdc_shared.ConnectorProvider)
	if !ok {
		log.Fatal(ok)
	}
	fmt.Println(connector.Name())
}
