# OPI Project PoC

[![OPI Integration CI](https://github.com/opiproject/opi-poc/actions/workflows/poc-integration.yml/badge.svg)](https://github.com/opiproject/opi-poc/actions/workflows/poc-integration.yml)
[![OPI Networking PoC CI](https://github.com/opiproject/opi-poc/actions/workflows/poc-networking.yml/badge.svg)](https://github.com/opiproject/opi-poc/actions/workflows/poc-networking.yml)
[![OPI Storage PoC CI](https://github.com/opiproject/opi-poc/actions/workflows/poc-storage.yml/badge.svg)](https://github.com/opiproject/opi-poc/actions/workflows/poc-storage.yml)
[![OPI Security PoC CI](https://github.com/opiproject/opi-poc/actions/workflows/poc-security.yml/badge.svg)](https://github.com/opiproject/opi-poc/actions/workflows/poc-security.yml)
[![Linters](https://github.com/opiproject/opi-poc/actions/workflows/linters.yml/badge.svg)](https://github.com/opiproject/opi-poc/actions/workflows/linters.yml)

This repo hosts OPI proofs of concept.  These PoCs are used to demonstrate that
the OPI project can work for some set of use cases.

## High level requirements

A PoC should demonstrate some aspect or aspects of OPI, for example a networking
PoC could demonstrate some networking application using OPI APIs, provisioning,
and lifecycle management.

In most cases, the PoC should

* be implemented by some set of containers
* include tests
* run on a variety of hardware configurations including virtual and multiple
  vendors' devices

## I Want To Contribute

This project welcomes contributions and suggestions.  We are happy to have the
Community involved via submission of **Issues and Pull Requests** (with
substantive content  or even just fixes). We are hoping for the documents,
test framework, etc. to become a community process with active engagement.
PRs can be reviewed by any number of people, and a maintainer may accept.

See [CONTRIBUTING](https://github.com/opiproject/opi/blob/main/CONTRIBUTING.md)
and [GitHub Basic Process](https://github.com/opiproject/opi/blob/main/doc-github-rules.md)
for more details.

## Current PoCs

* [Developer Platform / Integration Platform](integration/README.md)
* [Pure software networking](networking/README.md)
* [Storage](storage/README.md)
