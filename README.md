# rk
Simple tool to refactor a Kubernetes manifest set into kustomize format structured for dev, prod, release.

This is WIP.

HELLOWORLD: minimal golang app with tests that compiles
- [x] go package setup with folder structure, main, module init
- [x] makefile, gitignore, dockerignore, editorconfig
- [x] dummy refactor command using cobra including unit test
- [ ] github actions for tests / linting targeting PR's
- [ ] coveralls badge
- [ ] Precommit
- [ ] handle different log levels


MVP: take a non-kustomize manifest set and create a version using bases and overlays for dev, prod, staging
      in a namespace with an app set using crosscutting fields that renders correctly with `kubectl kustomize -k`
1. Commandline app with 2 validated args: source folder, namespace
2. Create basic folder structure
3. move existing yaml files
4. add missing kustomize files using golang templates

CLEANUP existing manifests: strip redundant namespace fields