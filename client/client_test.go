package client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/fabricator/fabricate/model"
	"github.com/fabricator/fabricate/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/api/meta"

	yaml "gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	yamlV2 "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	dynamicfakeclient "k8s.io/client-go/dynamic/fake"
	testclient "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type MockClientUtil struct {
	mock.Mock
}

func (m *MockClientUtil) ResourceMapping(cfg *rest.Config, gvk *schema.GroupVersionKind) *meta.RESTMapping {
	ret := m.Mock.Called(cfg, gvk)
	fmt.Println("mocked ResourceMapping method called")
	return ret.Get(0).(*meta.RESTMapping)
}

func (m *MockClientUtil) DynamicResourceInterface(cfg *rest.Config, mapping *meta.RESTMapping,
	obj *unstructured.Unstructured) dynamic.ResourceInterface {
	ret := m.Mock.Called(cfg, mapping, obj)
	fmt.Println("mocked DynamicResourceInterface method called")
	return ret.Get(0).(dynamic.ResourceInterface)
}

func TestDecodeGVKAndCreate(t *testing.T) {
	config := config()
	serializer := yamlV2.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	dyn := dynamicfakeclient.NewSimpleDynamicClient(runtime.NewScheme())
	mockClientUtil := new(MockClientUtil)

	for _, tt := range []struct {
		group    string
		version  string
		resource string
		kind     string
		yaml     string
		result   string
	}{
		{"", "v1", "configmaps", "ConfigMap", resource.ConfigYAML, "ConfigMap"},
		{"apps", "v1", "deployment", "Deployment", resource.DeploymentYAML, "Deployment"},
		{"", "v1", "services", "Service", resource.ServiceYAML, "Service"},
	} {
		mapping, gvk := mapping(tt.group, tt.version, tt.resource, tt.kind)
		serializer.Decode([]byte(tt.yaml), nil, obj)
		dr := dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
		mockClientUtil.Mock.
			On("ResourceMapping", config, &gvk).Return(mapping).
			On("DynamicResourceInterface", config, mapping, obj).Return(dr)

		namespace := model.InitNamespace(resource.Namespace{})
		configMap := model.InitConfigMap(resource.ConfigMap{})
		deployment := model.InitDeployment(resource.Deployment{})
		service := model.InitService(resource.Service{})

		fabricateService := InitFabricateService(config, namespace, configMap, deployment, service, testclient.NewSimpleClientset(), mockClientUtil)
		res := fabricateService.DecodeGVKAndCreate(tt.yaml, context.Background(), config)
		kind, found, isString := unstructured.NestedString(res.UnstructuredContent(), "kind")
		fmt.Println(found, isString)
		data, _ := yaml.Marshal(res)
		fmt.Println(string(data))
		assert.Equal(t, kind, tt.result, "They should be equal")
		assert.NotEqual(t, kind, "", "They should be not equal")
	}
}

func config() *rest.Config {
	home, _ := os.UserHomeDir()
	kubeConfig := filepath.Join(home, ".kube", "config")
	config, _ := clientcmd.BuildConfigFromFlags("", kubeConfig)
	return config
}

func mapping(g string, v string, r string, k string) (*meta.RESTMapping, schema.GroupVersionKind) {
	mapping := new(meta.RESTMapping)
	var gvr schema.GroupVersionResource
	gvr.Group = g
	gvr.Version = v
	gvr.Resource = r
	mapping.Resource = gvr
	var gvk schema.GroupVersionKind
	gvk.Group = g
	gvk.Version = v
	gvk.Kind = k
	mapping.GroupVersionKind = gvk
	mapping.Scope = &CustomRestScope{name: "namespace"}
	return mapping, gvk
}
