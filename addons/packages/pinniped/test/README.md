# Pinniped Package Unit Tests

The Go test in this directory will run unit tests across all pinniped package versions.

```text
.../packages/pinniped/a.b.c/  <- package root
  bundle/                     <- imgpkg bundle root
    ...
  fixtures/                   <- test fixtures directory
    values/                   <- input data values test cases
      test-case.yaml          <- single input data values test case
    out/                      <- expected ytt output
      test-case.yaml          <- expected ytt output for values/test-case.yaml data values
```

| To... | Command |
| ----- | ------- |
| Run all the tests | `make test` |
| Pass `go test` args | `make test GO_TEST_ARGS="-v"` |
| Run the tests for a specific package | `make test GO_TEST_FLAGS="-run TestTemplate/0.12.1"` |
| Run a specific test for multiple packages | `make test GO_TEST_FLAGS="-run TestTemplate/.*/mc-ldap-v1_5_0.yaml"` |
| Run a specific test for a specific packages | `make test GO_TEST_FLAGS="-run TestTemplate/0.12.1/mc-ldap-v1_5_0.yaml"` |
