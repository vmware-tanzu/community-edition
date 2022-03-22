// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package lint

import "fmt"

const (
	colorReset = "\033[0m" // Reset
	colorBlue  = "\033[34m"
	colorRed   = "\033[31m" // Fail
	passConst  = "Pass"
	failConst  = "Fail"
	//	colorPurple  = "\033[35m" // Not Identified
	colorGreen = "\033[32m"

//	colorMagenta = "\033[35m" // Pull Failed
)

func (llc *LinkLintConfig) ShowSummary() {
	pass := 0
	fail := 0
	for _, lls := range llc.LinkMap {
		if len(lls) > 0 {
			ll := lls[0]
			switch ll.Status {
			case string(passConst):
				pass += len(lls)
			case string(failConst):
				fail += len(lls)
			}
		}
	}
	fmt.Println("-------------------------------", string(colorGreen), "Summary", string(colorReset), "---------------------------------------------")
	fmt.Println("Total  Links                 :", len(llc.LinkMap))
	fmt.Println("Total  Occurrences           :", pass+fail)
	fmt.Println("Total", string(colorBlue), string(passConst), string(colorReset), "                :", pass)
	fmt.Println("Total", string(colorRed), string(failConst), string(colorReset), "                :", fail)
	fmt.Println("[Note:All totals for Pass|Fail are based on number of files]")
	fmt.Println("---------------------------------------------------------------------------------------")
}

func (llc *LinkLintConfig) ShowFailSummary() {
	fmt.Println("-------------------------------", string(colorGreen), "Fail Summary", string(colorReset), "---------------------------------------------")
	fail := Summary(llc, failConst)
	fmt.Println("Total", string(colorRed), string(failConst), string(colorReset), "          :", fail)
	fmt.Println("---------------------------------------------------------------------------------------")
}

func (llc *LinkLintConfig) ShowPassSummary() {
	fmt.Println("-------------------------------", string(colorGreen), "Pass Summary", string(colorReset), "---------------------------------------------")
	pass := Summary(llc, passConst)
	fmt.Println("Total", string(colorBlue), string(passConst), string(colorReset), "          :", pass)
	fmt.Println("---------------------------------------------------------------------------------------")
}

func Summary(llc *LinkLintConfig, statusConst string) int {
	count := 0
	for _, lls := range llc.LinkMap {
		if len(lls) > 0 {
			ll := lls[0]
			if ll.Status == statusConst {
				count += len(lls)
				fmt.Println("File Path:", ll.Path)
				fmt.Println("URL:", ll.Line)
				fmt.Println("Url Position:", ll.Position.Row, ":", ll.Position.Col)
				fmt.Println("Message:", ll.Message)
				fmt.Println("")
			}
		}
	}
	return count
}

func (llc *LinkLintConfig) OnPass(message, link string) {
	URLDetails(llc, colorBlue, link, message, passConst)
}

func (llc *LinkLintConfig) OnFail(message, link string) {
	URLDetails(llc, colorRed, link, message, failConst)
}

func URLDetails(llc *LinkLintConfig, color, link, message, statusConst string) {
	fmt.Println("Status:", color, statusConst, string(colorReset))
	fmt.Println("URL:    ", link)
	fmt.Println("Error:  ", message)
	fmt.Println("Total ", len(llc.LinkMap[link]), " file(s) contain(s) this URL")
	for i, ll := range llc.LinkMap[link] {
		fmt.Println("File Path:", ll.Path)
		fmt.Println("URL Position:", ll.Position.Row, ":", ll.Position.Col)
		llc.LinkMap[link][i].Message = message
		llc.LinkMap[link][i].Status = statusConst
	}
	fmt.Println()
}
