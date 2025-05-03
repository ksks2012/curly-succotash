# structure

- etc: setting file
- docs: document
- global: global variables
- internal (internal module):
 <!-- TODO: -->
- dao: data access object
- middleware
- model: database model control
- routers: api routes
- service: process business logic
- pkg: package
- storage: temp file
- scripts: build, install, analysis scripts
- third_party: third_party tools

# Golang Package

```
go mod tidy
```

# Go generate

- TODO:

# Build

```sh
go build curly-succotash/backend/cmd/${PROJECT_NAME}
```

# Test

```
go test ./testing/...
```