package cron

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"testing"
	"time"
)

func TestNewScheduler(t *testing.T) {
	schedulersManager := NewScheduler(gocron.NewScheduler(time.Now().Location()))

	schedulersManager.StartAsync() // Initialize the scheduler
	newJob := schedulersManager.Add("newScheduler")

	// Define the start and completed callbacks
	startCallback := func(j *JobInfo) {
		fmt.Println("Start Callback:", time.Now().Format(time.RFC3339), "UUID", j.UUID, "next run", j.NextRun)
	}

	completedCallback := func(j *JobInfo) {
		fmt.Println("Completed Callback:", time.Now().Format(time.RFC3339), "UUID", j.UUID, "next run", j.NextRun)
	}

	getJob := schedulersManager.Get(newJob)

	//cronExpression := crons.Schedule().Minutes(1)
	// Schedule a task to run every minute with start and completed callbacks
	//j, err := newScheduler.ScheduleTaskCron(cronExpression, startCallback, completedCallback, 2*time.Second)
	j, err := getJob.ScheduleTaskWithInterval("hey", "10s", 3, startCallback, completedCallback, 2*time.Second)
	if err != nil {
		fmt.Println("Error scheduling task:", err)
		return
	}

	fmt.Println("Task scheduled with UUID:", j.NextRun)

	// Wait for some time to allow tasks to run
	time.Sleep(6 * time.Second)
	//j.Job.
	err = getJob.CancelTask()
	fmt.Println(err)
	if err != nil {
		return
	}
	time.Sleep(2 * time.Minute)

	// Cancel the scheduled task by UUID
	//newScheduler.CancelTask()

	// Wait for a while to observe the cancellation
	time.Sleep(30 * time.Second)

	// Stop and remove the scheduler
	schedulersManager.Delete("uniqueUUID")

}
