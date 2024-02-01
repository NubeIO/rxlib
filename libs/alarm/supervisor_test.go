package alarm

import (
	"fmt"
	"github.com/NubeIO/rxlib/helpers/pprint"

	"testing"
)

func TestSupervisorWithAlarmsAndTransactions(t *testing.T) {
	// Create a new supervisor
	supervisor := NewSupervisor()

	// Create and add AlarmManager instances for two remote servers
	server1Manager := NewAlarmManager("bac 33")
	server2Manager := NewAlarmManager("bac 33")

	supervisor.AddRemoteManager("Server1", server1Manager)
	supervisor.AddRemoteManager("Server2", server2Manager)

	// Create and add alarms with transactions to the managers
	createAlarmWithTransactions(server1Manager, "Server1_Alarm1")
	createAlarmWithTransactions(server1Manager, "Server1_Alarm2")
	createAlarmWithTransactions(server2Manager, "Server2_Alarm1")

	// Retrieve AlarmManager instances from the supervisor
	manager1 := supervisor.GetRemoteManager("Server1")
	manager2 := supervisor.GetRemoteManager("Server2")

	// Verify that the retrieved managers are not nil
	if manager1 == nil {
		t.Errorf("Failed to retrieve AlarmManager for Server1")
	}

	if manager2 == nil {
		t.Errorf("Failed to retrieve AlarmManager for Server2")
	}

	// Check the count of managed servers
	allManagers := supervisor.All()
	if len(allManagers) != 2 {
		t.Errorf("Expected 2 managed servers, but got %d", len(allManagers))
	}

	// Print the list of managed servers and their corresponding AlarmManager instances
	fmt.Println("Managed Servers:")
	for serverName, manager := range allManagers {
		fmt.Printf("Server: %s\n", serverName)
		fmt.Printf("AlarmManager: %v\n", manager)
	}

	pprint.PrintJSON(supervisor.GetAllManagersTransactionsEntries())

}

func createAlarmWithTransactions(manager AlarmManager, alarmName string) {
	alarmBody := &AddAlarm{
		Title:      "t",
		ObjectType: "device",
		ObjectUUID: "12345",
	}
	alarm := manager.NewAlarm(10, alarmBody) // Create a new alarm with a transaction limit of 10

	// Add transactions to the alarm
	for i := 1; i <= 5; i++ {
		transaction := NewTransaction()
		alarm.AddTransaction(NewTransactionBody(AlarmStatusActive, AlarmSeverityCritical, "Transaction 1", "Transaction Body 1"), transaction)
	}

	// Add the alarm to the manager
	manager.Get(alarm.GetUUID())
}
