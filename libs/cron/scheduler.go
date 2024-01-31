package cron

import (
	"fmt"
	"github.com/NubeIO/rubix-rx/helpers"
	"github.com/go-co-op/gocron"
	"sync"
	"time"
)

type Schedulers interface {
	StartAsync()
	Stop()
	Add(name string) (uuid string)
	Get(uuid string) *Job
	Delete(uuid string)
}

// Scheduler represents a generic scheduler.
type Scheduler interface {
	ScheduleTaskCron(name, cron string, startCallback, completedCallback func(*JobInfo), delay time.Duration) (*JobInfo, error)
	ScheduleTaskWithInterval(name string, interval any, timeToTrigger int, startCallback, completedCallback func(*JobInfo), delay time.Duration) (*JobInfo, error)
	CancelTask() error
	GetUUID() string
	GetInfo() *JobInfo
}

type schedulers struct {
	scheduler []*Job
	mu        sync.Mutex
	cron      *gocron.Scheduler
}

func NewScheduler(cron *gocron.Scheduler) Schedulers {
	return &schedulers{
		cron:      cron,
		scheduler: []*Job{},
	}
}

func (inst *schedulers) StartAsync() {
	inst.cron.StartAsync()
}

func (inst *schedulers) Stop() {
	inst.cron.Stop()
}

type Job struct {
	Name     string
	uuid     string
	location string
	cron     *gocron.Scheduler
	job      *JobInfo
}

type Task struct {
	Name              string
	StartCallback     func(*JobInfo)
	CompletedCallback func(*JobInfo)
	Delay             time.Duration
	UUID              string
	jobInfo           *JobInfo // Store job info here

}

func (inst *schedulers) Add(name string) (uuid string) {
	inst.mu.Lock()
	defer inst.mu.Unlock()
	j := &Job{
		uuid: helpers.UUID(),
		Name: name,
	}
	j.cron = inst.cron
	inst.scheduler = append(inst.scheduler, j)
	return j.uuid
}

func (inst *schedulers) Get(uuid string) *Job {
	inst.mu.Lock()
	defer inst.mu.Unlock()
	for _, sch := range inst.scheduler {
		if sch.GetUUID() == uuid {
			return sch
		}
	}
	return nil
}

func (inst *schedulers) Delete(uuid string) {
	inst.mu.Lock()
	defer inst.mu.Unlock()

	for i, sch := range inst.scheduler {
		if sch.GetInfo().UUID == uuid {
			inst.scheduler = append(inst.scheduler[:i], inst.scheduler[i+1:]...)
			err := sch.CancelTask()
			if err != nil {
				fmt.Println(err)
				return
			}
		}

	}
}

type JobInfo struct {
	UUID             string    `json:"uuid"`
	Name             string    `json:"name"`
	NextRun          time.Time `json:"nextRun"`
	RunCount         int       `json:"runCount"`
	IsRunning        bool      `json:"isRunning"`
	ScheduledAtTimes []string  `json:"scheduledAtTimes"`
	LastRun          time.Time `json:"lastRun"`
}

func (inst *Job) ScheduleTaskCron(name, cron string, startCallback, completedCallback func(*JobInfo), delay time.Duration) (*JobInfo, error) {
	newUUID := helpers.UUID()
	task := &Task{
		Name:              name,
		StartCallback:     startCallback,
		CompletedCallback: completedCallback,
		Delay:             delay,
		UUID:              newUUID,
	}
	inst.uuid = newUUID

	taskFunc := func() {
		if task.StartCallback != nil {
			task.StartCallback(task.jobInfo) // Pass the JobInfo to the callback
		}

		<-time.After(task.Delay)

		if task.CompletedCallback != nil {
			task.CompletedCallback(task.jobInfo) // Pass the JobInfo to the callback
		}
	}

	job, err := inst.cron.Tag(task.UUID).Cron(cron).Do(taskFunc)
	if err != nil {
		return nil, err
	}
	job.ScheduledAtTimes()
	jobInfo := &JobInfo{
		UUID:             task.UUID,
		Name:             name,
		NextRun:          job.NextRun(),
		RunCount:         job.RunCount(),
		IsRunning:        job.IsRunning(),
		ScheduledAtTimes: job.ScheduledAtTimes(),
		LastRun:          job.LastRun(),
	}
	task.jobInfo = jobInfo // Store job info in Task
	return jobInfo, nil
}
func (inst *Job) ScheduleTaskWithInterval(name string, interval any, timeToTrigger int, startCallback, completedCallback func(*JobInfo), delay time.Duration) (*JobInfo, error) {
	newUUID := helpers.UUID()
	task := &Task{
		Name:              name,
		StartCallback:     startCallback,
		CompletedCallback: completedCallback,
		Delay:             delay,
		UUID:              newUUID,
	}
	inst.uuid = newUUID

	// Create a channel to signal the completion of each job run
	doneChan := make(chan struct{}, timeToTrigger)
	job, err := inst.cron.Tag(task.UUID).Every(interval).Do(func() {
		if task.StartCallback != nil {
			task.StartCallback(inst.job)
		}

		<-time.After(task.Delay)

		if task.CompletedCallback != nil {
			task.CompletedCallback(inst.job)
		}

		// Signal completion of a job run
		doneChan <- struct{}{}
	})

	if err != nil {
		return nil, err
	}

	go func() {
		for i := 0; i < timeToTrigger; i++ {
			<-doneChan
			fmt.Printf("Job completed %d time(cron)\n", i+1)
		}
		fmt.Println("Cancelling the task")
		err := inst.CancelTask()
		if err != nil {
			fmt.Println("Cancelling the task", err)
			return
		}
	}()

	inst.job = &JobInfo{
		Name:             name,
		UUID:             task.UUID,
		NextRun:          job.NextRun(),
		RunCount:         job.RunCount(),
		IsRunning:        job.IsRunning(),
		ScheduledAtTimes: job.ScheduledAtTimes(),
		LastRun:          job.LastRun(),
	}
	return inst.job, nil
}

func (inst *Job) CancelTask() error {
	err := inst.cron.RemoveByTag(inst.uuid)
	return err
}

func (inst *Job) GetUUID() string {
	return inst.uuid
}

func (inst *Job) GetInfo() *JobInfo {
	return inst.job
}
