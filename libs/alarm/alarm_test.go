package alarm

import (
	"fmt"
	pprint "github.com/NubeIO/rubix-rx/helpers/print"
	"testing"
)

func TestNewAlarmManager(t *testing.T) {
	alarmManager := NewAlarmManager()

	// Create a new Alarm with a sample limit size
	alarm1 := alarmManager.NewAlarm(10)

	// Create and add transactions to the Alarm
	transaction1 := NewTransaction()
	transaction2 := NewTransaction()

	alarm1.AddTransaction(transaction1)
	alarm1.AddTransaction(transaction2)

	// Retrieve and display all Alarms
	alarms := alarmManager.All()

	for s, transactions := range alarmManager.AllTransactions() {
		for i, transaction := range transactions {
			fmt.Println(s, i, transaction.GetUUID())

		}
		//fmt.Println(s, transactions)

	}

	for _, a := range alarms {
		fmt.Printf("Alarm UUID: %s\n", a.GetUUID())

		// Retrieve and display transactions for each Alarm
		transactions := a.GetTransactions()
		for _, t := range transactions {
			fmt.Printf("  Transaction UUID: %s, CreatedAt: %s  %s\n", t.GetUUID(), t.GetCreatedAt(), t.GetAlarmUUID())
		}
	}

	pprint.PrintJOSN(alarmManager.GetAllTransactionsEntries())
}
