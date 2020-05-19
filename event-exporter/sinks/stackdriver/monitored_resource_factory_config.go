package stackdriver

import (
	"cloud.google.com/go/compute/metadata"
	"fmt"
	"github.com/GoogleCloudPlatform/k8s-stackdriver/event-exporter/utils"
	"github.com/golang/glog"
	"strings"
)

type monitoredResourceFactoryConfig struct {
	resourceModel resourceModelVersion
	clusterName   string
	location      string
	projectID     string
}

func newMonitoredResourceFactoryConfig(resourceModelVersion string) (*monitoredResourceFactoryConfig, error) {
	clusterName, err := utils.GetClusterName()
	if err != nil {
		glog.Warningf("'cluster-name' label is not specified on the VM, defaulting to the empty value")
		clusterName = ""
	}
	clusterName = strings.TrimSpace(clusterName)

	projectID, err := utils.GetProjectID()
	if err != nil {
		return nil, fmt.Errorf("failed to get project id: %v", err)
	}

	location, err := utils.GetClusterLocation()
	location = strings.TrimSpace(location)
	if err != nil || location == "" {
		glog.Warningf("Failed to retrieve cluster location, falling back to local zone: %s", err)
		location, err = metadata.Zone()
		if err != nil {
			return nil, fmt.Errorf("error while getting cluster location: %v", err)
		}
	}

	sdResourceModel := getResourceModelVersion(resourceModelVersion)
	return &monitoredResourceFactoryConfig{
		resourceModel: sdResourceModel,
		clusterName:   clusterName,
		location:      location,
		projectID:     projectID,
	}, nil
}

func getResourceModelVersion(model string) resourceModelVersion {
	if resourceModelVersion(model) == newTypes {
		return newTypes
	}
	return oldTypes
}
