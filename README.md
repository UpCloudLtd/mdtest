# mdtest

[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/UpCloudLtd/mdtest/badge)](https://scorecard.dev/viewer/?uri=github.com%2FUpCloudLtd%2Fmdtest)

Tool for combining examples and test cases. Parses markdown files for test steps defined as code blocks and uses these to test command line applications.

## Usage

Build binary and test examples:

```sh
make
./bin/mdtest examples/
```

## Development

Use [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) when committing your changes.

To lint the code, run `golangci-lint run`. See its documentation for  [local installation instructions](https://golangci-lint.run/usage/install/#local-installation).

```sh
golangci-lint run
```

To test the code, run `go test ./...`.

```sh
go test ./...
```

To build the application and execute the tests in [examples/](./examples/) directory, run:

```sh
make
./bin/mdtest examples/
```
