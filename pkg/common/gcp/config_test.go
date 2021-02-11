package gcp

import "testing"

func TestInitBucketConfigEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("InitBucket should not panic: %v", r)
		}
	}()

	InitBucketConfig([]byte(""))
}
