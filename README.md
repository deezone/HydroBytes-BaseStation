# HydroBytes-BaseStation
The Base Station is a part of a collection of systems called
"[HydroBytes](https://github.com/deezone/HydroBytes)" that uses micro
controllers to manage and monitor plant health in an automated garden.

## Introduction

The "garden" is simply a backyard patio in Brooklyn, New York. Typically
there are only two seasons - cold and hot in New York City. By
automating an urban garden ideally the space will thrive with minimum
supervision. The amount of effort to automate is besides the point, everyone needs their vices.

- **[Water Station](https://github.com/deezone/HydroBytes-WaterStation)**
- **Base Station**
- **[Plant Station](https://github.com/deezone/HydroBytes-PlantStation)**

![brooklyn-20201115 garden layout](https://raw.githubusercontent.com/deezone/HydroBytes/master/resources/gardenBrooklynDiagram-20201115.jpg)

![Garden](https://github.com/deezone/HydroBytes-WaterManagement/blob/master/resources/garden-01.png)

### YouTube Channel

[![YouTube Channel](https://github.com/deezone/HydroBytes-WaterStation/blob/master/resources/youTube-TN.png?raw=true)](https://www.youtube.com/channel/UC00A_lEJD2Hcy9bw6UuoUBA "All of the HydroBytes videos")

### Notes

Development of a Go based API is based on instruction in the amazing
courses at **[Ardan Labs](https://education.ardanlabs.com/collections?category=courses)**.

#### Starting Web Server
```
> go run ./cmd/api
STATIONS API : 2021/01/30 23:34:33.625072 main.go:103: main : Started
STATIONS API : 2021/01/30 23:34:33.625255 main.go:110: main : Config :
--web-address=localhost:8000
--web-debug=localhost:6060
--web-read-timeout=5s
--web-write-timeout=5s
--web-shutdown-timeout=5s
--db-user=postgres
--db-host=localhost
--db-name=postgres
--db-disable-tls=true
--auth-key-id=1
--auth-private-key-file=private.pem
--auth-algorithm=RS256
--trace-url=http://localhost:9411/api/v2/spans
--trace-service=station-api
--trace-probability=1
STATIONS API : 2021/01/30 23:34:33.628227 main.go:198: main : API listening on localhost:8000
STATIONS API : 2021/01/30 23:34:33.628284 main.go:163: debug service listening on localhost:6060

^C
STATIONS API : 2021/01/30 23:35:18.575270 main.go:222: main : Start shutdown
STATIONS API : 2021/01/30 23:35:18.575380 main.go:245: main : Completed
```

- supported requests to `localhost:8000`:
  - `GET /v1/account/token`
  - `GET  /v1/station-types`
  - `GET  /v1/station-type/{id}`
  - `POST /v1/station-type`
  - `DELETE /v1/station-type/{id}`
  - `GET  /v1/station-type/{station-type-id}/stations`
  - `POST /v1/station-type/{station-type-id}/station`
  - `DELETE /v1/station/{id}`
  - `GET /v1/health`

- Debugging requests to `http://localhost:6060/debug/pprof/`

#### Admin tools

```
> go run ./cmd/admin -h migrate
Usage: admin [options] [arguments]

OPTIONS
  --web-address/$STATIONS_WEB_ADDRESS                      <string>    (default: localhost:8000)
  --web-debug/$STATIONS_WEB_DEBUG                          <string>    (default: localhost:6060)
  --web-read-timeout/$STATIONS_WEB_READ_TIMEOUT            <duration>  (default: 5s)
  --web-write-timeout/$STATIONS_WEB_WRITE_TIMEOUT          <duration>  (default: 5s)
  --web-shutdown-timeout/$STATIONS_WEB_SHUTDOWN_TIMEOUT    <duration>  (default: 5s)
  --db-user/$STATIONS_DB_USER                              <string>    (default: postgres)
  --db-password/$STATIONS_DB_PASSWORD                      <string>    (noprint,default: postgres)
  --db-host/$STATIONS_DB_HOST                              <string>    (default: localhost)
  --db-name/$STATIONS_DB_NAME                              <string>    (default: postgres)
  --db-disable-tls/$STATIONS_DB_DISABLE_TLS                <bool>      (default: false)
  --auth-key-id/$STATIONS_AUTH_KEY_ID                      <string>    (default: 1)
  --auth-private-key-file/$STATIONS_AUTH_PRIVATE_KEY_FILE  <string>    (default: private.pem)
  --auth-algorithm/$STATIONS_AUTH_ALGORITHM                <string>    (default: RS256)
  --trace-url/$STATIONS_TRACE_URL                          <string>    (default: http://localhost:9411/api/v2/spans)
  --trace-service/$STATIONS_TRACE_SERVICE                  <string>    (default: station-api)
  --trace-probability/$STATIONS_TRACE_PROBABILITY          <float>     (default: 1)
  --help/-h
  display this help message
```
- `adminAdd` to create a new admin account from CLI
```
> go run ./cmd/admin adminAdd Test2Admin supersecret
Admin account will be created with password "supersecret"
Continue? (1/0) 1
Account created with id: bd86e286-b1fc-415b-8f80-c30dcdac10bb
```

- `keygen`
```
> go run ./cmd/admin keygen private.pem
> cat private.pem
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAsdRBS6Cspo6uXKUetnEMifL7xM8pz0pSz7gDhEb/OH6eT8oO
yXpZlMvCOpXvnHCBZwa2C1zmYCW1v0nEFxyGYCFsMbmU09+7cz4qah5q9n6bTixB
...
Zz1gJz/yNc4n38M7SILauHLWmNBhoisb2axc1VX4X8D/oAwb+vDFlhR2QIAqTQ7N
Ikyxx1arCSpoDd1M48SS1+xkgaDqEZuEA+COyKUyZf41PtXL6MA73Q==
-----END RSA PRIVATE KEY-----
```

- `migrate` to update database with schema defined in code.
```
> go run ./cmd/admin migrate
Migrations complete
```

- `seed` populate the database tables with seed data for testing and development.
```
> go run ./cmd/admin seed
Seeding complete
```

#### Tests

- **Unit Tests**

**station_type** and **station**
```
> go test ./internal/station_type
ok  	github.com/deezone/HydroBytes-BaseStation/internal/station_type	13.515s
```

NOTE: test coverage reports:
```
alias gotwc='go test -coverprofile=coverage.out && go tool cover -html=coverage.out && rm coverage.out'
```

- **Functional tests**
```
# bust test cache
> go clean -testcache

> go test ./cmd/api/tests/station_tests
ok  	github.com/deezone/HydroBytes-BaseStation/cmd/api/tests/station_tests	3.248s

> go test ./cmd/api/tests/station_type_tests
ok  	github.com/deezone/HydroBytes-BaseStation/cmd/api/tests/station_type_tests	2.875s

>  go test ./cmd/api/tests/account_tests
ok  	github.com/deezone/HydroBytes-BaseStation/cmd/api/tests/account_tests	4.372s
```

#### Tracing
- http://localhost:9411/zipkin/
