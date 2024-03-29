module github.com/giantswarm/bridge-operator

go 1.14

require (
	github.com/giantswarm/apiextensions/v3 v3.39.0
	github.com/giantswarm/exporterkit v0.2.1
	github.com/giantswarm/k8sclient/v5 v5.12.0
	github.com/giantswarm/microendpoint v0.2.0
	github.com/giantswarm/microerror v0.3.0
	github.com/giantswarm/microkit v0.2.2
	github.com/giantswarm/micrologger v0.5.0
	github.com/giantswarm/operatorkit v1.2.0
	github.com/giantswarm/versionbundle v0.2.0
	github.com/google/go-cmp v0.5.9
	github.com/prometheus/client_golang v1.17.0
	github.com/spf13/viper v1.9.0
	k8s.io/apimachinery v0.18.19
	k8s.io/client-go v0.18.19
	sigs.k8s.io/cluster-api v0.3.16 // indirect
)

replace (
	// v3.3.X is required by sigs.k8s.io/controller-runtime. Can remove this replace when updated.
	github.com/coreos/etcd => github.com/coreos/etcd v3.3.25+incompatible
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
	golang.org/x/net => golang.org/x/net v0.19.0
	golang.org/x/text => golang.org/x/text v0.13.0
	sigs.k8s.io/cluster-api => github.com/giantswarm/cluster-api v0.3.13-gs
)
