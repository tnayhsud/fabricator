package resource

var ServiceYAML string = `
apiVersion: v1
kind: Service
metadata:
    name: service
    namespace: test
    labels:
        test-fabric: "true"
        sid: "1"
spec:
    selector:
        app: svc-1
        sid: "1"
    ports:
        - protocol: TCP
          port: 8080
          targetPort: 8081
`

type Port struct {
	Protocol   string `yaml:"protocol"`
	Port       int    `yaml:"port"`
	TargetPort int    `yaml:"targetPort"`
}

type Service struct {
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
	Spec struct {
		Selector struct {
			App string
			Sid string
		} `yaml:"selector"`
		Ports []Port `yaml:"ports"`
	} `yaml:"spec"`
}
