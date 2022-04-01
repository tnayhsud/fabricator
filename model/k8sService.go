package model

import (
	"fmt"
	"log"
	"strconv"

	"github.com/fabricator/fabricate/resource"
	yamlV2 "gopkg.in/yaml.v2"
)

type IService interface {
	UpdateServiceYAML(i int, lastSessionId int) string
}

type Service struct {
	resource resource.Service
}

func InitService(resource resource.Service) IService {
	service := new(Service)
	service.resource = resource
	return service
}

func (s *Service) UpdateServiceYAML(i int, sid int) string {
	err := yamlV2.Unmarshal([]byte(resource.ServiceYAML), &s.resource)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	svcName := fmt.Sprint(fmt.Sprint("sid-", sid), fmt.Sprint("-service-", i))
	s.resource.Metadata.Name = svcName
	s.resource.Metadata.Labels.Sid = strconv.Itoa(sid)
	s.resource.Spec.Selector.App = fmt.Sprint("svc-", i)
	s.resource.Spec.Selector.Sid = strconv.Itoa(sid)

	data, err := yamlV2.Marshal(&s.resource)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("Creating %s ...\n", svcName)
	return string(data)
}
