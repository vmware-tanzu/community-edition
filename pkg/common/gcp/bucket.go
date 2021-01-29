// Copyright 2020 VMware Tanzu Community Edition contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package gcp

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/adrg/xdg"
	"google.golang.org/api/option"
	klog "k8s.io/klog/v2"
)

// NewBucket generates a Bucket object
func NewBucket(byConfig []byte) (*Bucket, error) {

	cfg, err := InitBucketConfig(byConfig)
	if err != nil {
		klog.Errorf("InitBucketConfig failed. Err: %v", err)
		return nil, err
	}

	bucket := &Bucket{
		config:             cfg,
		remoteMetadataFile: filepath.Join(cfg.MetadataDirectory, cfg.VersionTag, cfg.MetadataFileName),
		localMetadataFile:  filepath.Join(xdg.DataHome, "tanzu-repository", cfg.MetadataDirectory, cfg.VersionTag, cfg.MetadataFileName),
		remoteReleaseFile:  filepath.Join(cfg.MetadataDirectory, cfg.ReleasesFileName),
	}

	klog.V(4).Infof("remoteMetadataFile = %s", bucket.remoteMetadataFile)
	klog.V(4).Infof("localMetadataFile = %s", bucket.localMetadataFile)
	klog.V(4).Infof("remoteReleaseFile = %s", bucket.remoteReleaseFile)

	return bucket, nil
}

func (b *Bucket) getBucket(ctx context.Context) (*storage.BucketHandle, error) {
	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		return nil, ErrBucketConnect
	}
	bkt := client.Bucket(b.config.TceBucket)
	if bkt == nil {
		return nil, ErrBucketConnect
	}
	return bkt, nil
}

// FetchMetadata grabs all the metadata relating to extensions
func (b *Bucket) fetchLocalMetadata() ([]byte, error) {
	klog.V(4).Info("Calling fetchLocalMetadata...")

	byFile, err := ioutil.ReadFile(b.localMetadataFile)
	if err != nil {
		klog.Errorf("Open failed. Err:", err)
		return nil, err
	}

	klog.V(2).Info("fetchLocalMetadata succeeded")
	klog.V(6).Info("Data:")
	klog.V(6).Info(string(byFile))

	return byFile, nil
}

// FetchMetadata grabs all the metadata relating to extensions
func (b *Bucket) fetchRemoteMetadata() ([]byte, error) {
	klog.V(4).Info("Calling fetchRemoteMetadata...")

	ctx := context.Background()

	bkt, err := b.getBucket(ctx)
	if err != nil {
		klog.Errorf("getBucket failed. Err: %v", err)
		return nil, err
	}

	obj := bkt.Object(b.remoteMetadataFile)
	if obj == nil {
		klog.Errorf("Object failed: %s", b.remoteMetadataFile)
		return nil, ErrBucketObject
	}

	r, err := obj.NewReader(ctx)
	if err != nil {
		klog.Errorf("NewReader failed. Err: %v", err)
		return nil, ErrBucketDownload
	}
	defer r.Close()

	byFile, err := ioutil.ReadAll(r)
	if err != nil {
		klog.Errorf("ReadAll failed. Err: ", err)
		return nil, err
	}

	//saving file locally
	toDir := filepath.Dir(b.localMetadataFile)

	err = os.MkdirAll(toDir, 0755)
	if err != nil {
		klog.Errorf("MkdirAll failed. Err: %v", err)
		return nil, err
	}

	file, err := os.Create(b.localMetadataFile)
	if err != nil {
		klog.Errorf("Create file failed. Err: %v", err)
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(byFile))
	if err != nil {
		klog.Errorf("Copy bits failed. Err: %v", err)
		return nil, err
	}
	//saving file locally

	klog.V(2).Info("fetchRemoteMetadata succeeded")
	klog.V(6).Info("Data:")
	klog.V(6).Info(string(byFile))

	return byFile, nil
}

// FetchMetadata grabs all the metadata relating to extensions
func (b *Bucket) FetchMetadata() ([]byte, error) {
	if _, err := os.Stat(b.localMetadataFile); os.IsNotExist(err) {
		klog.V(2).Infof("FetchMetadata remote")
		return b.fetchRemoteMetadata()
	}

	klog.V(2).Infof("FetchMetadata local")
	return b.fetchLocalMetadata()
}

// FetchRelease grabs all the metadata relating to releases
func (b *Bucket) FetchRelease() ([]byte, error) {
	ctx := context.Background()

	bkt, err := b.getBucket(ctx)
	if err != nil {
		klog.Errorf("getBucket failed. Err: %v", err)
		return nil, err
	}

	obj := bkt.Object(b.remoteReleaseFile)
	if obj == nil {
		klog.Errorf("Object failed: %s", b.remoteReleaseFile)
		return nil, ErrBucketObject
	}

	r, err := obj.NewReader(ctx)
	if err != nil {
		klog.Errorf("NewReader failed. Err: %v", err)
		return nil, ErrBucketDownload
	}
	defer r.Close()

	byFile, err := ioutil.ReadAll(r)
	if err != nil {
		klog.Errorf("ReadAll failed. Err: ", err)
		return nil, err
	}

	klog.V(4).Info("remoteReleaseFile succeeded")
	klog.V(6).Info("Data:")
	klog.V(6).Info(string(byFile))

	return byFile, nil
}

// Reset all directories to "factory"
func (b *Bucket) Reset() error {

	klog.V(2).Infof("bucket.Reset")
	metadatDirectory := filepath.Join(xdg.DataHome, "tanzu-repository", b.config.MetadataDirectory)
	klog.V(2).Infof("metadatDirectory = %s", metadatDirectory)

	err := os.RemoveAll(metadatDirectory)
	if err != nil {
		klog.Errorf("MkdirAll failed. Err: %v", err)
		return err
	}

	klog.V(2).Infof("bucket.Reset succeeded")
	return nil
}
