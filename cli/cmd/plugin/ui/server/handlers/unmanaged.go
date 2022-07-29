// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-openapi/runtime/middleware"

	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/pkg/containerruntime"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/models"
	"github.com/vmware-tanzu/community-edition/cli/cmd/plugin/ui/server/restapi/operations/unmanaged"
)

// unmanagedCluster is the cluster information we get from the CLI output.
type unmanagedCluster struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Status   string `json:"status"`
}

// checkUnmanagedPlugin verifies we can execute the "tanzu unmanaged-cluster" command.
func checkUnmanagedPlugin() error {
	cmd := exec.Command("tanzu", "unmanaged-cluster")
	err := cmd.Run()
	if err != nil {
		return errors.New("tanzu unmanaged-cluster could not be found")
	}
	return nil
}

func checkDockerContainerRunning() error {
	_, err := containerruntime.GetRuntimeInfo()
	return err
}

func runCommand(args ...string) (string, error) {
	cmdArgs := []string{"unmanaged-cluster"}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("tanzu", cmdArgs...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func parseCommandOutput(intput string, err error) (string, error) {
	if err == nil {
		var result = strings.Replace(intput, "\n", "", -1)
		result = strings.Replace(result, " ", "", -1)
		re := regexp.MustCompile(`(.*)(\[(?s).*\]$)`)
		return re.FindStringSubmatch(result)[2], nil
	}
	return "", err
}

// CreateUnmanagedCluster creates a new unmanaged cluster.
func (app *App) CreateUnmanagedCluster(params unmanaged.CreateUnmanagedClusterParams) middleware.Responder {
	fmt.Println("checking unmanaged plugin...")
	if err := checkUnmanagedPlugin(); err != nil {
		fmt.Println("ERROR: unmanaged plugin cannot be found")
		return unmanaged.NewDeleteUnmanagedClusterInternalServerError().WithPayload(Err(err))
	}
	fmt.Println("unmanaged plugin found")

	fmt.Println("checking docker container running...")
	if err := checkDockerContainerRunning(); err != nil {
		fmt.Println("ERROR: docker container cannot be found")
		return unmanaged.NewDeleteUnmanagedClusterInternalServerError().WithPayload(Err(err))
	}
	fmt.Println("docker container running")

	createParams := params.Params
	if createParams.Name == "" {
		return unmanaged.NewCreateUnmanagedClusterBadRequest().WithPayload(Err(fmt.Errorf("cluster name must be provided")))
	}

	args := []string{
		"create",
		createParams.Name,
	}

	if createParams.Provider != "" {
		args = append(args, "--provider", createParams.Provider)
	}

	if createParams.Cni != "" {
		args = append(args, "--cni", createParams.Cni)
	}

	if createParams.Podcidr != "" {
		args = append(args, "--pod-cidr", createParams.Podcidr)
	}

	if createParams.Servicecidr != "" {
		args = append(args, "--service-cidr", createParams.Servicecidr)
	}

	if len(createParams.Portmappings) > 0 {
		for _, pm := range createParams.Portmappings {
			args = append(args, "--port-map", pm)
		}
	}

	if createParams.Controlplanecount > 0 {
		args = append(args, "--control-plane-node-count", strconv.FormatInt(createParams.Controlplanecount, 10))
	}

	if createParams.Workernodecount > 0 {
		args = append(args, "--worker-node-count", strconv.FormatInt(createParams.Workernodecount, 10))
	}
	args = append(args, "--tty-activate")

	go executeAndRedirect(args...)

	creatingClusterPayload := &models.UnmanagedCluster{
		Name:     createParams.Name,
		Provider: createParams.Provider,
		Status:   "creating",
	}
	return unmanaged.NewCreateUnmanagedClusterOK().WithPayload(creatingClusterPayload)
}

// GetUnmanagedCluster gets details of a specific unmanaged cluster.
func (app *App) GetUnmanagedCluster(params unmanaged.GetUnmanagedClusterParams) middleware.Responder {
	if err := checkUnmanagedPlugin(); err != nil {
		return unmanaged.NewGetUnmanagedClusterInternalServerError().WithPayload(Err(err))
	}

	cluster, err := app.getUnmanagedCluster(params.ClusterName)
	if err != nil {
		return unmanaged.NewGetUnmanagedClusterInternalServerError().WithPayload(Err(err))
	}

	return unmanaged.NewGetUnmanagedClusterOK().WithPayload(cluster)
}

// getUnmanagedCluster gets details for a specific cluster.
func (app *App) getUnmanagedCluster(clusterName string) (*models.UnmanagedCluster, error) {
	clusters, err := app.getUnmanagedClusters()
	if err != nil {
		return nil, err
	}

	for _, cluster := range clusters {
		if cluster.Name == clusterName {
			return cluster, nil
		}
	}

	return nil, fmt.Errorf("unmanaged cluster %q could not be found", clusterName)
}

// getUnmanagedClusters gets a list of all unmanaged clusters.
func (app *App) getUnmanagedClusters() ([]*models.UnmanagedCluster, error) {
	jsonOutput, err := parseCommandOutput(runCommand("list", "-o", "json"))
	if err != nil {
		return nil, err
	}

	var clusters []unmanagedCluster
	err = json.Unmarshal([]byte(jsonOutput), &clusters)
	if err != nil {
		return nil, fmt.Errorf("unable to parse unmanaged cluster information: %s", err.Error())
	}

	results := []*models.UnmanagedCluster{}
	for _, cluster := range clusters {
		results = append(results, &models.UnmanagedCluster{
			Name:     cluster.Name,
			Provider: cluster.Provider,
			Status:   cluster.Status,
		})
	}

	return results, nil
}

// GetUnmanagedClusters gets all unmanaged clusters.
func (app *App) GetUnmanagedClusters(params unmanaged.GetUnmanagedClustersParams) middleware.Responder {
	if err := checkUnmanagedPlugin(); err != nil {
		return unmanaged.NewGetUnmanagedClusterInternalServerError().WithPayload(Err(err))
	}

	clusters, err := app.getUnmanagedClusters()
	if err != nil {
		return unmanaged.NewGetUnmanagedClusterInternalServerError().WithPayload(Err(err))
	}

	return unmanaged.NewGetUnmanagedClustersOK().WithPayload(clusters)
}

// DeleteUnmanagedCluster triggers the deletion of an unmanaged cluster.
func (app *App) DeleteUnmanagedCluster(params unmanaged.DeleteUnmanagedClusterParams) middleware.Responder {
	if err := checkUnmanagedPlugin(); err != nil {
		return unmanaged.NewDeleteUnmanagedClusterInternalServerError().WithPayload(Err(err))
	}

	// TODO: Add streaming of delete output
	// TODO #2: There is some talk of adding an "are you sure" prompt to the
	// tanzu uc rm command. If so, when we update to that release we will need
	// to add a "-y" argument here (if that's what ends up being implemented) to
	// tell it to not prompt for confirmation.
	_, err := runCommand("delete", params.ClusterName)
	if err != nil {
		return unmanaged.NewDeleteUnmanagedClusterInternalServerError().WithPayload(Err(err))
	}

	return unmanaged.NewDeleteUnmanagedClusterOK()
}

func sendTanzuCommand(tanzuArgs []string) {
	doSendLog([]byte(formatTanzuCommandMessage(tanzuArgs)))
}

func customTrim(src string) string {
	// We want to trim all spaces from the RIGHT side, and just "\r\n" from the left.
	// For the right side, we COULD use TrimRight(), but
	// that function requires us to specify the white characters,
	// so instead we use TrimSpace to find the central non-whitespace string
	// and then pull it out of the main string ALONG with any left side whitespace.
	noSpaceMessage := strings.TrimSpace(src)
	xStartNoSpaceMessage := strings.Index(src, noSpaceMessage)
	rightTrimmedMessage := src[0 : xStartNoSpaceMessage+len(noSpaceMessage)]
	// THEN we remove \r\n from the left
	return strings.TrimLeft(rightTrimmedMessage, "\r\n")
}

// Our sendLog is a wrapper to the SendLog method.
// The SendLog method is expecting a byte array of a valid JSON object, so our sendLog method
// - splits the raw byte array into several messages along "\n" boundaries,
// - wraps the resulting messages in JSON as a formatted "log" message, and
// - calls doSendLog (which calls SendLog).
func sendLog(msg []byte) {
	// split the message on "\n" boundaries
	rawMessages := strings.Split(string(msg), "\n")
	// remove messages that only contain white space
	var messages []string
	for x := range rawMessages {
		trimmedMessage := customTrim(rawMessages[x])
		if len(trimmedMessage) > 0 {
			messages = append(messages, trimmedMessage)
		}
	}
	// Some "messages" are ONLY white space. When we encounter them, we send a "ping" message.
	// The interpretation: the white space "message" is just the process saying "I'm still alive"
	if len(messages) == 0 {
		doSendLog([]byte(formatPingMessage()))
		return
	}

	for _, singleMessage := range messages {
		fmt.Printf("%s: Sending LOG message: [%s]\n", time.Now(), singleMessage)
		doSendLog([]byte(formatLogMessage(singleMessage)))
		// NOTE: "RESULT: FAILURE" and "RESULT: SUCCESS" are magic strings output from the unmanaged-cluster plugin
		if singleMessage == "RESULT: FAILURE" {
			sendResultFailure()
		} else if singleMessage == "RESULT: SUCCESS" {
			sendResultSuccess()
		}
	}
}

func sendResultFailure() {
	// NOTE: "failed" is a magic string recognized by the UI front end
	doSendLog([]byte(formatResultMessage("failed")))
}

func sendResultSuccess() {
	// NOTE: "successful" is a magic string recognized by the UI front end
	doSendLog([]byte(formatResultMessage("successful")))
}

func doSendLog(msg []byte) {
	SendLog(msg)
	sleepAfterSendingWebSocketMessage()
}

func sleepAfterSendingWebSocketMessage() {
	// A hack to avoid losing messages, which happens if we write another message too quickly to the websocket
	time.Sleep(25 * time.Millisecond)
}

// NOTE: all the formatXXX methods are structuring JSON in a special way that is expected by the front end
func formatResultMessage(result string) string {
	return fmt.Sprintf("{\"type\":\"progress\", \"data\":{\"status\":%q,\"message\":\"RESULT: %s\"}}", result, result)
}

func formatLogMessage(logMessage string) string {
	// Because we are sending the raw messages wrapped in JSON, we need to escape any double quotes, hence %q
	return fmt.Sprintf("{\"type\":\"log\", \"data\":{\"logType\":\"output\",\"message\":%q}}", logMessage)
}

func formatPingMessage() string {
	return "{\"type\":\"ping\"}"
}

func formatTanzuCommandMessage(args []string) string {
	tanzuCmd := "tanzu " + strings.Join(args, " ")
	return fmt.Sprintf("{\"type\":\"log\", \"data\":{\"logType\":\"tanzu command\",\"message\":%q}}", tanzuCmd)
}

// Adapted from copyAndCapture at https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
// Copies data from one stream to the other AND sends it as a Log message to the web socket
func copyAndLog(w io.Writer, r io.Reader) error {
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			d := buf[:n]
			sendLog(d)
			_, err := w.Write(d)
			if err != nil {
				return err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return err
		}
	}
}

func executeAndRedirect(args ...string) {
	cmdArgs := []string{"unmanaged-cluster"}
	cmdArgs = append(cmdArgs, args...)

	fmt.Println("Tanzu CLI command: 'tanzu " + strings.Join(cmdArgs, " ") + "'")
	sendTanzuCommand(cmdArgs)
	cmd := exec.Command("tanzu", cmdArgs...)

	var errStdout, errStderr error
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	fmt.Println("starting command with cmd.Start()")
	err := cmd.Start()
	if err != nil {
		fmt.Printf("cmd.Start() failed with '%s'\n", err)
		return
	}

	// cmd.Wait() should be called only after we finish reading from stdoutIn and stderrIn.
	// wg ensures that we finish
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		errStdout = copyAndLog(os.Stdout, stdoutIn)
		wg.Done()
	}()

	errStderr = copyAndLog(os.Stderr, stderrIn)

	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		fmt.Printf("failed to capture stdout or stderr\n")
	}
}
