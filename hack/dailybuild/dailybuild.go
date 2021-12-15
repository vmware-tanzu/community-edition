// Copyright 2021 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gorilla/feeds"
)

const (
	projectID  string = "vmware-cna"
	bucketName string = "tce-cli-plugins-staging"
	gcpFile    string = "gcptokenfile.json"
	rssFile    string = "../../rss.xml"
)

type Builds struct {
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

func uploadTarball(platform, arch, tag, extension string) (string, error) {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", gcpFile)

	client, err := storage.NewClient(context.Background())
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return "", err
	}

	assetFilename := fmt.Sprintf("tce-%s-%s-%s.%s", platform, arch, tag, extension)
	uploadPathTmp := filepath.Join("build-daily", time.Now().Format("2006-01-02"))
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
	// upload
	err := os.WriteFile(gcpFile, []byte(gcpToken), 0644)
	if err != nil {
		fmt.Printf("WriteFile failed: %v\n", err)
		return nil, err
	}

	urlLinuxAmd64, err := uploadTarball("linux", "amd64", tag, "tar.gz")
	if err != nil {
		fmt.Printf("uploadLinux failed: %v\n", err)

		errFile := os.Remove(gcpFile)
		if err != nil {
			fmt.Printf("Remove failed: %v\n", errFile)
		}

		return nil, err
	}

	urlWindowsAmd64, err := uploadTarball("windows", "amd64", tag, "zip")
	if err != nil {
		fmt.Printf("uploadWindows failed: %v\n", err)

		errFile := os.Remove(gcpFile)
		if err != nil {
			fmt.Printf("Remove failed: %v\n", errFile)
		}

		return nil, err
	}

	urlDarwinAmd64, err := uploadTarball("darwin", "amd64", tag, "tar.gz")
	if err != nil {
		fmt.Printf("uploadDarwinAmd64 failed: %v\n", err)

		errFile := os.Remove(gcpFile)
		if err != nil {
			fmt.Printf("Remove failed: %v\n", errFile)
		}

		return nil, err
	}

	_, err = uploadTarball("darwin", "arm64", tag, "tar.gz")
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
		LinuxAmd64:   urlLinuxAmd64,
		WindowsAmd64: urlWindowsAmd64,
		DarwinAmd64:  urlDarwinAmd64,
	}

	return builds, nil
}

func createRssFeed(builds *Builds) error {
	// update rss
	now, err := time.Parse(time.RFC3339, "2013-01-16T21:52:35-05:00")
	if err != nil {
		fmt.Printf("time.Parse failed: %v\n", err)
		return err
	}
	tz := time.FixedZone("PST", -8*60*60)
	now = now.In(tz)

	dateDailyBuild := fmt.Sprintf("Daily Build %s", time.Now().Format("2006-01-02"))
	buildMessage := fmt.Sprintf("<p>You can find the daily build here:<br><a href=%q>%s</a><br><a href=%q>%s</a><br><a href=%q>%s</a></p>",
		builds.LinuxAmd64, builds.LinuxAmd64, builds.WindowsAmd64, builds.WindowsAmd64, builds.DarwinAmd64, builds.DarwinAmd64)

	feed := &feeds.Feed{
		Title:       "Tanzu Community Edition daily build",
		Link:        &feeds.Link{Href: "https://github.com/vmware-tanzu/community-edition"},
		Description: "Tanzu Community Edition daily build",
		Author:      &feeds.Author{Name: "tce-automation", Email: "tce-core-eng@vmware.com"},
		Created:     now,
		Copyright:   "This work is copyright © VMware",
	}

	feed.Items = []*feeds.Item{
		{
			Title:       dateDailyBuild,
			Link:        &feeds.Link{Href: "https://github.com/vmware-tanzu/community-edition"},
			Description: dateDailyBuild,
			Author:      &feeds.Author{Name: "tce-automation", Email: "tce-core-eng@vmware.com"},
			Created:     now,
			Content:     buildMessage,
		},
	}

	rss, err := feed.ToRss()
	if err != nil {
		fmt.Printf("unexpected error encoding RSS: %v", err)
		return err
	}

	errFile := os.Remove(rssFile)
	if err != nil {
		fmt.Printf("Remove failed: %v\n", errFile)
	}
	err = os.WriteFile(rssFile, []byte(rss), 0666)
	if err != nil {
		fmt.Printf("WriteFile failed: %v\n", err)
		return err
	}

	return nil
}

func uploadAndUpdateRss(tag, gcpToken string) error {
	builds, err := uploadAll(tag, gcpToken)
	if err != nil {
		fmt.Printf("uploadAll failed: %v\n", err)
		return err
	}

	err = createRssFeed(builds)
	if err != nil {
		fmt.Printf("createRssFeed failed: %v\n", err)
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

	flag.Parse()

	if tag == "" {
		fmt.Printf("A tag must be provided\n")
		os.Exit(1)
	}
	if gcpToken == "" {
		fmt.Printf("A token must be provided\n")
		os.Exit(1)
	}

	err := uploadAndUpdateRss(tag, gcpToken)
	if err != nil {
		fmt.Printf("uploadAndUpdateRss failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Uploaded daily build succeeded\n")
	os.Exit(0)
}
