# gomicroservice-search

The search service for my gomicroservice project

Prerequisites:

- make
- Go 1.13.x
- docker

## Run Locally(Docker)

```sh
make docker
```

## Run Tests

The pipeline is robust containing:

- static code analysis
- safesql check
- unit tests
- integration tests using cucumber
- benchmarking
