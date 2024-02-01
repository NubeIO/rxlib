package alarm

import (
	"fmt"
	pprint "github.com/NubeIO/rubix-rx/helpers/print"
	"testing"
)

func TestCoordinatorWithAlarmsAndTransactions(t *testing.T) {
	// Create a new coordinator
	coordinator := NewCoordinator()

	// Create and add AlarmManager instances for two remote servers
	server1Manager := NewAlarmManager("bac 33")
	server2Manager := NewAlarmManager("bac 33")

	coordinator.AddRemoteAlarmManager("Server1", server1Manager)
	coordinator.AddRemoteAlarmManager("Server2", server2Manager)

	// Create and add alarms with transactions to the managers
	createAlarmWithTransactions(server1Manager, "Server1_Alarm1")
	createAlarmWithTransactions(server1Manager, "Server1_Alarm2")
	createAlarmWithTransactions(server2Manager, "Server2_Alarm1")

	// Retrieve AlarmManager instances from the coordinator
	manager1 := coordinator.GetRemoteAlarmManager("Server1")
	manager2 := coordinator.GetRemoteAlarmManager("Server2")

	// Verify that the retrieved managers are not nil
	if manager1 == nil {
		t.Errorf("Failed to retrieve AlarmManager for Server1")
	}

	if manager2 == nil {
		t.Errorf("Failed to retrieve AlarmManager for Server2")
	}

	// Check the count of managed servers
	allManagers := coordinator.All()
	if len(allManagers) != 2 {
		t.Errorf("Expected 2 managed servers, but got %d", len(allManagers))
	}

	// Print the list of managed servers and their corresponding AlarmManager instances
	fmt.Println("Managed Servers:")
	for serverName, manager := range allManagers {
		fmt.Printf("Server: %s\n", serverName)
		fmt.Printf("AlarmManager: %v\n", manager)
	}

	pprint.PrintJOSN(coordinator.GetAllManagersTransactionsEntries())

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
