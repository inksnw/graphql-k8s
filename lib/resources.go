package lib

import (
	"github.com/phuslu/log"
	"strings"
)

type ResourceType struct {
	ResourceName string
	Group        string
	Version      string
	Resource     string
	Kind         string
}

func ResourceTypes() []ResourceType {
	apiGroups, err := ClientSet.ServerGroups()
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}
	prefix := "io.k8s.api"
	Types := make([]ResourceType, 0)
	for _, apiGroup := range apiGroups.Groups {
		for _, groupVersion := range apiGroup.Versions {
			resources, err := ClientSet.ServerResourcesForGroupVersion(groupVersion.GroupVersion)
			if err != nil {
				log.Fatal().Msgf("%s", err)
			}
			for _, resource := range resources.APIResources {
				//跳过子资源
				if strings.Contains(resource.Name, "/") {
					continue
				}
				var openApiDefinitionId string
				if apiGroup.Name == "" {
					list := []string{prefix, "core", groupVersion.Version, resource.Kind}
					openApiDefinitionId = strings.Join(list, ".")
				} else if apiGroup.Name == "apps" {
					list := []string{prefix, apiGroup.Name, groupVersion.Version, resource.Kind}
					openApiDefinitionId = strings.Join(list, ".")
				}
				if openApiDefinitionId != "" {
					ins := ResourceType{
						ResourceName: openApiDefinitionId,
						Group:        apiGroup.Name,
						Version:      groupVersion.Version,
						Resource:     resource.Name,
						Kind:         resource.Kind,
					}
					Types = append(Types, ins)
				}
			}
		}
	}
	return Types
}
