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
	"strings"

	"sigs.k8s.io/kubetest2/pkg/exec"
)

const (
	// kubectl = "./kubernetes/client/bin/kubectl"
	// kubectl    = "/home/prow/go/src/k8s.io/cloud-provider-gcp/cluster/kubectl.sh"
	// kubeconfig = "--kubeconfig=${ARTIFACTS}/kubetest2-kubeconfig"
	kubectl    = "${HOME}/devel/cloud-provider-gcp/cluster/kubectl.sh"
	kubeconfig = "--kubeconfig=${HOME}/devel/cloud-provider-gcp/_artifacts/kubetest2-kubeconfig"
)

// APIServerURL obtains the URL of the k8s master from kubectl
func APIServerURL() (string, error) {
	fmt.Println("STARTING API SERVER URL")
	lsresult, err := wrapAndRunWithLogging(
		"set -o xtrace; ls ${ARTIFACTS}; echo \"----\"; cat ${ARTIFACTS}/kubetest2-kubeconfig; echo \"-----\"; kubectl --kubeconfig=${ARTIFACTS}/kubetest2-kubeconfig config view -o jsonpath=\"{.current-context}\"")
	fmt.Println("RESULT")
	fmt.Println(lsresult)
	fmt.Println("ERROR")
	fmt.Println(err)
	fmt.Println("END")

	fmt.Println("COMMAND")
	command := []string{kubectl, kubeconfig, "config", "view", "-o", "jsonpath=\"{.current-context}\""}
	fmt.Println(command)
	fmt.Println("END, RUNNING")

	// kubecontext, err := execAndResult(kubectl, kubeconfig, "config", "view", "-o", "jsonpath=\"{.current-context}\"")
	kubecontext, err := execAndResult(command[0], command[1:]...)
	if err != nil {
		return "", fmt.Errorf("Could not get kube context: %v", err)
	}

	fmt.Println("KUBECONTEXT")
	fmt.Println(kubecontext)
	fmt.Println("END")

	fmt.Println("CLUSTERNAME COMMAND")
	clusternameCommand := []string{kubectl, kubeconfig, "config", "view", "-o",
		fmt.Sprintf("jsonpath=\"{.contexts[?(@.name == \\\"%s\\\")].context.cluster}\"", kubecontext)}
	fmt.Println(clusternameCommand)
	fmt.Println("END, RUNNING")
	// clustername, err := execAndResult(clusternameCommand[0], clusternameCommand[1:]...)
	clustername, err := wrapAndRunWithLogging(clusternameCommand[0], clusternameCommand[1:]...)
	if err != nil {
		return "", fmt.Errorf("Could not get cluster name: %v", err)
	}

	fmt.Println("APISERVER COMMAND")
	//apiServerURL, err := execAndResult(kubectl, kubeconfig, "config", "view", "-o",
	//	fmt.Sprintf("jsonpath={.clusters[?(@.name == %s)].cluster.server}", clustername))
	apiServerURL, err := wrapAndRunWithLogging(kubectl, kubeconfig, "config", "view", "-o",
		fmt.Sprintf("jsonpath=\"{.clusters[?(@.name == \\\"%s\\\")].cluster.server}\"", clustername))
	if err != nil {
		return "", err
	}
	return apiServerURL, nil
}

func wrapInBashC(command string, args ...string) (string, []string) {
	if len(args) == 0 {
		return "bash", []string{"-c", command}
	}

	return "bash", []string{
		"-c",
		fmt.Sprintf("%s %s", command, strings.Join(args, " ")),
	}
}

// execAndResult runs command with args and returns the entire output (or error)
func execAndResult(command string, args ...string) (string, error) {
	// cmd := exec.Command(command, args...)
	command, args = wrapInBashC(command, args...)
	fmt.Printf("COMMAND: %+v\n", command)
	fmt.Printf("ARGS: %+v\n", args)
	cmd := exec.Command(command, args...)

	cmd.SetStderr(os.Stderr)
	bytes, err := exec.Output(cmd)
	return string(bytes), err
}

func wrapAndRunWithLogging(command string, args ...string) (string, error) {
	fmt.Printf("Before wrapping, command is: %#v\n", command)
	fmt.Printf("Before wrapping, args is: %#v\n", args)
	wrappedCommand, wrappedArgs := wrapInBashC(command, args...)
	fmt.Printf("After wrapping, command is: %#v\n", wrappedCommand)
	fmt.Printf("After wrapping, args is: %#v\n", wrappedArgs)

	cmd := exec.Command(wrappedCommand, wrappedArgs...)
	cmd.SetStderr(os.Stderr)
	fmt.Println("running")
	bytes, err := exec.Output(cmd)
	if err != nil {
		fmt.Printf("error from run: %s\n", err)
	}
	fmt.Println()
	fmt.Printf("bytes: %#v\n", bytes)
	fmt.Printf("str result: %s\n", bytes)
	fmt.Println("done")
	fmt.Println()

	return string(bytes), err
}
