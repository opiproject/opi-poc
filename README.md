# OPI Project PoC

[![Integration](https://github.com/opiproject/opi-poc/actions/workflows/poc-integration.yml/badge.svg)](https://github.com/opiproject/opi-poc/actions/workflows/poc-integration.yml)
[![Linters](https://github.com/opiproject/opi-poc/actions/workflows/linters.yml/badge.svg)](https://github.com/opiproject/opi-poc/actions/workflows/linters.yml)
[![GitHub stars](https://img.shields.io/github/stars/opiproject/opi-poc.svg?style=flat-square&label=github%20stars)](https://github.com/opiproject/opi-poc)
[![GitHub Contributors](https://img.shields.io/github/contributors/opiproject/opi-poc.svg?style=flat-square)](https://github.com/opiproject/opi-poc/graphs/contributors)

This repo hosts OPI proofs of concept.  These PoCs are used to demonstrate that
the OPI project can work for some set of use cases.  As a result, OPI CI is
built on top of these PoCs.

## High level requirements

A PoC should demonstrate some aspect or aspects of OPI, for example a networking
PoC could demonstrate some networking application using OPI APIs, provisioning,
and lifecycle management.

In most cases, the PoC should

* be implemented by some set of containers
* include tests
* run on a variety of hardware configurations including virtual and multiple
  vendors' devices

## Current PoCs

* [Developer Platform / Integration Platform](integration/README.md) aka The CI
* [Storage](storage/README.md)
* [Security](security/README.md)
* [Demos](demos/README.md)

## LAB

* [Lab hardware planning](lab/goals-and-requirements.md)
* [LAB current state - Phase 1](lab/README.md)

## I Want To Contribute

This project welcomes contributions and suggestions.  We are happy to have the
Community involved via submission of **Issues and Pull Requests** (with
substantive content  or even just fixes). We are hoping for the documents,
test framework, etc. to become a community process with active engagement.
PRs can be reviewed by any number of people, and a maintainer may accept.

See [CONTRIBUTING](https://github.com/opiproject/opi/blob/main/CONTRIBUTING.md)
and [GitHub Basic Process](https://github.com/opiproject/opi/blob/main/doc-github-rules.md)
for more details.
