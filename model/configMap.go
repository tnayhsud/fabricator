package model

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/fabricator/fabricate/resource"
	yamlV2 "gopkg.in/yaml.v2"
)

type IConfigMap interface {
	UpdateConfigYAML(i int, count int, lastSessionId int) string
}

type ConfigMap struct {
	resource resource.ConfigMap
}

func InitConfigMap(resource resource.ConfigMap) IConfigMap {
	configMap := new(ConfigMap)
	configMap.resource = resource
	return configMap
}

func (c *ConfigMap) UpdateConfigYAML(i int, count int, sid int) string {
	serviceList := createServiceList(count, sid)
	err := yamlV2.Unmarshal([]byte(resource.ConfigYAML), &c.resource)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	cmName := fmt.Sprint(fmt.Sprint("sid-", sid), fmt.Sprint("-configmap-", i))
	c.resource.Metadata.Name = cmName
	c.resource.Metadata.Labels.Sid = strconv.Itoa(sid)
	c.resource.Data.Services = strings.Join(append(serviceList[:i-1], serviceList[i:]...), ",")
	data, err := yamlV2.Marshal(&c.resource)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("Creating %s ...\n", cmName)
	return string(data)
}

func createServiceList(count int, lastSessionId int) []string {
	serviceList := make([]string, count)
	for i := 1; i <= count; i++ {
		serviceList[i-1] = fmt.Sprint(fmt.Sprint("sid-", lastSessionId), fmt.Sprint("-service-", i))
	}
	return serviceList
}
