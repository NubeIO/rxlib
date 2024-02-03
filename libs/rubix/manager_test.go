package rubix

import (
	"fmt"
	"github.com/NubeIO/rxlib/helpers/pprint"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()

	// Create a new network
	networkParams := AddNetwork{
		Name:             "MyNetwork",
		ObjectUUID:       "NetworkUUID",
		ParentObjectUUID: "ParentUUID",
	}
	network := NewNetwork(&networkParams)

	// Add the network to the manager
	manager.AddNetwork(network)

	// Create a new mapping
	mappingParams := &AddMapping{
		UUID:                 "MappingUUID",
		SourceUUID:           "SourceUUID",
		SourceConnectionUUID: "SourceConnectionUUID",
		TargetUUID:           "TargetUUID",
		TargetConnectionUUID: "TargetConnectionUUID",
	}
	mapping := NewMapping(mappingParams)

	// Add the mapping to the network
	network.NewMapping(mapping)

	// Create a new connection
	connection := &Connection{
		UUID:           "ConnectionUUID",
		ConnectionUUID: "ConnectionUUID",
		ObjectUUID:     "ObjectUUID",
		PortID:         "PortID",
		Enable:         true,
		IsError:        false,
		Created:        time.Now(),
		LastOk:         nil,
		LastFail:       nil,
		FailCount:      0,
	}

	// Add the connection to the mapping
	err := network.AddConnection("MappingUUID", connection)
	if err != nil {
		fmt.Println("Error adding connection:", err)
	}

	// Edit the connection
	editedConnection := &Connection{
		UUID:           "ConnectionUUID",
		ConnectionUUID: "ConnectionUUID",
		ObjectUUID:     "ObjectUUID",
		PortID:         "EditedPortID",
		Enable:         true,
		IsError:        false,
		Created:        time.Now(),
		LastOk:         nil,
		LastFail:       nil,
		FailCount:      0,
	}

	err = network.EditConnection("MappingUUID", editedConnection)
	if err != nil {
		fmt.Println("Error editing connection:", err)
	}

	// Delete the connection
	err = network.DeleteConnection("MappingUUID", "ConnectionUUID")
	if err != nil {
		fmt.Println("Error deleting connection:", err)
	}

	// Retrieve all mappings for the network
	mappings := network.AllMapping()
	fmt.Println("Mappings for the network:")
	for _, m := range mappings {
		fmt.Printf("Mapping UUID: %s\n", m.UUID)
	}
	pprint.PrintJSON(manager.All())
}
