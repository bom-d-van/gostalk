package gostalk

import (
	"time"
)

type delayedJobs struct {
	amount int
}

func newDelayedJobs() (jobs *delayedJobs) {
	return &delayedJobs{
		amount: 0,
	}
}

func (jobs *delayedJobs) Len() int {
	return jobs.amount
}

func (jobs *delayedJobs) putJob(job *job, afterFunc func()) {
	duration := job.delayEndsAt.Sub(time.Now())
	job.jobHolder = jobs
	job.timer = time.AfterFunc(duration, afterFunc)
}

func (jobs *delayedJobs) buryJob(job *job) {
	job.timer.Stop()
	job.tube.buried.putJob(job)
}

func (jobs *delayedJobs) deleteJob(job *job) {
	job.timer.Stop()
}

func (jobs *delayedJobs) touchJob(job *job) {}

// TODO: implement!
func (jobs *delayedJobs) kickJobs(bound int) (actual int) {
	return
}

// TODO: implement!
func (jobs *delayedJobs) peekJob(request *jobPeekRequest) {
	request.success <- nil
}
