package flink

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// A Job is a representation for a Flink Job
type Job struct {
	ID     string `json:"jid"`
	Name   string `json:"name"`
	Status string `json:"state"`
}

type retrieveJobsOverviewResponse struct {
	RunningJobs   []string `json:"jobs-running"`
	FinishedJobs  []string `json:"jobs-finished"`
	CancelledJobs []string `json:"jobs-cancelled"`
	FailedJobs    []string `json:"jobs-failed"`
}

// RetrieveJobs returns all the jobs on the Flink cluster
func (c FlinkRestClient) RetrieveJobs() ([]Job, error) {
	res, err := c.Client.Get(c.constructURL("jobs"))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	jobs := []Job{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return jobs, err
	}

	if res.StatusCode != 200 {
		return jobs, fmt.Errorf("Unexpected response status %v with body %v", res.StatusCode, string(body[:]))
	}

	overviewResponse := retrieveJobsOverviewResponse{}
	err = json.Unmarshal(body, &overviewResponse)
	if err != nil {
		return jobs, fmt.Errorf("Unable to parse API response as valid JSON: %v", string(body[:]))
	}

	for _, resp := range [][]string{overviewResponse.RunningJobs, overviewResponse.FinishedJobs, overviewResponse.CancelledJobs, overviewResponse.FailedJobs} {
		for _, job := range resp {
			res, err := c.Client.Get(c.constructURL(fmt.Sprintf("%s/%s", "jobs", job)))
			if err != nil {
				fmt.Errorf("Failed to fetch job (%s) info. Status %v with body %v", job, res.StatusCode, string(body[:]))
				continue
			}
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			jobResponse := Job{}
			err = json.Unmarshal(body, &jobResponse)
			jobs = append(jobs, jobResponse)
		}
	}

	return jobs, nil
}
