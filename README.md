[![CircleCI](https://dl.circleci.com/status-badge/img/gh/giantswarm/bridge-operator/tree/master.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/giantswarm/bridge-operator/tree/master)

# bridge-operator

`bridge-operator` is deployed as a DaemonSet on all KVM installations. It emits metrics for orphaned flannel resources.

## Getting the Project

Download the latest release:
https://github.com/giantswarm/bridge-operator/releases/latest

Clone the git repository: https://github.com/giantswarm/bridge-operator.git

Download the latest docker image from here:
https://quay.io/repository/giantswarm/bridge-operator


### How to build

```
go build github.com/giantswarm/bridge-operator
```

## Architecture

The operator is pretty simple, just spawn a server to check the health of the workload. On the other hand there is a collector that lists all clusters and check if there is any leftover from a deleted cluster.

## Contact

- Mailing list: [giantswarm](https://groups.google.com/forum/!forum/giantswarm)
- Bugs: [issues](https://github.com/giantswarm/bridge-operator/issues)

## Contributing & Reporting Bugs

See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches, the
contribution workflow as well as reporting bugs.

For security issues, please see [the security policy](SECURITY.md).


## License

bridge-operator is under the Apache 2.0 license. See the [LICENSE](LICENSE) file
for details.


## Credit
- https://golang.org
