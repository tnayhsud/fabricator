package resource

const DeploymentYAML string = `
apiVersion: apps/v1
kind: Deployment
metadata:
    name: nginx-deployment
    namespace: test
    labels:
        app: svc-1
        test-fabric: "true"
        sid: "1"
spec:
    replicas: 1
    selector:
        matchLabels:
            app: svc-1
    template:
        metadata:
            labels:
                app: svc-1
                test-fabric: "true"
                sid: "1"
        spec:
            containers:
                - name: container
                  image: dvaishnav/ping:0.0.1
                  env:
                      -   name: SERVICE_LIST
                          valueFrom:
                              configMapKeyRef:
                                  name: config-map
                                  key: services
                  ports:
                      - containerPort: 8081
`

type EnvVar struct {
	Name      string `yaml:"name"`
	ValueFrom struct {
		ConfigMapKeyRef struct {
			Name string `yaml:"name"`
			Key  string `yaml:"key"`
		} `yaml:"configMapKeyRef"`
	} `yaml:"valueFrom"`
}

type containerPort map[string]int

type Container struct {
	Name  string          `yaml:"name"`
	Image string          `yaml:"image"`
	Env   []EnvVar        `yaml:"env"`
	Ports []containerPort `yaml:"ports"`
}

type Deployment struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
		Labels    struct {
			App        string
			TestFabric string `yaml:"test-fabric"`
			Sid        string `yaml:"sid"`
		} `yaml:"labels"`
	} `yaml:"metadata"`
	Spec struct {
		Replicas int `yaml:"replicas"`
		Selector struct {
			MatchLabels struct {
				App string
			} `yaml:"matchLabels"`
		} `yaml:"selector"`
		Template struct {
			Metadata struct {
				Labels struct {
					App        string
					TestFabric string `yaml:"test-fabric"`
					Sid        string
				} `yaml:"labels"`
			} `yaml:"metadata"`
			Spec struct {
				Containers []Container `yaml:"containers"`
			} `yaml:"spec"`
		} `yaml:"template"`
	} `yaml:"spec"`
}
