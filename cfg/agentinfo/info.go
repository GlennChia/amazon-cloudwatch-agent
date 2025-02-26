// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package agentinfo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aws/amazon-cloudwatch-agent/cfg/envconfig"
)

const versionFilename = "CWAGENT_VERSION"

// We will fall back to a major version if no valid version file is found
const fallbackVersion = "1"

var isRunningAsRoot = func() bool {
	return os.Getuid() == 0
}

var (
	VersionStr    string
	BuildStr      string = "No Build Date"
	InputPlugins  []string
	OutputPlugins []string

	userAgent string
)

func Version() string {
	if VersionStr != "" {
		return VersionStr
	}

	version, err := readVersionFile()
	if err != nil {
		return fallbackVersion
	}

	VersionStr = version
	return version
}

func Build() string {
	return BuildStr
}

func Plugins() string {
	outputs := strings.Join(OutputPlugins, " ")
	inputs := strings.Join(InputPlugins, " ")

	if !isRunningAsRoot() {
		inputs += " run_as_user" // `inputs` is never empty, or agent will not start
	}

	return fmt.Sprintf("inputs:(%s) outputs:(%s)", inputs, outputs)
}

func UserAgent() string {
	if userAgent == "" {
		ua := os.Getenv(envconfig.CWAGENT_USER_AGENT)
		if ua == "" {
			ua = fmt.Sprintf("%s %s", FullVersion(), Plugins())
		}
		userAgent = ua
	}
	return userAgent
}

func FullVersion() string {
	return fmt.Sprintf("CWAgent/%s (%s; %s; %s) %s", Version(), runtime.Version(), runtime.GOOS, runtime.GOARCH, Build())
}

func readVersionFile() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("cannot get the path for current executable binary: %v", err)
	}
	curPath := filepath.Dir(ex)
	versionFilePath := filepath.Join(curPath, versionFilename)
	if _, err := os.Stat(versionFilePath); err != nil {
		return "", fmt.Errorf("the agent version file %s does not exist: %v", versionFilePath, err)
	}

	byteArray, err := ioutil.ReadFile(versionFilePath)
	if err != nil {
		return "", fmt.Errorf("issue encountered when reading content from file %s: %v", versionFilePath, err)
	}

	//TODO we may consider to do a format checking for the Version value.
	return strings.Trim(string(byteArray), " \n\r\t"), nil
}
