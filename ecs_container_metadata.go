package ecstaskports

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type ecsInstanceTaskInjectedFileMetadata struct {
	PortMappings []struct {
		ContainerPort int    `json:"ContainerPort"`
		HostPort      int    `json:"HostPort"`
		BindIP        string `json:"BindIp"`
		Protocol      string `json:"Protocol"`
	}
}

func readInstanceMetadataFile(filePath string) ([]PortMapping, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read file %s", filePath)
	}

	var metadata ecsInstanceTaskInjectedFileMetadata
	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not unmarshal JSON in %s", filePath)
	}

	var addresses []PortMapping
	for _, addr := range metadata.PortMappings {
		addresses = append(
			addresses,
			PortMapping{
				ContainerPort: addr.ContainerPort,
				HostPort:      addr.HostPort,
			},
		)
	}

	return addresses, nil
}
