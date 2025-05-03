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
go get github.com/gin-gonic/gin
go get github.com/mattn/go-sqlite3
go get github.com/jung-kurt/gofpdf
go get github.com/gin-gonic/gin
go get github.com/mattn/go-sqlite3
go get gopkg.in/natefinch/lumberjack.v2
go get github.com/fsnotify/fsnotify
go get github.com/spf13/viper
<!-- go get github.com/jinzhu/gorm -->
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