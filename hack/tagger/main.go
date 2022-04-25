// Copyright 2022 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-containerregistry/pkg/crane"
	log "github.com/sirupsen/logrus"
)

// yeah global clients are bad, but this is a script
var myClient = &http.Client{Timeout: 30 * time.Second}

// Artifacts is the returned object
type Artifacts []struct {
	AdditionLinks struct {
		BuildHistory struct {
			Absolute bool   `json:"absolute"`
			Href     string `json:"href"`
		} `json:"build_history"`
		Vulnerabilities struct {
			Absolute bool   `json:"absolute"`
			Href     string `json:"href"`
		} `json:"vulnerabilities"`
	} `json:"addition_links"`
	Digest     string `json:"digest"`
	ExtraAttrs struct {
		Architecture string      `json:"architecture"`
		Author       interface{} `json:"author"`
		Created      time.Time   `json:"created"`
		Os           string      `json:"os"`
	} `json:"extra_attrs"`
	Icon              string      `json:"icon"`
	ID                int         `json:"id"`
	Labels            interface{} `json:"labels"`
	ManifestMediaType string      `json:"manifest_media_type"`
	MediaType         string      `json:"media_type"`
	ProjectID         int         `json:"project_id"`
	PullTime          time.Time   `json:"pull_time"`
	PushTime          time.Time   `json:"push_time"`
	References        interface{} `json:"references"`
	RepositoryID      int         `json:"repository_id"`
	Size              int         `json:"size"`
	Tags              []Tag       `json:"tags"`
	Type              string      `json:"type"`
}

type Tag struct {
	ArtifactID   int       `json:"artifact_id"`
	ID           int       `json:"id"`
	Immutable    bool      `json:"immutable"`
	Name         string    `json:"name"`
	PullTime     time.Time `json:"pull_time"`
	PushTime     time.Time `json:"push_time"`
	RepositoryID int       `json:"repository_id"`
	Signed       bool      `json:"signed"`
}

// harbor endpoint used to get all artifacts for a project.
// it will list untagged artifacts
const (
	urlTemp = "https://projects.registry.vmware.com/api/v2.0/projects/tce/repositories/%s/artifacts"
)

var (
	// adding a new package? change this, recompile, and re-run
	projects = []string{
		"app-toolkit",
		"cartographer",
		"cartographer-catalog",
		"cert-injection-webhook",
		"cert-manager",
		"contour",
		"external-dns",
		"fluent-bit",
		"fluxcd-kustomize-controller",
		"fluxcd-source-controller",
		"gatekeeper",
		"grafana",
		"harbor",
		"knative-serving",
		"kpack",
		"local-path-storage",
		"multus-cni",
		"prometheus",
		"velero",
		"whereabouts",
	}
)

func main() {
	// setup logger
	formatter := &log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	log.SetFormatter(formatter)
	log.Infoln("starting")

	RunTagger()
}

// RunTagger is the process of finding artifact's tags and determining if the required
// SHA-based tag is required
func RunTagger() {
	log.Infoln("Tagger check started")
	// for every projct, look-up its artifacts
	for _, project := range projects {
		artifactList := &Artifacts{}
		artifactUrl := fmt.Sprintf(urlTemp, project)
		err := ListAllArtifacts(artifactUrl, artifactList)
		if err != nil {
			log.Errorf("Failed to get artifacts for %s. Skipping. Reason: %s", artifactUrl, err.Error())
			break
		}

		// for each artifact, check to see if its missing the required SHA tag
		for _, a := range *artifactList {
			log.Debugf("(%s) evaluating: %s", project, a.Digest)

			// if the digest is malformed, error log and skip this artifact
			shaSplit := strings.Split(a.Digest, ":")
			if len(shaSplit) < 2 {
				log.Errorf("Failed to parse SHA for digest. Skipping. Origianl Value: %s", a.Digest)
				break
			}

			// if the sha value is less than 10 characters, its invalid. Error log and skip this artifact
			SHAValue := shaSplit[1]
			if len(SHAValue) < 10 {
				log.Errorf("SHA Value was invalid. Original Value: %s; Parsed value: %s", a.Digest, SHAValue)
				break
			}

			// the required tag is the first 10 characters of the SHA value
			requiredTag := SHAValue[0:10]
			imgUrl := fmt.Sprintf("%s/%s@%s", "projects.registry.vmware.com/tce", project, a.Digest)

			// check for the requiredTag in the artifacts tags
			// if its missing, add it
			if !SHATagPresent(requiredTag, a.Tags) {
				log.Infof("Attempting to add SHA-based tag %s to %s", requiredTag, imgUrl)

				// add the requiredTag
				err := AddSHATag(imgUrl, requiredTag)
				if err != nil {
					log.Errorf("Failed to add SHA-based tag to %s. Reason: %s", imgUrl, err)
					break
				}

				log.Infof("Added SHA-based tag %s to %s", requiredTag, imgUrl)

			} else {

				log.Debugf("(%s) Required SHA tag (%s) exists on %s", project, requiredTag, imgUrl)

			}

		}
	}
	log.Infoln("Tagger check finished")
}

// ListAllArtifacts retrieves every known artifact (tagged or not) in a repository
func ListAllArtifacts(url string, a *Artifacts) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(a)
}

// SHATagPresent returns true if the first 10 characters of the SHA are a tag
// on the OCI artifact, otherwise it returns false.
func SHATagPresent(requiredTag string, existingTags []Tag) bool {
	requiredTagExists := false
	for _, t := range existingTags {
		if !requiredTagExists {
			requiredTagExists = strings.EqualFold(t.Name, requiredTag)
		}
		log.Debugf("tag: %s", t.Name)
	}
	return requiredTagExists
}

// AddSHATag adds a new image tag to an artifact
func AddSHATag(imageUrl string, extraTag string) error {
	err := crane.Tag(imageUrl, extraTag, crane.WithContext(context.TODO()))
	if err != nil {
		return err
	}
	return nil
}
