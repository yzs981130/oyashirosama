package scheduler

import (
	"errors"
	"github.com/onsi/gomega"
	"testing"
)


//TODO: replace with gomega

type Test struct {
	deploymentName string
	status string
}

var tests = []Test {
	{"sample-deployment", "RunningStatus"},
	{"yet-another-sample-deployment", "SchedulingStatus"},
	{"another-sample-deployment", "WaitingStatus"},
	{"yet-another-sample-deployment", "TerminatedStatus"},
	{"yet-another-sample-deployment", "FinishedStatus"},
	{"yet-another-sample-deployment", "OtherErrStatus"},
}


func TestJobStatus(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	var set JobStatus
	set.init()
	for _, v := range tests {
		if v.status == "OtherErrStatus" {
			g.Expect(set.Set(v.deploymentName, v.status)).To(gomega.HaveOccurred())
			continue
		}
		g.Expect(set.Set(v.deploymentName, v.status)).NotTo(gomega.HaveOccurred())
		g.Expect(set.Information[v.deploymentName]).To(gomega.Equal(v.status))
		g.Expect(set.Get(v.deploymentName)).To(gomega.Equal(set.Information[v.deploymentName]))
	}
	_, err := set.Get("unknown deployment name")
	g.Expect(err).To(gomega.MatchError(errors.New("not found")))
}

func TestSetJobStatus(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	for _, v := range tests {
		if v.status == "OtherErrStatus" {
			g.Expect(SetJobStatus(v.deploymentName, v.status)).To(gomega.HaveOccurred())
			continue
		}
		g.Expect(SetJobStatus(v.deploymentName, v.status)).NotTo(gomega.HaveOccurred())
	}

}

func TestGetJobStatus(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	g.Expect(GetJobStatus("sample-deployment")).To(gomega.Equal("RunningStatus"))
	g.Expect(GetJobStatus("another-sample-deployment")).To(gomega.Equal("WaitingStatus"))
	g.Expect(GetJobStatus("yet-another-sample-deployment")).To(gomega.Equal("FinishedStatus"))
}