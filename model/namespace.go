package model

import (
	"log"
	"strconv"

	"github.com/fabricator/fabricate/resource"
	yamlV2 "gopkg.in/yaml.v2"
)

type INamespace interface {
	UpdateNamespaceYAML(sid int) string
}

type Namespace struct {
	resource resource.Namespace
}

func InitNamespace(resource resource.Namespace) INamespace {
	namespace := new(Namespace)
	namespace.resource = resource
	return namespace
}

func (n *Namespace) UpdateNamespaceYAML(sid int) string {
	err := yamlV2.Unmarshal([]byte(resource.NamespaceYAML), &n.resource)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	n.resource.Metadata.Labels.SId = strconv.Itoa(sid)
	data, err := yamlV2.Marshal(&n.resource)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	return string(data)
}
