package project

var (
	description string = "The bridge-operator handles network bridges in KVM hardware machines."
	gitSHA      string = "n/a"
	name        string = "bridge-operator"
	source      string = "https://github.com/giantswarm/bridge-operator"
	version            = "n/a"
)

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
}

func Name() string {
	return name
}

func Source() string {
	return source
}

func Version() string {
	return version
}
