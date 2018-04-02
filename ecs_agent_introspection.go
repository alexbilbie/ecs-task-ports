package ecstaskports

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/pkg/errors"
)

func discoverWithTaskIntrospection() ([]PortMapping, error) {
	im, err := getECSInstanceMetadata()
	if err != nil {
		return nil, err
	}

	dockerID, err := getDockerContainerID()
	if err != nil {
		return nil, err
	}

	tm, err := getECSTaskMetadata(dockerID)
	if err != nil {
		return nil, err
	}

	return getECSTaskPorts(im.Cluster, tm.Arn)
}

// ECS agent introspection URL as per https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-agent-introspection.html
const introspectionURL = "http://172.17.0.1:51678/v1/"

type ecsInstanceMetadata struct {
	Cluster              string `json:"Cluster"`
	ContainerInstanceArn string `json:"ContainerInstanceArn"`
	Version              string `json:"Version"`
}

type ecsInstanceTaskIntrospectionMetadata struct {
	Arn           string `json:"Arn"`
	DesiredStatus string `json:"DesiredStatus"`
	KnownStatus   string `json:"KnownStatus"`
	Family        string `json:"Family"`
	Version       string `json:"Version"`
	Containers    []struct {
		DockerID   string `json:"DockerId"`
		DockerName string `json:"DockerName"`
		Name       string `json:"Name"`
	} `json:"Containers"`
}

func getECSInstanceMetadata() (ecsInstanceMetadata, error) {
	var metadata ecsInstanceMetadata

	client := &http.Client{
		Timeout: time.Second * 3,
	}
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/metadata", introspectionURL), nil)
	resp, err := client.Do(req)
	if err != nil {
		return metadata, errors.Wrap(err, "Get instance metadata error")
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, &metadata)
	if err != nil {
		return metadata, errors.Wrap(err, "Unmarshal instance metadata error")
	}

	return metadata, nil
}

func getECSTaskMetadata(dockerID string) (ecsInstanceTaskIntrospectionMetadata, error) {
	var task ecsInstanceTaskIntrospectionMetadata

	client := &http.Client{
		Timeout: time.Second * 3,
	}
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/tasks?dockerid=%s", introspectionURL, dockerID),
		nil,
	)
	resp, err := client.Do(req)
	if err != nil {
		return task, errors.Wrap(err, "Get task metadata error")
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, &task)
	if err != nil {
		return task, errors.Wrap(err, "Unmarshal task metadata error")
	}

	return task, nil
}

func getECSTaskPorts(cluster string, taskARN string) ([]PortMapping, error) {
	var addresses []PortMapping

	s, _ := session.NewSession(&aws.Config{Region: aws.String(GetEC2InstanceRegion())})
	ecsClient := ecs.New(s)
	result, err := ecsClient.DescribeTasks(&ecs.DescribeTasksInput{
		Cluster: aws.String(cluster),
		Tasks: []*string{
			aws.String(taskARN),
		},
	})

	if err != nil {
		return addresses, errors.Wrap(err, "ECS DescribeTasks error")
	}

	var t *ecs.Task
	for _, task := range result.Tasks {
		if aws.StringValue(task.TaskArn) == taskARN {
			t = task
		}
	}

	if t == nil {
		return addresses, errors.New("No described task matched task ARN")
	}

	for _, nb := range t.Containers[0].NetworkBindings {
		addresses = append(
			addresses,
			PortMapping{
				ContainerPort: int(aws.Int64Value(nb.ContainerPort)),
				HostPort:      int(aws.Int64Value(nb.HostPort)),
			},
		)
	}

	return addresses, nil
}
