package ecstaskports

import (
	"context"
	"net/http"

	"io/ioutil"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

const (
	instanceMetadataEndpointLocalIPv4 = "http://169.254.169.254/latest/meta-data/local-ipv4"
	instanceMetadataEndpointAZ        = "http://169.254.169.254/latest/meta-data/placement/availability-zone"
)

// GetEC2InstancePrivateIPAddress returns the private IP address of an EC2 instance
func GetEC2InstancePrivateIPAddress() ([]byte, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	resp, err := ctxhttp.Get(ctx, http.DefaultClient, instanceMetadataEndpointLocalIPv4)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

// GetEC2InstanceAvailabilityZone returns the availability zone that an EC2 instance is running in
func GetEC2InstanceAvailabilityZone() ([]byte, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	resp, err := ctxhttp.Get(ctx, http.DefaultClient, instanceMetadataEndpointAZ)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

// GetEC2InstanceRegion automatically discovers the AWS region from the runtime environment
func GetEC2InstanceRegion() string {
	region := os.Getenv("AWS_REGION")

	if region == "" {
		region = os.Getenv("AWS_DEFAULT_REGION")
	}

	if region == "" {
		data, err := GetEC2InstanceAvailabilityZone()
		if err == nil {
			region = strings.TrimSpace(string(data))
			if len(region) > 0 {
				region = region[0 : len(region)-1]
			}
		}
	}

	if region == "" {
		region = "us-east-1"
	}

	return region
}
