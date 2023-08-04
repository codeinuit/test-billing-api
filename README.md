# Billing API 
[![Go Report Card](https://goreportcard.com/badge/github.com/codeinuit/test-billing-api)](https://goreportcard.com/report/github.com/codeinuit/test-billing-api)

A billing API with user management wrote in Go for a test.

## Setup
### Configuration
```
# Postgres configuration
POSTGRES_HOST=
POSTGRES_PORT=
POSTGRES_USER=
POSTGRES_PASS=

# port, default to 8080
PORT=
```

### Running with Docker
The project comes with a Docker Compose v3 that runs PostgreSQL database and the API. You can update the configuration in the `docker-compose.yml` file.

```
$ git clone git@github.com:codeinuit/test-billing-api.git
$ cd test-billing-api
$ docker compose up -d --build
```

### Running on Linux
On Linux, to clone and build this application, you’ll need Git and [Golang](https://go.dev/doc/install) installed. You’ll also need a PostgreSQL database running aside and export the configuration as environnement variable as specified above.

```
$ git clone git@github.com:codeinuit/test-billing-api.git
$ cd test-billing-api

# clean the project, install dependancies and build the binary
$ make

# run the application
$ ./bin/billing-api
```

## About the development and possible enhancement
This project has been made for a technical test in less than a day.
- Docker, log implementation : 2h
- Endpoints, connectors : 4h
- README : 45min

There is also some points that could be improved in the future
- **The database schema could have a dedicated table for transactions**. That could allow users to pay up an invoice in multiple transactions, but also we could keep traces on this.
- **Database could be an implentation**. Instead of relying directly on PostgreSQL we could use an implentation for supporting updates, or for using a different database / library
- **Handling Docker stop signals**
- **Use a configuration library** to support `.env` files and/or CLI arguments
- **Add versionning on routes** to handle breaking changes
- **Query optimisation**
- **Use `FLOAT` instead of `BIGINT`** on amounts or balance variables on the database schema


## Library used
- [gin-gonic/gin](https://github.com/gin-gonic/gin)
- [sirupsen/logrus](https://github.com/sirupsen/logrus)
- [lib/pq](github.com/lib/pq)
