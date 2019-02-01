package scheduler

import "errors"

//TODO: thread-safe

type void struct{}
var Exists = struct{}{}

func (js *JobStatus) init() *JobStatus {
	js.Information = make(map[string]string)
	js.RunningStatus = make(map[string]void)
	js.WaitingStatus = make(map[string]void)
	js.SchedulingStatus = make(map[string]void)
	js.FinishedStatus = make(map[string]void)
	js.TerminatedStatus = make(map[string]void)
	return js
}

func (js *JobStatus) Set(name string, status string) error {
	js.CheckContains(name)
	switch status {
	case "RunningStatus":
		js.RunningStatus[name] = Exists
		js.Information[name] = status
	case "WaitingStatus":
		js.WaitingStatus[name] = Exists
		js.Information[name] = status
	case "SchedulingStatus":
		js.WaitingStatus[name] = Exists
		js.Information[name] = status
	case "FinishedStatus":
		js.FinishedStatus[name] = Exists
		js.Information[name] = status
	case "TerminatedStatus":
		js.TerminatedStatus[name] = Exists
		js.Information[name] = status
	default:
		err := errors.New("unknown status")
		log.Error(err, "unknown status when setting", "status", status)
		return err
	}
	return nil
}

func (js *JobStatus) Get(deploymentName string) (string, error) {
	if val, ok := js.Information[deploymentName]; ok {
		return val, nil
	}
	return "", errors.New("not found")
}


func (js *JobStatus) CheckContains(name string) bool {
	val, ok := js.Information[name]
	if ok {
		log.Info("found exist job status pair", name, val)
	}
	return ok
}


/*
func (js *JobStatus) CheckIn(name string, status string) bool {
	var ret bool
	switch status {
	case "RunningStatus":
		_, ret = js.RunningStatus[name]
	case "WaitingStatus":
		_, ret = js.WaitingStatus[name]
	case "SchedulingStatus":
		_, ret = js.WaitingStatus[name]
	case "FinishedStatus":
		_, ret = js.FinishedStatus[name]
	case "TerminatedStatus":
		_, ret = js.TerminatedStatus[name]
	default:
		return false
	}
	return ret
}
*/


type JobStatus struct {
	RunningStatus map[string]void
	WaitingStatus map[string]void
	SchedulingStatus map[string]void
	FinishedStatus map[string]void
	TerminatedStatus map[string]void

	Information map[string]string
}
