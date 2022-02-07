// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"cloud.google.com/go/storage"
)

const (
	regexReadmeLinuxAmd64   string = "\\[Linux AMD64.*\\]\\(.*\\)"
	regexReadmeDarwinAmd64  string = "\\[Darwin AMD64.*\\]\\(.*\\)"
	regexReadmeWindowsAmd64 string = "\\[Windows AMD64.*\\]\\(.*\\)"

	projectID      string = "vmware-cna"
	bucketName     string = "tce-cli-plugins-staging"
	gcpFile        string = "gcptokenfile.json"
	fullPathReadme string = "../../README.md"
)

type Builds struct {
	BuildDate    string
	DarwinAmd64  string
	DarwinArm64  string
	WindowsAmd64 string
	LinuxAmd64   string
}

type ClientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

// UploadFile uploads an object
func (c *ClientUploader) UploadFile(filename string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*1200)
	defer cancel()

	// open the src file
	uploadFile := filepath.Join("..", "..", "build", filename)
	file, err := os.OpenFile(uploadFile, os.O_RDWR, 0755)
	if err != nil {
		fmt.Printf("OpenFile failed: %v\n", err)
		return err
	}

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath + filename).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		fmt.Printf("io.Copy: %v\n", err)
		return err
	}
	if err := wc.Close(); err != nil {
		fmt.Printf("Writer.Close: %v\n", err)
	}

	return nil
}

func uploadTarball(buildDate, platform, arch, tag, extension string) (string, error) {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", gcpFile)

	client, err := storage.NewClient(context.Background())
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return "", err
	}

	assetFilename := fmt.Sprintf("tce-%s-%s-%s.%s", platform, arch, tag, extension)
	uploadPathTmp := filepath.Join("build-daily", buildDate)
	uploadPath := fmt.Sprintf("%s/", uploadPathTmp)

	uploader := &ClientUploader{
		cl:         client,
		bucketName: bucketName,
		projectID:  projectID,
		uploadPath: uploadPath,
	}

	err = uploader.UploadFile(assetFilename)
	if err != nil {
		fmt.Printf("UploadFile failed: %v\n", err)
		return "", err
	}

	urlUpload := fmt.Sprintf("https://storage.googleapis.com/%s/%s%s", bucketName, uploadPath, assetFilename)
	return urlUpload, nil
}

func uploadAll(tag, gcpToken string) (*Builds, error) {
	err := os.WriteFile(gcpFile, []byte(gcpToken), 0644)
	if err != nil {
		fmt.Printf("WriteFile failed: %v\n", err)
		return nil, err
	}

	buildDate := time.Now().Format("2006-01-02")

	urlLinuxAmd64, err := uploadTarball(buildDate, "linux", "amd64", tag, "tar.gz")
	if err != nil {
		fmt.Printf("uploadLinux failed: %v\n", err)

		errFile := os.Remove(gcpFile)
		if err != nil {
			fmt.Printf("Remove failed: %v\n", errFile)
		}

		return nil, err
	}

	urlWindowsAmd64, err := uploadTarball(buildDate, "windows", "amd64", tag, "zip")
	if err != nil {
		fmt.Printf("uploadWindows failed: %v\n", err)

		errFile := os.Remove(gcpFile)
		if err != nil {
			fmt.Printf("Remove failed: %v\n", errFile)
		}

		return nil, err
	}

	urlDarwinAmd64, err := uploadTarball(buildDate, "darwin", "amd64", tag, "tar.gz")
	if err != nil {
		fmt.Printf("uploadDarwinAmd64 failed: %v\n", err)

		errFile := os.Remove(gcpFile)
		if err != nil {
			fmt.Printf("Remove failed: %v\n", errFile)
		}

		return nil, err
	}

	_, err = uploadTarball(buildDate, "darwin", "arm64", tag, "tar.gz")
	if err != nil {
		fmt.Printf("uploadDarwinArm64 failed: %v\n", err)

		errFile := os.Remove(gcpFile)
		if err != nil {
			fmt.Printf("Remove failed: %v\n", errFile)
		}

		return nil, err
	}

	errFile := os.Remove(gcpFile)
	if err != nil {
		fmt.Printf("Remove failed: %v\n", errFile)
	}

	builds := &Builds{
		BuildDate:    buildDate,
		LinuxAmd64:   urlLinuxAmd64,
		WindowsAmd64: urlWindowsAmd64,
		DarwinAmd64:  urlDarwinAmd64,
	}

	return builds, nil
}

