/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kubectl

import (
	"fmt"
	"os"

	"sigs.k8s.io/kubetest2/pkg/exec"
)

const (
	// kubectl = "./kubernetes/client/bin/kubectl"
	kubectl    = "/home/prow/go/src/k8s.io/cloud-provider-gcp/cluster/kubectl.sh"
	kubeconfig = "--kubeconfig=${ARTIFACTS}/kubetest2-kubeconfig"
)

// APIServerURL obtains the URL of the k8s master from kubectl
func APIServerURL() (string, error) {
	fmt.Println("STARTING API SERVER URL")
	lsresult, err := execAndResult("bash", "-c",
		"set -o xtrace; ls ${ARTIFACTS}; cat ${ARTIFACTS}/kubetest2-kubeconfig")
	fmt.Println("RESULT")
	fmt.Println(lsresult)
	fmt.Println("ERROR")
	fmt.Println(err)
	fmt.Println("END")

	fmt.Println("COMMAND")
	command := []string{kubectl, kubeconfig, "config", "view", "-o", "jsonpath=\"{.current-context}\""}
	fmt.Println("END, RUNNING")

	// kubecontext, err := execAndResult(kubectl, kubeconfig, "config", "view", "-o", "jsonpath=\"{.current-context}\"")
	kubecontext, err := execAndResult(command[0], command[1:]...)
	if err != nil {
		return "", fmt.Errorf("Could not get kube context: %v", err)
	}

	fmt.Println("KUBECONTEXT")
	fmt.Println(kubecontext)
	fmt.Println("END")

	clustername, err := execAndResult(kubectl, kubeconfig, "config", "view", "-o",
		fmt.Sprintf("jsonpath=\"{.contexts[?(@.name == %s)].context.cluster}\"", kubecontext))
	if err != nil {
		return "", fmt.Errorf("Could not get cluster name: %v", err)
	}

	apiServerURL, err := execAndResult(kubectl, kubeconfig, "config", "view", "-o",
		fmt.Sprintf("jsonpath={.clusters[?(@.name == %s)].cluster.server}", clustername))
	if err != nil {
		return "", err
	}
	return apiServerURL, nil
}

// execAndResult runs command with args and returns the entire output (or error)
func execAndResult(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.SetStderr(os.Stderr)
	bytes, err := exec.Output(cmd)
	return string(bytes), err
}
