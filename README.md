# Luke Stratta - Flatpeak Take Home Task

Thank you for taking the time to look over my implementation of this task.

If you want to try out the application, head to [Getting Started](#getting-started).

## Contents

### Getting Started

1. [About the project](#about-the-project)
2. [Project structure and extras](#project-structure-and-extras)
3. [Getting started](#getting-started)
4. [Concsessions & assumptions](#concessions--assumptions-made)

### API Documentation

#### Slots API endpoints

1. [GET /slots](#get-slots)

## About the Project

This application fetches the time slots with the lowest carbon intensity over a specified period. 

It will either:

1. Provide the time slots as individual time periods that are the lowest out of the time period provided. Time slots returned aren't necessarily adjacent to each other in order.
2. Provide one time slot that shows the average over the specified time period, with weighting based on how long you have specified. The time slot returned is from the time you initiated the request.

### Project Structure and Extras

The `main.go` file lives in the `cmd/` directory. This is the main application directory and the entrypoint for the application.

Nearly everything else lives in the `internal/` directory.

There is a `docker/` directory that hosts the Docker Compose yaml file for the application.

A Makefile is present to add some convenience aliases. See the Makefile for all commands.

Air is used as a hot-reload support tool for development. It helps when quickly making changes in the code and automatically watches for changes, builds a binary, and then runs it.

A Dockerfile to build a container image is available.


## Getting Started 

### Prerequisites

You must have the following installed:

- Go@v1.26.1 (minimum)
- Docker (or alternative, using the `docker` alias so the Makefile can be used if you use something like Podman or Colima)

### Project Dependenices

There are no required project dependencies to download as everything was created using the standard Go library. 

I recommend the following optional dependencies:

- [Air](https://github.com/air-verse/air) for fast-reload development

```bash
go install github.com/air-verse/air@latest
```

Using Nix? A flake is included for repeatable builds. Run:

```bash
nix develop
```

## Running the application

### Using Go

```bash
make run
```

### Using Air

```bash 
make start # starts up a hot-reload server
```

### Using Docker:

This will build and run the Docker container

```bash
make docker-run
```

### Running tests
```bash
make test
```

### Assumptions Made & Improvements I would make

I had to make a couple of assumptions:

1. Only the lowest out of the entire 24 hour period should be returned, even if they are not adjacent.
2. When asking specifiying `continuous=true` there should be no gap between the time periods processed. 

I have listed a few of the things I would do if I had more time:

1. The averaging functions are a little messy and I would refactor these to be able to try and utilise the same averaging functions across multliple calculations. 
2. I only added a logging middleware, so adding extra middleware would be required
3. There are many more tests that I could write but I wanted to add some simple ones to validate the calculate functions
4. There are a couple of situations I could extract code from functions to make it less tightly coupled


### GET /slots

The `/slots` endpoint accepts two query parameters: 
```
duration (number): 1 <= x <= 1440: default = 30

continuous (bool): true|false: default = false
```

The `/slots` endpoint defaults to a `duration=30` and `continuous=30` 

Here are some example test commands:

```bash
# defaults to duration=30 and continuous = false
curl --H 'application/json' --url "http://localhost:7777/slots"

curl --H 'application/json' --url "http://localhost:7777/slots?duration=30&continuous=true"

# returns 400 Bad Request
curl --H 'application/json' --url "http://localhost:7777/slots?duration=1500&continuous=true" 
curl --H 'application/json' --url "http://localhost:7777/slots?duration=-1&continuous=true" 

# defaults to duration=30 and continuous = false
curl --H 'application/json' --url "http://localhost:7777/slots?duration=0" 

```

