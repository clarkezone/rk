# rk

[![Coverage Status](https://coveralls.io/repos/github/clarkezone/rk/badge.svg?branch=main)](https://coveralls.io/github/clarkezone/rk?branch=main)

Simple tool to refactor a Kubernetes manifest set
into kustomize format structured for dev, prod, release.

This is WIP. Very WIP!

HELLOWORLD: minimal golang app with tests that compiles

- [x] go package setup with folder structure, main, module init
- [x] makefile, gitignore, dockerignore, editorconfig
- [x] dummy refactor command using cobra including unit test
- [x] github actions for tests / linting targeting PR's
- [x] enable coverage in VSCode UI
- [x] go linting
- [x] all other linting
- [x] coveralls badge
- [x] Precommit
- [ ] handle different log levels
- [ ] gate on codecoverage thresholds
- [ ] converge CI build with make script
- [ ] show build badges ()

Scenario: take a non-kustomize manifest set and create a version using bases and
overlays for dev, prod, staging in a namespace with an app set using
crosscutting fields that renders correctly with `kubectl kustomize -k`

1. `rk -makeoverlay dev, staging, prod -source .`
2. Create basic folder structure
3. move existing yaml files
4. add missing kustomize files using golang templates

Goal: Inner loop

- [x] setup .exe to point at actual thing we're trying to solve hard coded
- [x] setup input test dir
- [x] setup build command in make file
- [x] setup UT infra: empty test, prep/clean, dirs
- [x] build dir creation logic
- [x] integration test that calls `kubectl kustomize ./`
- [x] build manifest templated add logic
- [x] verify against jekyll project
- [ ] verify on mac (switch to local)

- [x] Goal: make UT pass for `DoMakeOverlay`
- [ ] Goal: add suggestions to make authoring overrides easier in overlay layer
- [ ] Refactor `writeOverlayKustTemplate` so that path concatenation is inside function
- [ ] write two patch files into overlay directory (set memory, increase replicas)
- [ ] Add commented out reference in respective overlay folder

- [ ] Goal: get CD using same matrixed build versions
- [ ] Goal: Make command get it's stuff from cmdline to call `DoMakeOverlay`

# figure out

1. how to get dir where cmd is running
2. how to combine multiple cobra commands / switch structures
3. how to test multiple cobra commands / switch structures
4. how to return args in cobraCommand
5. how to incorporate release into github action
6. how to run integration test in GA or circleci

Backlog:
If there is a namespace file used to create namespace, patch that with passed in namespace
If there is no namespace file, createone
CLEANUP existing manifests: strip redundant namespace fields
