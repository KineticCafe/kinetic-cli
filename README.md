# Kinetic CLI

This is the start of a CLI program that will be used for managing various
pieces of kinetic platform infrastructure. This is generated from
[emmanuelgautier/go-cli-template][] and partially based on [KineticCafe/kci][]
(used for infrastructure management).

<p align="left">
    <a href="https://github.com/KineticCafe/kinetic-cli/actions/workflows/ci.yml"><img src="https://github.com/KineticCafe/kinetic-cli/actions/workflows/ci.yml/badge.svg?branch=main&event=push" alt="CI Tasks for Go Cli template"></a>
    <!-- <a href="https://goreportcard.com/report/github.com/KineticCafe/kinetic-cli"><img src="https://goreportcard.com/badge/github.com/KineticCafe/kinetic-cli" alt="Go Report Card"></a> -->
    <!-- <a href="https://pkg.go.dev/github.com/KineticCafe/kinetic-cli"><img src="https://pkg.go.dev/badge/www.github.com/KineticCafe/kinetic-cli" alt="PkgGoDev"></a> -->
</p>

There are many more features to add and not enough time to add them.

## Usage

1. Clone or download this repository:

```bash
git clone https://github.com/KineticCafe/kinetic-cli.git
cd kinetic-cli
```

2. Build the CLI tool:

```bash
go build -o kinetic-cli
```

3. Run the CLI tool with the `--help` flag to see the available commands:

```bash
./kinetic-cli --help
```

## License

We have retained the licence of the Go CLI Template (MIT).

[emmanuelgautier/go-cli-template]: https://github.com/emmanuelgautier/go-cli-template
[kineticcafe/kci]: https://github.com/kineticcafe/kci
