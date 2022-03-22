// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package lint

import "fmt"

const (
	colorReset        = "\033[0m"  // Reset
	colorGreen        = "\033[32m" // General Titles
	colorBlue         = "\033[34m" // Pass
	colorRed          = "\033[31m" // Fail
	colorLightMagenta = "\033[95m" // Not Identified
	colorMagenta      = "\033[35m" // Pull Failed
)
const (
	Pass          = "Pass"
	Fail          = "Fail"
	NotIdentified = "Not Identified"
	PullFail      = "Pull Failed"
)

func (imc *ImageLintConfig) ShowOverallSummary() {
	pass := 0
	fail := 0
	notIdentified := 0
	pullFailed := 0
	for _, imgs := range imc.ImageMap {
		if len(imgs) > 0 {
			img := imgs[0]
			switch img.Status {
			case Pass:
				pass += len(imgs)
			case Fail:
				fail += len(imgs)
			case NotIdentified:
				notIdentified += len(imgs)
			case PullFail:
				pullFailed += len(imgs)
			}
		}
	}
	fmt.Println("-------------------------------", string(colorGreen), "Summary", string(colorReset), "-----------------------------------------------")
	fmt.Println("Total  Images                :", len(imc.ImageMap))
	fmt.Println("Total  Occurrences           :", pass+fail+notIdentified+pullFailed)
	fmt.Println("Total", string(colorBlue), "Pass", string(colorReset), "                :", pass)
	fmt.Println("Total", string(colorRed), "Fail", string(colorReset), "                :", fail)
	fmt.Println("Total", string(colorMagenta), "Pull Failed", string(colorReset), "         :", pullFailed)
	fmt.Println("Total", string(colorLightMagenta), "Not Identified", string(colorReset), "      :", notIdentified)
	fmt.Println("[Note:All totals for Pass|Fail|Not Identified|Pull Failed are based on number of occurrences]")
	fmt.Println("--------------------------------------------------------------------------------------------")
}

func (imc *ImageLintConfig) ShowSummary(stype string) {
	headingColor := string(colorGreen)
	color := ""
	switch stype {
	case Pass:
		color = string(colorBlue)
	case Fail:
		color = string(colorRed)
	case NotIdentified:
		color = string(colorLightMagenta)
	case PullFail:
		color = string(colorMagenta)
	default:
		return
	}
	fmt.Println("-------------------------------", headingColor, stype, " Summary", string(colorReset), "---------------------------------------------")
	count := 0
	for image, imgs := range imc.ImageMap {
		if len(imgs) > 0 {
			img := imgs[0]
			if img.Status == stype {
				count += len(imgs)
				fmt.Println("Image:    ", image)
				fmt.Println("File Path:", img.Path)
				fmt.Println("Image Position:", img.Position.Row, ":", img.Position.Col)
				fmt.Println("")
			}
		}
	}
	fmt.Println("Total", color, stype, string(colorReset), ":", count)
}

func (imc *ImageLintConfig) ShowDetailedSummary(detailedSummary string) {
	switch detailedSummary {
	case "Fail", "fail", "FAIL":
		imc.ShowSummary("Fail")
	case "Pass", "pass", "PASS":
		imc.ShowSummary("Pass")
	case "Not Identified", "not identified", "NOT IDENTIFIED", "Not identified", "not Identified":
		imc.ShowSummary("Not Identified")
	case "Pull Failed", "", "PULL FAILED", "Pull failed", "pull Failed":
		imc.ShowSummary("Pull Failed")
	case "ALL", "all":
		imc.ShowSummary("Fail")
		imc.ShowSummary("Pass")
		imc.ShowSummary("Pull Failed")
		imc.ShowSummary("Not Identified")
	}
}

func (imc *ImageLintConfig) OnEvent(etype, message, image string) {
	color := ""
	switch etype {
	case Pass:
		color = string(colorBlue)
	case Fail:
		color = string(colorRed)
	case NotIdentified:
		color = string(colorLightMagenta)
	case PullFail:
		color = string(colorMagenta)
	default:
		return
	}
	fmt.Println("Status: ", color, etype, string(colorReset))
	fmt.Println("Image:   ", image)
	fmt.Println("Message: ", message)
	fmt.Println("Total ", len(imc.ImageMap[image]), " file(s) contain(s) this image")
	for i, img := range imc.ImageMap[image] {
		fmt.Println("File Path:", img.Path)
		fmt.Println("Image Position:", img.Position.Row, ":", img.Position.Col)
		imc.ImageMap[image][i].Status = etype
	}
	fmt.Println()
}
