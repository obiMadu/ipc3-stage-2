# SIMPLE CRUD REST API
## InternPulse Cohort 3, Stage 2 Project

## Table of Contents

- [SIMPLE CRUD REST API](#simple-crud-rest-api)
  - [InternPulse Cohort 3, Stage 2 Project](#internpulse-cohort-3-stage-2-project)
  - [Table of Contents](#table-of-contents)
  - [1. Introduction](#1-introduction)
  - [2. API Documentation](#2-api-documentation)
    - [2.1 How to Call the API](#21-how-to-call-the-api)
    - [2.2 Supported CRUD Operations](#22-supported-crud-operations)
  - [3. Request and Response Formats](#3-request-and-response-formats)
    - [3.1 Request Formats](#31-request-formats)
    - [3.2 Response Formats](#32-response-formats)
  - [4. Sample API Calls](#4-sample-api-calls)
  - [5. Setting Up and Running the API (Locally or othewise)](#5-setting-up-and-running-the-api-locally-or-othewise)
    - [5.1 Environment Variables](#51-environment-variables)
    - [5.2 Docker Compose Setup](#52-docker-compose-setup)
    - [5.3 Run API Locally](#53-run-api-locally)
  - [6. Some Additional Notes](#6-some-additional-notes)

---

## 1. Introduction

This project is a simple REST compliant API written in Go.

## 2. API Documentation

This API is also documented [here](https://documenter.getpostman.com/view/29936566/2sA3XV9KXa) on Postman.

### 2.1 How to Call the API

The API can be accessed via HTTP requests. It exposes endpoints for various CRUD operations.

Sample API base URL: `http://example.com/api`

**Current [Active]** base URL: `https://ips2.obi.ninja/api`

### 2.2 Supported CRUD Operations

The API supports the following CRUD operations:

- **CREATE**: `POST /users`
- **READ**: `GET /users` & `GET /users/{userID}` & `GET /users?username={username}`
- **UPDATE**: `PUT /users/{userID}` & `PUT /users?username={username}`
- **DELETE**: `DELETE /users/{userID}` & `DELETE /users?username={username}`

## 3. Request and Response Formats

### 3.1 Request Formats

- **GET Request:** `GET` `/users` | `GET` `/users/{userID}` | `GET` `/users?username={username}`
  - Body (no-data)
  
- **CREATE Request:** `POST` `/users`
  - Body (Json):
    - name (string, required, **must-be-unique**): This is the Username of the new User.
    - email (string, required, **must-be-unique**): The Email address of User the new User.
    - fullname (string, optional): Optional Fullname of the User.
    
- **UPDATE Request:** `PUT` `/users/{userID}` | `PUT` `/users?username={username}`
  - Body (Json): Only any one of the following fields is required
    - name (string, **must-be-unique**): This is the new Username for the User.
    - email (string, **must-be-unique**): The new Email address for the User.
    - fullname (string): The new Fullname for the user.

- **DELETE Request** `DELETE` `/users/{userID}` | `DELETE` `/users?username={username}`
  - Body (no-data)

### 3.2 Response Formats

- **REST Compliant** (Success Code: 200 OK, Error Codes: 4xx or 5xx)
  - Body (Text):
    - status (string) (success,error): Status of the response.
    - message (string): Summary of the response.
    - data (object): Information returned by the API.
    - error (object): (Mostly always nil) Rare error object retuned at server errors.
  - Error Codes
    - The API returns 4xx errors for bad, malformed, incomplete or improper requests.
    - The API returns 5xx errors for server errors.

## 4. Sample API Calls

Visit the Postman Documentation [here](https://documenter.getpostman.com/view/29936566/2sA3XV9KXa) to view sample `success` and `error` requests.

## 5. Setting Up and Running the API (Locally or othewise)

### 5.1 Environment Variables

To run the API locally or via Docker Compose, remember to set the following environment variables. If working locally you can create a `.env` file with a `key:value` format containing the variables below and their values, the program will automatically pick those up at run time. On Docker, set these variables on your Compose file, the `docker-compose.yml` in the project source has good examples.

- `POSTGRES_DSN`: The DSN string for the PosgreSQL database connection. It's of the format `"host=localhost port=5432 user=postgres password=password dbname=users sslmode=disable"`.

`!important:` When setting environment variables on your Docker Compose file, do NOT enclose the variable values in quotes, EVEN IF said value contains spaces. Docker Compose will add the quotes as part of your string, causing confusion for the program.

### 5.2 Docker Compose Setup

I've created and attached a Docker compose file containing instructions for building an image for this API, and a PostgresSQL database with default configurations. This file should be enough to quickly get the API up and running.

At the project root run:
   
```sh
docker compose up --build
```

The API should now be accessible at `http://localhost:8080`

### 5.3 Run API Locally
If you have a Postgres database setup locally, you can head over to the [Release page](https://github.com/obiMadu/ipc3-stage-2/releases), download the binary for your operating system and run the API.


## 6. Some Additional Notes

- This repository contains Github Actions workflows for Continuous Integration.
