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
- [x] verify on mac (switch to local)

Scenario: `createlayers` functionality

- [x] Goal: make UT pass for `DoMakeOverlay`
- [ ] Goal: Make command get it's stuff from cmdline to call `DoMakeOverlay`
  - [x] Add Cobra root command
  - [x] Add version command
    - [x] Fix version string
  - [x] Add layers create with flags and help
    - [x] Fix output flag
    - [x] create output dir if doesn't exist
    - [x] ignore it if inside source
    - [x] Support no output flag passed in with absolute sourcedir
    - [x] Support . for sourcedir
    - [x] source should have at least 1 k8s manifest.. if none do not create output folders
    - [ ] updateoutput with confirm
  - [x] Make command to install with version
- [ ] Goal: fix CI official build based on tags
- [ ] Goal: add suggestions to make authoring overrides easier in overlay layer
  - [ ] Add commented out reference in respective overlay folder
  - [x] Refactor `writeOverlayKustTemplate` so that path concatenation is inside function
  - [x] write patch file into overlay directory (increase replicas)
  - [x] write patch file into overlay directory (set memory for all containers)
  - [x] update integration test to test overlays - [x] Helper to call dyff library using source and dest path for single file compare - [x] Helper to recurse of tree calling above
  - [x] implement functions to get deployment name and container names from manifests
    - [x] get deployment name using kyaml
    - [x] get container names using kyaml
    - [x] write unit test for function to get deployment and container names

Scenario: prep MVP release

- [ ] Goal: get CD using same matrixed build versions. Hook up download counter to build
- [ ] Goal: BadgeApp clean
- [ ] Goal: write documentation for initial scenarios
- [ ] handle different log levels
- [ ] gate on codecoverage thresholds
- [ ] converge CI build with make script
- [ ] show build badges ()

# figure out

1. how to test multiple cobra commands / switch structures
2. how to run integration test in GA or circleci

Backlog:
Scenario: implement function to strip ?hard-coded app name/namespaces from original manifests
write unit test for namespace stripping with yaml comparison with known failure case
Scenario: `kubectl create deployement` then clean
write integration test for MakeOverlay that compares output with known good yaml that runs in prod to test functional correctness
command to add a kustomization.yaml and populate with all validd manifests. If one exists, update manifest list based on what's in dir
Add colorization to output
Output shows treestructure that was created
If there is a namespace file used to create namespace, patch that with passed in namespace
If there is no namespace file, createone
CLEANUP existing manifests: strip redundant namespace fields
Enable limits optionally via flag
