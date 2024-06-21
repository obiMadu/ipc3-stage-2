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
  - [6. Some Additional Notes](#6-some-additional-notes)

---

## 1. Introduction

This project is a simple REST compliant API written in Go.

## 2. API Documentation

This API is also documented [here](https://documenter.getpostman.com/view/29936566/2sA3XV9KXa) on Postman.

### 2.1 How to Call the API

The API can be accessed via HTTP requests. It exposes endpoints for various CRUD operations.

Sample API base URL: `http://example.com/api`

**Current [Active]** base URL: `https://ip.obi.ninja/api`

### 2.2 Supported CRUD Operations

The API supports the following CRUD operations:

- **CREATE**: `POST /api/users`
- **READ**: `GET /api` & `GET /api/{userID}` & `GET /api?username={username}`
- **UPDATE**: `PUT /api/{userID}` & `PUT /api?username={username}`
- **DELETE**: `DELETE /api{userID}` & `DELETE /api?username={username}`

## 3. Request and Response Formats

### 3.1 Request Formats

- **GET Request:** `GET` `/api` | `GET` `/api/{userID}` | `GET` `/api?username={username}`
  - Body (no-data)
  
- **CREATE Request:** `POST` `/api`
  - Body (Json):
    - name (string, required, **must-be-unique**): This is the Username of the new User.
    - email (string, required, **must-be-unique**): The Email address of User the new User.
    - fullname (string, optional): Optional Fullname of the User.
    
- **UPDATE Request:** `PUT` `/api/{userID}` | `PUT` `/api?username={username}`
  - Body (Json): Only any one of the following fields is required
    - name (string, **must-be-unique**): This is the new Username for the User.
    - email (string, **must-be-unique**): The new Email address for the User.
    - fullname (string): The new Fullname for the user.

- **DELETE Request** `DELETE` `/api/{userID}` | `DELETE` `/api?username={username}`
  - Body (no-data)

### 3.2 Response Formats

- **REST Compliant** (Success Code: 200 OK, Error Codes: 4xx or 5xx)
  - Body (Text):
    - status (string) (success,error): Status of the response.
    - message (string): Summary of the response.
    - data (object): Information returned by the API.
    - error (object): (Mostly always nil) Rare error object retuned at server errors.

## 4. Sample API Calls

Visit the Postman Documentation [here](https://documenter.getpostman.com/view/29936566/2sA3XV9KXa) to view sample `success` and `error` requests.
    ```

## 5. Setting Up and Running the API (Locally or othewise)

### 5.1 Environment Variables

To run the API via Docker Compose, remember to set the following environment variables on your compose file:

- `POSTGRES_DSN`: The DSN string for the PosgreSQL database connection.

### 5.2 Docker Compose Setup

I've created and attached a Docker compose file containing instructions for building an image for this API, and a PostgresSQL database with default configurations. This file should be enough to quickly get the API up and running.

At the project root run:
   
```sh
docker compose up --build
```

The API should now be accessible at `http://localhost:8080`


## 6. Some Additional Notes

- This repository contains Github Actions worflows for Continuous Integration.
