// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"cloud.google.com/go/storage"
)

const (
	// GH discussion related consts
	githubDiscussionURI string = "https://api.github.com/graphql"

	// GCP related consts
	projectID  string = "vmware-cna"
	bucketName string = "tce-cli-plugins-staging"
	gcpFile    string = "gcptokenfile.json"

	// README.md related consts
	fullPathReadme          string = "../../README.md"
	regexReadmeLinuxAmd64   string = "\\[Linux AMD64.*\\]\\(.*\\)"
	regexReadmeDarwinAmd64  string = "\\[Darwin AMD64.*\\]\\(.*\\)"
	regexReadmeWindowsAmd64 string = "\\[Windows AMD64.*\\]\\(.*\\)"
)

type Result struct {
	Data struct {
		Repository struct {
			Discussions struct {
				Nodes []struct {
					ID    string `json:"id"`
					Title string `json:"title"`
				} `json:"nodes"`
			} `json:"discussions"`
		} `json:"repository"`
	} `json:"data"`
}

type Builds struct {
	BuildYear    int
	BuildMonth   string
	BuildDate    string
	DarwinAmd64  string
	DarwinArm64  string
	WindowsAmd64 string
	LinuxAmd64   string
	PreviousHash string
	CurrentHash  string
}

type ClientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

// the upload function
func uploadAll(tag, gcpToken, previous, current string) (*Builds, error) {
	err := os.WriteFile(gcpFile, []byte(gcpToken), 0644)
	if err != nil {
		fmt.Printf("WriteFile failed: %v\n", err)
		return nil, err
	}

	myNow := time.Now()
	year, month, _ := myNow.Date()
	buildDate := myNow.Format("2006-01-02")

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
		BuildYear:    year,
		BuildMonth:   month.String(),
		BuildDate:    buildDate,
		LinuxAmd64:   urlLinuxAmd64,
		WindowsAmd64: urlWindowsAmd64,
		DarwinAmd64:  urlDarwinAmd64,
		PreviousHash: previous,
		CurrentHash:  current,
	}

	return builds, nil
}

// the readme update function
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

// the post discussion function
func doDailyBuildNotification(token string, builds *Builds) error {
	discussID, err := getDiscussionID(token, builds)
	if err != nil {
		fmt.Printf("Discussion for this month doesnt exist. Create it!")

		err = createDiscussion(token, builds)
		if err != nil {
			fmt.Printf("createDiscussion failed: %v\n", err)
			return err
		}

		discussID, err = getDiscussionID(token, builds)
		if err != nil {
			fmt.Printf("getDiscussionID failed: %v\n", err)
			return err
		}
	}
	fmt.Printf("ID: %s\n", discussID)

	err = postComment(token, discussID, builds)
	if err != nil {
		fmt.Printf("postComment failed: %v\n", err)
		return err
	}

	return nil
}

// helper bucket functions
func (c *ClientUploader) UploadFile(filename string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*1200)
	defer cancel()

	// open the src file
	uploadFile := filepath.Join("..", "..", "release", filename)
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

// helper discussion functions start
func createDiscussion(token string, builds *Builds) error {
	titleTemplate := "%s %d: Daily Development Builds"
	bodyTemplate := "# Downloadable assets:\\n\\n" +
		"You can find published daily development builds for %s %d here. These builds are intentionally unsigned " +
		"as they aren't apart of any supported release. These builds are for development/test purposes only. " +
		"To get the latest, sort the replies in this discussion from newest to oldest.\\n\\n"

	jsonTemplate := "{\n" +
		"\"query\": \"mutation {   createDiscussion(input: {repositoryId: \\\"MDEwOlJlcG9zaXRvcnkzMDM4MDIzMzI\\\", categoryId: \\\"DIC_kwDOEhun3M4COOTN\\\", body: \\\"%s\\\", title: \\\"%s\\\"}) {     discussion {       id     }   } }\"\n" +
		"}\n"

	titleStr := fmt.Sprintf(titleTemplate, builds.BuildMonth, builds.BuildYear)
	bodyStr := fmt.Sprintf(bodyTemplate, builds.BuildMonth, builds.BuildYear)
	authStr := fmt.Sprintf("token %s", token)
	jsonStr := fmt.Sprintf(jsonTemplate, bodyStr, titleStr)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*1200)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", githubDiscussionURI, bytes.NewBuffer([]byte(jsonStr)))
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

