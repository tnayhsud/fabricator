package resource

var ConfigYAML string = `
apiVersion: v1
kind: ConfigMap
metadata:
    name: config-map-1
    namespace: test
    labels:
        test-fabric: "true"
        sid: "1"
data:
    services: "svc-2,svc-3,svc-4,svc-5"
`

type ConfigMap struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
		Labels    struct {
			TestFabric string `yaml:"test-fabric"`
			Sid        string `yaml:"sid"`
		} `yaml:"labels"`
	} `yaml:"metadata"`
	Data struct {
		Services string `yaml:"services"`
	} `yaml:"data"`
}
