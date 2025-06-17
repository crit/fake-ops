# Fake Ops

- What? Fake responses from remote services.
- Why? I wanted it.
- How? Go, Bubbletea, fsnotify, Gin

__NOTICE:__ This is a development aid and not intended to be used to run production workloads. 

## Usage

```shell
cd examples/ && fake-ops
```

![fake-ops screenshot](https://i.critrussell.net/2025/0TzwkbemCqWTYre.png)

Flags:

- `--services` Directory of service files. See [examples/services](examples/services).
  - default: `./services`
- `--results` Directory of http result files. See [examples/results](examples/results).
  - default: `./results`

## Install

```shell
go install .
```

_OR_

```shell
go install github.com/crit/fake-ops
```

## Creating Services

### HTTP Service File

```yaml
name: products   # Must match a folder containing response files.
type: http       # Indicates this is an HTTP server service.
port: 3002       # Port to run the HTTP server on.
skip: true       # If true, skips running the service but lists it.
```

### App Service File

Example uses the temporal cli published by [Temporal.io](https://docs.temporal.io/cli)

```yaml
name: temporal   # Name of the service.
type: app        # Indicates this is an Application service.
port: 7233       # Port that the application will run on for communication.
stdout: false    # If true, stdout will be piped to fake-ops.
stderr: false    # If true, stderr will be piped to fake-ops.
skip: false      # If true, skips running the service but lists it.
exec: temporal server start-dev --port {port} # Command to run this application. {port} will be substituted at runtime.
```

## Creating HTTP Response Files

See [examples/results/users](examples/results/users)

```yaml
# GET /users/:id 200 application/json
{
  "status": "SUCCESS",
  "message": "User found.",
  "data": {
    "user": {
      "id": "1",
      "name": "Alice Johnson",
      "email": "alice.johnson@example.com",
      "role": "admin"
    }
  }
}
```

### Anatomy of a Response File

First line contains the data needed to serve this file's contents. Space delimited.

1. Must start with `#`
2. HTTP Method
3. Route including any path parameters.
4. HTTP Status Code
5. Content-Type of the response.

All content after the first line is used as the response body.

__NOTE:__ yaml in this case is used for syntax highlighting of JSON responses. You can choose any file
format that suites your needs. See [examples/results/static](examples/results/static) for more variety.

### Hot Reloading

`fake-ops` monitors all response files and their parent directory for changes. It will reload the HTTP service when 
one of the files is changed or a new one created.
