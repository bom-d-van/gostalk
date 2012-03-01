package gostalker

/*
dsal | tube has jobs that are ready or not.  tube has channels to connections wanting jobs and new jobs coming in.
dsal | At any point there's either:  no jobs and no workers, no jobs and some workers, jobs and no workers or jobs and workers.  Just need to cover those cases and transitions between them.
*/

import (
  "container/heap"
)

type jobReserveRequest struct {
  success chan *Job
  cancel  chan bool
}

type Tube struct {
  name     string
  ready    *readyJobs
  reserved *reservedJobs
  buried   *buriedJobs
  delayed  *delayedJobs

  jobDemand chan *jobReserveRequest
  jobSupply chan *Job

  statUrgent   int
  statReady    int
  statReserved int
  statDelayed  int
  statBuried   int
}

func newTube(name string) (tube *Tube) {
  tube = &Tube{
    name:      name,
    ready:     newReadyJobs(),
    reserved:  newReservedJobs(),
    buried:    newBuriedJobs(),
    delayed:   newDelayedJobs(),
    jobDemand: make(chan *jobReserveRequest),
    jobSupply: make(chan *Job),
  }

  go tube.handleDemand()

  return
}

func (tube *Tube) handleDemand() {
  for {
    if tube.ready.Len() > 0 {
      select {
      case job := <-tube.jobSupply:
        tube.put(job)
      case request := <-tube.jobDemand:
        select {
        case request.success <- tube.reserve():
        case <-request.cancel:
          request.cancel <- true // propagate to the other tubes
        }
      }
    } else {
      tube.put(<-tube.jobSupply)
    }
  }
}

func (tube *Tube) reserve() (job *Job) {
  job = heap.Pop(tube.ready).(*Job)
  tube.statReady = tube.ready.Len()
  heap.Push(tube.reserved, job)
  tube.statReserved = tube.reserved.Len()
  return
}

func (tube *Tube) put(job *Job) {
  heap.Push(tube.ready, job)
  tube.statReady = tube.ready.Len()
  if job.priority < 1024 {
    tube.statUrgent += 1
  }
}
