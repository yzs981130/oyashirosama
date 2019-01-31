package scheduler

import "errors"

type void struct{}
var Exists void

func (js *JobStatus) Set(name string, status string) error {
	js.CheckContains(name)
	js.Information[name] = status
	switch status {
	case "RunningStatus":
		js.RunningStatus[name] = Exists
	case "WaitingStatus":
		js.WaitingStatus[name] = Exists
	case "SchedulingStatus":
		js.WaitingStatus[name] = Exists
	case "FinishedStatus":
		js.FinishedStatus[name] = Exists
	case "TerminatedStatus":
		js.TerminatedStatus[name] = Exists
	default:
		err := errors.New("unknown status")
		log.Error(err, "unknown status when setting", "status", status)
		return err
	}
	return nil
}

func (js *JobStatus) CheckContains(name string) bool {
	val, ok := js.Information[name]
	if ok {
		log.Info("found exist job status pair", name, val)
	}
	return ok
}



type JobStatus struct {
	RunningStatus map[string]void
	WaitingStatus map[string]void
	SchedulingStatus map[string]void
	FinishedStatus map[string]void
	TerminatedStatus map[string]void

	Information map[string]string
}
