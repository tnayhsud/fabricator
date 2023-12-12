package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	testclient "k8s.io/client-go/kubernetes/fake"

	"github.com/fabricator/fabricate/model"
	"github.com/fabricator/fabricate/resource"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	yamlV2 "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type IFabricateService interface {
	ServerSideApply(ctx context.Context, count int, replicas int)
	DecodeGVKAndCreate(resourceYAML string, ctx context.Context, cfg *rest.Config) *unstructured.Unstructured
	CreateOrPatchResource(ctx context.Context, cfg *rest.Config, obj *unstructured.Unstructured, gvk *schema.GroupVersionKind) *unstructured.Unstructured
}

type FabricateService struct {
	cfg             *rest.Config
	namespace       model.INamespace
	configMap       model.IConfigMap
	deployment      model.IDeployment
	service         model.IService
	serializer      runtime.Serializer
	staticK8sClient interface{}
	clientUtil      IClientUtil
}

func InitFabricateService(cfg *rest.Config, namespace model.INamespace, configMap model.IConfigMap,
	deployment model.IDeployment, service model.IService, staticK8sClient interface{}, clientUtil IClientUtil) IFabricateService {
	fabricateService := new(FabricateService)
	fabricateService.cfg = cfg
	fabricateService.staticK8sClient = staticK8sClient
	fabricateService.namespace = namespace
	fabricateService.configMap = configMap
	fabricateService.deployment = deployment
	fabricateService.service = service
	fabricateService.serializer = yamlV2.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	fabricateService.clientUtil = clientUtil
	return fabricateService
}

/*
	method 	ServerSideApply
	desc 	initiates the process of resource creation/patching using SSA
	param 	count int
	param 	replicas int
*/
func (f *FabricateService) ServerSideApply(ctx context.Context, count int, replicas int) {
	var sid int
	switch f.staticK8sClient.(type) {
	case *kubernetes.Clientset:
		ns, err := f.staticK8sClient.(*kubernetes.Clientset).CoreV1().Namespaces().Get(ctx, "test", metav1.GetOptions{})
		sid = createNS(err, f, ctx, sid, ns)
	case testclient.Clientset:
		ns, err := f.staticK8sClient.(*testclient.Clientset).CoreV1().Namespaces().Get(ctx, "test", metav1.GetOptions{})
		sid = createNS(err, f, ctx, sid, ns)
	}
	for i := 1; i <= count; i++ {
		configMapYAML := f.configMap.UpdateConfigYAML(i, count, sid)
		f.DecodeGVKAndCreate(configMapYAML, ctx, f.cfg)
	}
	for i := 1; i <= count; i++ {
		deploymentYAML := f.deployment.UpdateDeploymentYAML(i, replicas, sid)
		f.DecodeGVKAndCreate(deploymentYAML, ctx, f.cfg)
	}
	for i := 1; i <= count; i++ {
		serviceYAML := f.service.UpdateServiceYAML(i, sid)
		f.DecodeGVKAndCreate(serviceYAML, ctx, f.cfg)
	}
}

func createNS(err error, f *FabricateService, ctx context.Context, sid int, ns *v1.Namespace) int {
	if err != nil {
		sid += 1
		fmt.Printf("Current session id: %d\n", sid)
		fmt.Println("Creating 'test' namepace...")
		f.DecodeGVKAndCreate(resource.NamespaceYAML, ctx, f.cfg)
		fmt.Println("'test' namespace created")
	} else {
		fmt.Println("Found 'test' namepace...")
		sid, _ = strconv.Atoi(ns.ObjectMeta.Labels["sid"])
		sid += 1
		fmt.Printf("Current session id: %d\n", sid)
		namespaceYAML := f.namespace.UpdateNamespaceYAML(sid)
		fmt.Println("Updating session-id in 'test' namespace")
		f.DecodeGVKAndCreate(namespaceYAML, ctx, f.cfg)
		fmt.Println("Updated session-id")
	}
	return sid
}

/*
	func 	decodeGVKAndCreate
	desc 	decodes the YAML string into unstructured struct which return group, version & kind
			for the resource defined in the YAMLand invokes createOrPatchResource
	param 	count int
	param 	replicas int
*/
func (f *FabricateService) DecodeGVKAndCreate(resourceYAML string, ctx context.Context, cfg *rest.Config) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	_, gvk, err := f.serializer.Decode([]byte(resourceYAML), nil, obj)
	if err != nil {
		fmt.Println("Error", err)
		return nil
	}
	return f.CreateOrPatchResource(ctx, cfg, obj, gvk)
}

/*
	func 	createOrPatchResource
	desc	Marshals resource obj to json and applies to k8s cluster with dynamicClient
	param 	ctx context.Context
	param 	cfg *rest.Config
	param 	obj *unstructured.Unstructured
	param 	gvk *schema.GroupVersionKind
*/
func (f *FabricateService) CreateOrPatchResource(ctx context.Context, cfg *rest.Config, obj *unstructured.Unstructured,
	gvk *schema.GroupVersionKind) *unstructured.Unstructured {
	// find preferred resource mapping for the gvk
	mapping := f.clientUtil.ResourceMapping(cfg, gvk)

	// obtain dynamic resource interface for the mapping
	dr := f.clientUtil.DynamicResourceInterface(cfg, mapping, obj)

	data, _ := json.Marshal(obj)

	// Patch the resource if already exists
	res, err := dr.Patch(ctx, obj.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{
		FieldManager: "sample-controller",
	})
	if err != nil {
		fmt.Println(err)
		// Create the resource if not found
		if strings.Contains(string(fmt.Sprint(err)), "not found") {
			fmt.Println("Creating...")
			res, err = dr.Create(ctx, obj, metav1.CreateOptions{
				FieldManager: "sample-controller",
			})
			if err != nil {
				fmt.Println("Error while creating resource: ", err)
			}
		}
	}
	return res

}
