package model

import (
	"fmt"
	"log"
	"strconv"

	"github.com/fabricator/fabricate/resource"
	yamlV2 "gopkg.in/yaml.v2"
)

type IDeployment interface {
	UpdateDeploymentYAML(i int, replicas int, lastSessionId int) string
}

type Deployment struct {
	resource resource.Deployment
}

func InitDeployment(resource resource.Deployment) IDeployment {
	deployment := new(Deployment)
	deployment.resource = resource
	return deployment
}

func (d *Deployment) UpdateDeploymentYAML(i int, replicas int, sid int) string {
	err := yamlV2.Unmarshal([]byte(resource.DeploymentYAML), &d.resource)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	depName := fmt.Sprint(fmt.Sprint("sid-", sid), fmt.Sprint("-deployment-", i))
	d.resource.Metadata.Name = depName
	d.resource.Metadata.Labels.Sid = strconv.Itoa(sid)
	if d.resource.Spec.Replicas != replicas {
		d.resource.Spec.Replicas = replicas
	}
	d.resource.Spec.Selector.MatchLabels.App = fmt.Sprint("svc-", i)
	d.resource.Spec.Template.Metadata.Labels.App = fmt.Sprint("svc-", i)
	d.resource.Spec.Template.Metadata.Labels.Sid = strconv.Itoa(sid)
	d.resource.Spec.Template.Spec.Containers[0].Env[0].ValueFrom.ConfigMapKeyRef.Name =
		fmt.Sprint(fmt.Sprint("sid-", sid), fmt.Sprint("-configmap-", i))

	data, err := yamlV2.Marshal(&d.resource)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("Creating %s ...\n", depName)
	return string(data)
}
