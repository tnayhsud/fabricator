package client

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

type IClientUtil interface {
	DynamicResourceInterface(cfg *rest.Config, mapping *meta.RESTMapping, obj *unstructured.Unstructured) dynamic.ResourceInterface
	ResourceMapping(cfg *rest.Config, gvk *schema.GroupVersionKind) *meta.RESTMapping
}

type ClientUtil struct {
}

func InitClientUtil() IClientUtil {
	return &ClientUtil{}
}

/*
	func 	dynamicResourceInterface
	desc 	creates resource interface, for the given mapping, that provides method such as Create, Update, Patch etc
	param	cfg *rest.Config
	param	mapping *meta.RESTMapping
	param	obj *unstructured.Unstructured
	returns	dynamic.ResourceInterface
*/
func (c ClientUtil) DynamicResourceInterface(cfg *rest.Config, mapping *meta.RESTMapping,
	obj *unstructured.Unstructured) dynamic.ResourceInterface {
	dyn, err := dynamic.NewForConfig(cfg)
	if err != nil {
		fmt.Println(err)
	}
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		/* namespaced resources should specify the namespace */
		return dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else { /* for cluster-wide resources */
		return dyn.Resource(mapping.Resource)
	}
}

/*
	func 	resourceMapping
	desc 	queries the API server to find preferred resource mapping for the provided group kind
	param	cfg *rest.Config
	param	gvk *schema.GroupVersionKind
	returns	*meta.RESTMapping
*/
func (c ClientUtil) ResourceMapping(cfg *rest.Config, gvk *schema.GroupVersionKind) *meta.RESTMapping {
	// creates a new DiscoveryClient to discover supported resources in the API server
	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		fmt.Println(err)
	}
	// returns a DeferredDiscoveryRESTMapper that will lazily query the provided client for discovery information
	// to do REST mappings
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))
	// call to RESTMapping identifies a preferred resource mapping for the provided group kind
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)

	if err != nil {
		fmt.Println(err)
	}
	return mapping
}

type CustomRestScope struct {
	name meta.RESTScopeName
}

func (rs *CustomRestScope) Name() meta.RESTScopeName {
	return rs.name
}
