# redis
base on go-redis library

<a href="https://github.com/goapt/redis/actions"><img src="https://github.com/goapt/redis/workflows/build/badge.svg" alt="Build Status"></a>
<a href="https://codecov.io/gh/goapt/redis"><img src="https://codecov.io/gh/goapt/redis/branch/master/graph/badge.svg" alt="codecov"></a>
<a href="https://goreportcard.com/report/github.com/goapt/redis"><img src="https://goreportcard.com/badge/github.com/goapt/redis" alt="Go Report Card
"></a>
<a href="https://pkg.go.dev/github.com/goapt/redis"><img src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square" alt="GoDoc"></a>
<a href="https://opensource.org/licenses/mit-license.php" rel="nofollow"><img src="https://badges.frapsoft.com/os/mit/mit.svg?v=103"></a>



## Usage

```shell script
go get github.com/goapt/redis
```

```go
configs := make(map[string]redis.Config)

configs["default"] = Config{
    Server: "127.0.0.1:6379",
}

redis.Connect(configs)
client := redis.NewRedisWithName("default")
```