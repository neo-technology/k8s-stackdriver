package utils

import (
	"cloud.google.com/go/compute/metadata"
	"fmt"
	"os"
)

func GetProjectID() (string, error) {
	return getProperty("PROJECT_ID", func() (string, error) {
		return metadata.ProjectID()
	})
}

func GetClusterName() (string, error) {
	return getProperty("CLUSTER_NAME", func() (string, error) {
		return metadata.InstanceAttributeValue("cluster-name")
	})
}

func GetClusterLocation() (string, error) {
	return getProperty("CLUSTER_LOCATION", func() (string, error) {
		return metadata.InstanceAttributeValue("cluster-location")
	})
}

func getProperty(key string, fetchPropFromMetadataFunc func() (string, error)) (string, error) {
	value, exists := os.LookupEnv(key)
	if exists {
		return value, nil
	}

	if !metadata.OnGCE() {
		return "", fmt.Errorf("not running on GCE and %s is not defined", key)
	}

	value, err := fetchPropFromMetadataFunc()
	if err != nil {
		return "", fmt.Errorf("failed to fetch property from metadata server: %v", err)
	}

	return value, nil
}