func postComment(token, discussID string, builds *Builds) error {
	bodyTemplate := "# Downloadable assets:\\n\\n" +
		"Linux AMD64 - %s\\n" +
		"Darwin AMD64 - %s\\n" +
		"Windows AMD64 - %s\\n\\n"
	if builds.PreviousHash != builds.CurrentHash {
		changelog := fmt.Sprintf("https://github.com/vmware-tanzu/community-edition/compare/%s...%s\\n", builds.PreviousHash, builds.CurrentHash)

		bodyTemplate += "# Full Changelog From Last Daily Build\\n\\n"
		bodyTemplate += changelog
	}

	jsonTemplate := "{\n" +
		"\"query\": \"mutation {   addDiscussionComment(input: {discussionId: \\\"%s\\\", body: \\\"%s\\\"}) { comment {       id     }   } }\"\n" +
		"}\n"

	bodyStr := fmt.Sprintf(bodyTemplate, builds.LinuxAmd64, builds.DarwinAmd64, builds.WindowsAmd64)
	authStr := fmt.Sprintf("token %s", token)
	jsonStr := fmt.Sprintf(jsonTemplate, discussID, bodyStr)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*1200)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", githubDiscussionURI, bytes.NewBuffer([]byte(jsonStr)))
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

func getDiscussionID(token string, builds *Builds) (string, error) {
	authStr := fmt.Sprintf("token %s", token)
	jsonStr := "{\n" +
		"\"query\": \"query {   repository(owner: \\\"vmware-tanzu\\\", name: \\\"community-edition\\\") {     discussions(first: 25) {       nodes {         id title       }     }   } }\"\n" +
		"}\n"

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*1200)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", githubDiscussionURI, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		fmt.Printf("http.NewRequest failed. Err: %v\n", err)
		return "", err
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
		return "", err
	}
	defer resp.Body.Close()

	fmt.Printf("StatusCode: %d\n", resp.StatusCode)
	fmt.Printf("Status: %s\n", resp.Status)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("resp.StatusCode failed: %d\n", resp.StatusCode)
		return "", errors.New(resp.Status)
	}

	byFile, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ReadAll failed. Err: %v\n", err)
		return "", err
	}

	fmt.Printf("DATA:\n\n%s\n\n", string(byFile))

	result := Result{}
	err = json.Unmarshal(byFile, &result)
	if err != nil {
		fmt.Printf("json.Unmarshal failed. Err: %v\n", err)
		return "", err
	}

	titleMatch := fmt.Sprintf("%s %d:", builds.BuildMonth, builds.BuildYear)
	for _, article := range result.Data.Repository.Discussions.Nodes {
		if strings.Index(article.Title, titleMatch) == 0 {
			fmt.Printf("ID: %s\n", article.ID)
			return article.ID, nil
		}
	}

	fmt.Printf("ID NOT FOUND\n")
	return "", errors.New("ID not found")
}

// do it all
func doWork(ghToken, gcpToken, tag, previous, current string) error {
	builds, err := uploadAll(tag, gcpToken, previous, current)
	if err != nil {
		fmt.Printf("uploadAll failed: %v\n", err)
		return err
	}

	err = updateReadme(builds)
	if err != nil {
		fmt.Printf("updateReadme failed: %v\n", err)
		return err
	}

	err = doDailyBuildNotification(ghToken, builds)
	if err != nil {
		fmt.Printf("doDailyBuildNotification failed: %v\n", err)
		return err
	}

	return nil
}

func main() {
	var tag string
	flag.StringVar(&tag, "tag", "", "Build tag")

	var previous string
	flag.StringVar(&previous, "previous", "", "Previous hash")

	var current string
	flag.StringVar(&current, "current", "", "Current hash")

	flag.Parse()

	var gcpToken string
	if v := os.Getenv("GCP_BUCKET_SA"); v != "" {
		gcpToken = v
	}
	var ghToken string
	if v := os.Getenv("GITHUB_TOKEN"); v != "" {
		ghToken = v
	}

	if tag == "" {
		fmt.Printf("A tag must be provided\n")
		os.Exit(1)
	}
	if previous == "" {
		fmt.Printf("A previous hash must be provided\n")
		os.Exit(1)
	}
	if current == "" {
		fmt.Printf("A current hash must be provided\n")
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

	err := doWork(ghToken, gcpToken, tag, previous, current)
	if err != nil {
		fmt.Printf("doWork failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Uploaded daily build succeeded\n")
	os.Exit(0)
}
