package flink

import (
	"fmt"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

// Cancel terminates a running job specified by job ID
func (c FlinkRestClient) Cancel(jobID string) error {
	req, err := retryablehttp.NewRequest("PATCH", c.constructURL(fmt.Sprintf("jobs/%v", jobID)), nil)
	if err != nil {
		return err
	}

	res, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 202 {
		return fmt.Errorf("Unexpected response status %v", res.StatusCode)
	}

	return nil
}