func updateReadme(builds *Builds) error {
	readme, err := os.ReadFile(fullPathReadme)
	if err != nil {
		fmt.Printf("Failed to parse README file at %s: %v\n", fullPathReadme, err)
		return err
	}

	replaceLinuxAmd64 := fmt.Sprintf("[Linux AMD64 - %s](%s)", builds.BuildDate, builds.LinuxAmd64)
	var expLinuxAmd64 = regexp.MustCompile(regexReadmeLinuxAmd64)
	replace1 := expLinuxAmd64.ReplaceAllString(string(readme), replaceLinuxAmd64)

	replaceDarwinAmd64 := fmt.Sprintf("[Darwin AMD64 - %s](%s)", builds.BuildDate, builds.DarwinAmd64)
	var expDarwinAmd64 = regexp.MustCompile(regexReadmeDarwinAmd64)
	replace2 := expDarwinAmd64.ReplaceAllString(replace1, replaceDarwinAmd64)

	replaceWindowsAmd64 := fmt.Sprintf("[Windows AMD64 - %s](%s)", builds.BuildDate, builds.WindowsAmd64)
	var expWindowsAmd64 = regexp.MustCompile(regexReadmeWindowsAmd64)
	replaceFinal := expWindowsAmd64.ReplaceAllString(replace2, replaceWindowsAmd64)

	err = os.Remove(fullPathReadme)
	if err != nil {
		fmt.Printf("Remove failed: %v\n", err)
		return err
	}
	err = os.WriteFile(fullPathReadme, []byte(replaceFinal), 0644)
	if err != nil {
		fmt.Printf("WriteFile failed: %v\n", err)
		return err
	}

	return nil
}

func postDiscussion(token string, builds *Builds) error {
	titleTemplate := "Daily Development Build for %s"
	bodyTemplate := "Downloadable assets:\\n\\n" +
		"Linux AMD64 - %s\\n" +
		"Darwin AMD64 - %s\\n" +
		"Windows AMD64 - %s\\n"
	jsonTemplate := "{\n" +
		"\"query\": \"mutation {   createDiscussion(input: {repositoryId: \\\"MDEwOlJlcG9zaXRvcnkzMDM4MDIzMzI\\\", categoryId: \\\"DIC_kwDOEhun3M4CA9pd\\\", body: \\\"%s\\\", title: \\\"%s\\\"}) {     discussion {       id     }   } }\"\n" +
		"}\n"

	url := "https://api.github.com/graphql"
	titleStr := fmt.Sprintf(titleTemplate, builds.BuildDate)
	bodyStr := fmt.Sprintf(bodyTemplate, builds.LinuxAmd64, builds.DarwinAmd64, builds.WindowsAmd64)
	authStr := fmt.Sprintf("token %s", token)
	jsonStr := fmt.Sprintf(jsonTemplate, bodyStr, titleStr)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*1200)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		fmt.Printf("http.NewRequest failed. Err: %v\n", err)
		return err
	}
	req.Header.Set("Authorization", authStr)
	req.Header.Set("Content-Type", "application/json")

	/* #nosec G402 */
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("client.Do failed. Err: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	fmt.Printf("StatusCode: %d\n", resp.StatusCode)
	fmt.Printf("Status: %s\n", resp.Status)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("resp.StatusCode failed: %d\n", resp.StatusCode)
		return errors.New(resp.Status)
	}

	return nil
}

func uploadAndUpdateRss(tag, ghToken, gcpToken string) error {
	builds, err := uploadAll(tag, gcpToken)
	if err != nil {
		fmt.Printf("uploadAll failed: %v\n", err)
		return err
	}

	err = updateReadme(builds)
	if err != nil {
		fmt.Printf("updateReadme failed: %v\n", err)
		return err
	}

	err = postDiscussion(ghToken, builds)
	if err != nil {
		fmt.Printf("postDiscussion failed: %v\n", err)
		return err
	}

	return nil
}

func main() {
	var tag string
	flag.StringVar(&tag, "tag", "", "The release tag to add")

	var gcpToken string
	if v := os.Getenv("GCP_BUCKET_SA"); v != "" {
		gcpToken = v
	}
	var ghToken string
	if v := os.Getenv("GITHUB_TOKEN"); v != "" {
		ghToken = v
	}

	flag.Parse()

	if tag == "" {
		fmt.Printf("A tag must be provided\n")
		os.Exit(1)
	}
	if gcpToken == "" {
		fmt.Printf("A GCP token must be provided\n")
		os.Exit(1)
	}
	if ghToken == "" {
		fmt.Printf("A GitHub token must be provided\n")
		os.Exit(1)
	}

	err := uploadAndUpdateRss(tag, ghToken, gcpToken)
	if err != nil {
		fmt.Printf("uploadAndUpdateRss failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Uploaded daily build succeeded\n")
	os.Exit(0)
}
