# Dynamo Backend Challenge
[![Code Coverage](https://codecov.io/gh/nvnoskov/dynamo-backend/branch/master/graph/badge.svg)](https://codecov.io/gh/nvnoskov/dynamo-backend)

## 
This project based on Go RESTful API Starter Kit

## Install


```shell
git clone https://github.com/nvnoskov/dynamo-backend.git

cd go-rest-api

# start a PostgreSQL database server in a Docker container
make db-start

# start migrations
make migrate

# seed the database with some test data
make testdata

# run the RESTful API server
make run

# run tests
make test
```

At this time, you have a RESTful API server running at `http://127.0.0.1:8080`. It provides the following endpoints:

* `GET /healthcheck`: a healthcheck service provided for health checking purpose (needed when implementing a server cluster)
* `POST /v1/login`: authenticates a user and generates a JWT
* `POST /v1/register`: register a user
* `GET /v1/flights`: returns a paginated list of the flights 
* `GET /v1/flights/:id`: returns the detailed information of an flight
* `POST /v1/flights`: creates a new flight
* `PUT /v1/flights/:id`: updates an existing flight
* `DELETE /v1/flights/:id`: deletes an flight

Try the URL `http://localhost:8080/healthcheck` in a browser, and you should see something like `"OK v1.0.0"` displayed.


```shell
# register the user via: POST /v1/register
curl -L -X POST 'http://localhost:8080/v1/register' -H 'Content-Type: application/json' --data-raw '{
    "username": "BOEING",
    "email": "BOEING@email.com",
    "password": "123"
}'

# authenticate the user via: POST /v1/login
curl -L -X POST 'http://localhost:8080/v1/login' -H 'Content-Type: application/json' --data-raw '{
    "username": "BOEING",
    "password": "123"
}'
# should return a JWT token like: {"token":"...JWT token here..."}

# create new flight
curl -L -X POST 'http://localhost:8080/v1/flights' -H 'Authorization: Bearer ...JWT token here...' -H 'Content-Type: application/json' --data-raw '{
   "name": "BOEING 737-400",
   "number": "UR-CSV",
   "departure": "MALMÖ, SWEDEN",
   "departure_time": "2020-10-01T14:36:38Z",
   "destination": "MERZIFON, TURKEY",
   "arrival_time": "2020-10-01T17:36:38Z",
   "fare": "100EUR"
}'

# with the above JWT token, access the flight resources, such as: GET /v1/flights
curl -X GET -H "Authorization: Bearer ...JWT token here..." http://localhost:8080/v1/flights
# should return a list of flight records in the JSON format

# Search by parameters departure_time format 2020-10-01. Will search records from 2020-10-01 00:00:00 to 2020-10-01 23:59:59
curl -X GET -H "Authorization: Bearer ...JWT token here..." http://localhost:8080/v1/flights?departure_time=2020-10-01

# Search by parameters 
curl -X GET -H "Authorization: Bearer ...JWT token here..." http://localhost:8080/v1/flights?departure=MALMÖ, SWEDEN2

```



## Deployment

The application can be run as a docker container. You can use `make build-docker` to build the application 
into a docker image. The docker container starts with the `cmd/server/entryscript.sh` script which reads 
the `APP_ENV` environment variable to determine which configuration file to use. For example,
if `APP_ENV` is `qa`, the application will be started with the `config/qa.yml` configuration file.

You can also run `make build` to build an executable binary named `server`. Then start the API server using the following
command,

```shell
./server -config=./config/prod.yml
```

## TODO
 - Add partial search for the `name` parameter (like%)
 - Swagger Docs
