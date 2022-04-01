package resource

const NamespaceYAML string = `
apiVersion: v1
kind: Namespace
metadata:
    name: test
    labels:
        sid: "1"
`

type Namespace struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name   string `yaml:"name"`
		Labels struct {
			SId string `yaml:"sid"`
		} `yaml:"labels"`
	} `yaml:"metadata"`
}
