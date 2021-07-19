// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fluent-bit Addon E2E Test", func() {
	Specify("check fluent-bit is running", func() {
		var result int
		signal := make(chan *exec.Cmd)
		timeOut := make(chan bool)
		go func() {
			r, err := cmdHelperUp.CliRunnerChan("kubectl", nil, signal, cmdHelperUp.CommandArgs["k8s-port-forward"]...)
			println(r)
			Expect(err).To(HaveOccurred())
		}()
		command := <-signal
		tr := time.AfterFunc(ApiCallTimeout, func() {
			if command != nil {
				command.Process.Kill()
			}
			close(signal)
			println("command timed out")
			Expect(result).Should(Equal(200))
		})
		defer tr.Stop()

		for {
			resp, _ := http.Get(cmdHelperUp.CommandArgs["fluent-bit-health-check"][1])
			if resp != nil {
				result = resp.StatusCode
				if resp.StatusCode == 200 {
					fmt.Println("Test successfully done")
					if command != nil {
						command.Process.Kill()
					}
					Expect(resp.StatusCode).Should(Equal(200))
				}
			}
			if result != 0 {
				close(timeOut)
				close(signal)
				if command != nil {
					command.Process.Kill()
				}
				break
			}
		}
	})
})
