package alarm

import (
	"fmt"
	"github.com/NubeIO/rxlib/helpers/pprint"

	"testing"
)

func TestNewAlarmManager(t *testing.T) {
	// Create a new Manager
	alarmManager := NewAlarmManager("abc 123")
	alarmBody := &AddAlarm{
		Title:      "t",
		ObjectType: "device",
		ObjectUUID: "12345",
	}
	// Create two alarms with a sample limit size
	alarm1 := alarmManager.NewAlarm(10, alarmBody)
	alarm2 := alarmManager.NewAlarm(5, alarmBody)

	// Add transactions to alarm1
	transaction1 := NewTransaction()
	transaction2 := NewTransaction()
	transaction3 := NewTransaction()

	alarm1.AddTransaction(NewTransactionBody(StatusActive, SeverityCritical, "Transaction 1", "Transaction Body 1"), transaction1)
	alarm1.AddTransaction(NewTransactionBody(StatusActive, SeverityWarning, "Transaction 2", "Transaction Body 2"), transaction2)
	alarm1.AddTransaction(NewTransactionBody(StatusClosed, SeverityInfo, "Transaction 3", "Transaction Body 3"), transaction3)

	// Add transactions to alarm2
	transaction4 := NewTransaction()
	transaction5 := NewTransaction()

	alarm2.AddTransaction(NewTransactionBody(StatusActive, SeverityCritical, "Transaction 4", "Transaction Body 4"), transaction4)
	alarm2.AddTransaction(NewTransactionBody(StatusAcknowledged, SeverityError, "Transaction 5", "Transaction Body 5"), transaction5)

	// Print all transactions as TransactionEntry
	transactionsEntries := alarmManager.GetAllTransactionsEntries()
	for alarmUUID, transactionEntries := range transactionsEntries {
		fmt.Printf("Alarm UUID: %s\n", alarmUUID)
		for _, transactionEntry := range transactionEntries {
			fmt.Printf("  Transaction UUID: %s\n", transactionEntry.UUID)
			fmt.Printf("  Status: %s\n", transactionEntry.Status)
			fmt.Printf("  Severity: %s\n", transactionEntry.Severity)
			fmt.Printf("  Target: %s\n", transactionEntry.Target)
			fmt.Printf("  Title: %s\n", transactionEntry.Title)
			fmt.Printf("  Body: %s\n", transactionEntry.Body)
			fmt.Printf("  Created At: %s\n", transactionEntry.CreatedAt)
			fmt.Printf("  Last Updated: %s\n", transactionEntry.LastUpdated)
			fmt.Println("--------------------")
		}
	}

	// Retrieve all transactions
	allTransactions := alarmManager.AllTransactions()
	fmt.Printf("All Transactions:\n%v\n", allTransactions)

	// Drop alarm2
	alarmManager.Drop(alarm2.GetUUID())

	// Retrieve all alarms
	allAlarms := alarmManager.All()
	fmt.Printf("All Alarms:\n%v\n", allAlarms)

	// Delete transactions by UUID
	deleteUUIDs := map[string]string{
		transaction1.GetUUID(): "alarmUUID1",
		transaction2.GetUUID(): "alarmUUID1",
		transaction3.GetUUID(): "alarmUUID1",
	}
	alarmManager.DeleteTransactions(deleteUUIDs)

	// Print transactions after deletion
	allTransactionsAfterDeletion := alarmManager.AllTransactions()
	fmt.Printf("All Transactions After Deletion:\n%v\n", allTransactionsAfterDeletion)

	pprint.PrintJSON(alarmManager.GetAllTransactionsEntries())
}
