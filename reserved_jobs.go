package gostalk

import (
  "code.google.com/p/go-priority-queue/prio"
)

type reservedJobsItem Job

func (i *reservedJobsItem) Less(j prio.Interface) bool {
  return i.reserveEndsAt.Before(j.(*reservedJobsItem).reserveEndsAt)
}

func (i *reservedJobsItem) Index(n int) {
  i.index = n
}

type reservedJobs struct {
  prio.Queue
}

func newReservedJobs() (jobs *reservedJobs) {
  return &reservedJobs{}
}

func (jobs *reservedJobs) getJob() (job *Job) {
  job = (*Job)(jobs.Pop().(*reservedJobsItem))
  job.jobHolder = nil
  return
}

func (jobs *reservedJobs) putJob(job *Job) {
  job.jobHolder = jobs
  jobs.Push((*reservedJobsItem)(job))
}

func (jobs *reservedJobs) deleteJob(job *Job) {
  jobs.Remove(job.index)
}
