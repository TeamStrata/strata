# Strata

## Docker Setup
The root of the project contains a Dockerfile and compose.yml file. 
This allows a quick installation using the following command within this project's root source code directory.
<br>
To run the following commands, you must have [Docker-Compose](https://docs.docker.com/compose/install/) installed. Alternative, you can use the [Go Setup](#go-setup) procedure. 

```sh
docker compose up
```
The container can be shut down and removed, while keeping the built image, using the following command:

```sh
docker compose down
```

## Ollama Setup
After starting the docker containers, the strata\_ollama contaner is uninitialized. The ollama container hosts an HTTP server on port 11434, which allows access to the AI model over a JSON REST API. A `Modelfile` has been provided in the root of this project, which can be installed using the command:
```sh
strata exec -it strata_ollama ollama create stratabot -f /root/.ollama/Modelfile
```

### Executing commands in a Running Docker Container
A shell within the Strata and Postgres containers can be created using the following respective commands:
```sh
# Strata
docker exec -it strata bash

# Postgres
docker exec -it strata_psql bash
```
To run any command within the running docker instance, replace `bash` with the command to execute. Refer to the [Running Tests](#running-tests) section where this is utilized to run the go command within the docker instance. 
<br>
To access the 'psql' command within postgresql, ensure the username and database are set to 'strata' (by passing parameters `-U` and `-d` respectively) such as the following command will perform:

```sh
docker exec -it strata_psql psql -U strata -d strata -f /*.d/1_schema.sql
```

## Go Setup
The following set of commands can be used to install required packages, then build the 'strata' executable in the project root directory. These commands require your terminal's current working directory is the project's root directory. Once the executable is built, the docker container can be restarted. It will automatically use the new `strata` executable in favour of the fallback stored in the docker container filesystem, at `/bin/strata`. 
<br>
These commands can be prefixed by `docker exec -it strata` if you do not have go installed. 
```sh
# Install packages:
go install ./cmd

# Build project
go build -o strata ./cmd

# Result binary is placed in ./strata
``` 

## Running Tests
To run tests within the docker container, you can use the command below:

```sh
docker exec -it strata go test -v ./pkg/tests
```
