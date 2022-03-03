# rk
Simple tool to refactor a Kubernetes manifest set into kustomize format.

This is WIP.

HELLOWORLD: minimal golang app with tests that compiles

CICD:
Linting, test running, mac / linux binary publishing as releases

MVP: take a non-kustomize manifest set and create a version using bases and overlays for dev, prod, staging
      in a namespace with an app set using crosscutting fields that renders correctly with `kubectl kustomize -k`
1. Commandline app with 2 validated args: source folder, namespace
2. Create basic folder structure
3. move existing yaml files
4. add missing kustomize files using golang templates

CLEANUP existing manifests: strip redundant namespace fields