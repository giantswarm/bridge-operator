module github.com/giantswarm/bridge-operator

go 1.14

require (
	github.com/giantswarm/apiextensions v0.2.0
	github.com/giantswarm/exporterkit v0.2.0
	github.com/giantswarm/k8sclient v0.2.0
	github.com/giantswarm/microendpoint v0.2.0
	github.com/giantswarm/microerror v0.2.1
	github.com/giantswarm/microkit v0.2.0
	github.com/giantswarm/micrologger v0.3.1
	github.com/giantswarm/operatorkit v0.2.0
	github.com/giantswarm/versionbundle v0.2.0
	github.com/google/go-cmp v0.4.0
	github.com/prometheus/client_golang v1.3.0
	github.com/spf13/viper v1.6.2
	k8s.io/apimachinery v0.16.6
	k8s.io/client-go v0.16.6
)

replace github.com/gorilla/websocket v1.4.0 => github.com/gorilla/websocket v1.4.2

// v3.3.X is required by sigs.k8s.io/controller-runtime. Can remove this replace when updated.
replace github.com/coreos/etcd => github.com/coreos/etcd v3.3.25+incompatible
