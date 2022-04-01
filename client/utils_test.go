package client

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/fabricator/fabricate/client"

// 	"github.com/fabricator/fabricate/resource"

// 	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
// 	"k8s.io/apimachinery/pkg/runtime"
// 	yamlV2 "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
// 	dynamicfakeclient "k8s.io/client-go/dynamic/fake"
// )

// func TestDynamicResourceInterface(t *testing.T) {
// 	config := config()
// 	mapping, gvk := mapping()
// 	serializer := yamlV2.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
// 	obj := &unstructured.Unstructured{}
// 	serializer.Decode([]byte(resource.DeploymentYAML), nil, obj)
// 	dyn := dynamicfakeclient.NewSimpleDynamicClient(runtime.NewScheme())

// 	clientUtil := client.InitClientUtil()
// 	dr := clientUtil.DynamicResourceInterface(config, mapping, obj)
// 	fmt.Println(dr)
// 	// assert.Equal(t, 123, 123, "they should be equal")
// }
