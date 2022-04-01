# Test Fabricator

Test Fabricator is a GO package to create K8s resources for generating network traffic for testing/monitoring purpose.

- Built using [client-go] library, [apimachinery]
- Uses a docker image to enable communication between different pods

# Usage of the package

As of now, the package can be used to generate network traffic in K8s cluster and then the data can be pushed to Prometheus, Grafana or any other monitoring tool like skywalking, to quickly generate meaningful data to populate graphs, tables or charts and showcase in any POC.

# Example 

```sh
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fabricator/fabricate" 

	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	userHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home directory: %v\n", err)
		os.Exit(1)
	}
    // Filepath to K8s config 
	kubeConfig := filepath.Join(userHome, ".kube", "config")
    masterURL := ""
	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeConfig)
	if err != nil {
		fmt.Printf("Error building config: %v\n", err)
		os.Exit(1)
	}
    // Requires an integer to specify number of resources i.e. number of config maps, deployments and services
    count := 5
    // Requires an integer to specify number of pods to create per deployment
    replicas := 3
    // Required kube config to authenticate with the kube-apiserver
	fabricator := fabricate.New(config, count, replicas)
	fabricator.Fabricate()
}
```

[//]: # (These are reference links used in the body of this note and get stripped out when the markdown processor does its job. There is no need to format nicely because it shouldn't be seen.)

   [client-go]: <https://github.com/kubernetes/client-go>
   [apimachinery]: <https://github.com/kubernetes/apimachinery>