# SAS Golang Generics Development Kit

## Overview

A library for developing in go with generics. See the doc.go file in each
package for information about that package.

### Packages

* [collections](./pkg/collections/doc.go)
* [condition](./pkg/condition/doc.go)
* [di](./pkg/di/doc.go)
* [embedres](./pkg/embedres/doc.go)
* [errors](./pkg/errors/doc.go)
* [filters](./pkg/filters/doc.go)
* [folders](./pkg/folders/doc.go)
* [jsonutils](./pkg/jsonutils/doc.go)
* [maputils](./pkg/maputils/doc.go)
* [maybe](./pkg/maybe/doc.go)
* [pointer](./pkg/pointer/doc.go)
* [processutils](./pkg/processutils/doc.go)
* [result](./pkg/result/doc.go)
* [sliceutils](./pkg/sliceutils/doc.go)
* [stack](./pkg/stack/doc.go)
* [streamutils](./pkg/streamutils/doc.go)
* [stringutils](./pkg/stringutils/doc.go)
* [timeutils](./pkg/timeutils/doc.go)

### Installation

    go get github.com/sassoftware/sas-ggdk

## Contributing

We welcome your contributions! Please read [CONTRIBUTING.md](CONTRIBUTING.md)
for details on how to submit contributions to this project.

### Clone

```bash
git clone https://github.com/sassoftware/sas-ggdk.git
cd sas-ggdk
```

### Install tools

This project uses [bingo](https://github.com/bwplotka/bingo) for managing local
tools. See the bingo project for instructions on installing bingo. Install
supporting tools used by this project with

```bash
make tools
```

### Run tests

Contributions must ensure linting and all tests continue to pass and any new
code is covered.

```bash
make lint test
```

A coverage report can be produced in `./build/reports/coverage.html` with the
following commands.

```bash
make test-with-coverage
```

### Commits

This project currently uses loose [semantic versioning](https://semver.org/) and
[conventional commits](https://www.conventionalcommits.org/en/v1.0.0/). The API
is undergoing rapid evolution until 1.0 is released. Efforts will be made to
communicate braking changes but please pin to a specific version until then.

## License

This project is licensed under the [Apache 2.0 License](LICENSE).

