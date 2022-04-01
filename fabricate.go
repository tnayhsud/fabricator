package fabricate

import (
	"context"

	"github.com/fabricator/fabricate/client"
	"github.com/fabricator/fabricate/model"
	"github.com/fabricator/fabricate/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Fabricate interface {
	Fabricate()
}

// type fabricate implementing Fabricate interface
type fabricate struct {
	config   *rest.Config
	count    int
	replicas int
}

/*
   func    New
   desc:   returns a pointer of type fabricate
   param   config *rest.Config: a pointer of type rest.Config to authenticate with k8s api-server
   param   count int: Requires an integer to specify number of resources i.e. number of config maps,
           deployments and services
   param   replicas int: Requires an integer to specify number of pods to create per deployment
           returns a pointer of type fabricate
*/
func New(config *rest.Config, count int, replicas int) Fabricate {
	return &fabricate{config, count, replicas}
}

func CreateStaticClient(config *rest.Config) *kubernetes.Clientset {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return client
}

/*
   method  Fabricate
   desc:   creates new SSAClient with the context and kube config provided and invokes ServerSideApply
           with count and replicas
*/
func (f *fabricate) Fabricate() {
	context := context.Background()
	namespace := model.InitNamespace(resource.Namespace{})
	configMap := model.InitConfigMap(resource.ConfigMap{})
	deployment := model.InitDeployment(resource.Deployment{})
	service := model.InitService(resource.Service{})
	clientUtil := client.InitClientUtil()
	client := client.InitFabricateService(f.config, namespace, configMap, deployment, service, CreateStaticClient(f.config), clientUtil)
	client.ServerSideApply(context, f.count, f.replicas)
}
