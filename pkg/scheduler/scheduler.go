package scheduler

import (
	"errors"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("scheduler")
var JobStatusData JobStatus


func SetJobStatus(deploymentName string, status string) error {
	log.Info("SetJobStatus begin with", "deploymentName", deploymentName, "status", status)
	return JobStatusData.Set(deploymentName, status)
}

func GetJobStatus(deploymentName string) (string, error) {
	if exists := JobStatusData.CheckContains(deploymentName); exists {
		return JobStatusData.Information[deploymentName], nil
	}
	log.Info("not found deployment", "deploymentname", deploymentName)
	return "", errors.New("not found")
}