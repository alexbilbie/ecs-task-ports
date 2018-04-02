package ecstaskports

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/pkg/errors"
)

func getDockerContainerID() (string, error) {
	b, err := ioutil.ReadFile("/proc/1/cpuset")
	if err != nil {
		return "", errors.Wrap(err, "Could not read Docker ID from /proc/1/cpuset")
	}

	return strings.TrimSpace(path.Base(string(b))), nil
}
