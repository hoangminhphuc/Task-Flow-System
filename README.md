# Task-Flow-System

A sample Go project demonstrating Clean Architecture with user authentication, registration, and API documentation. This project uses Gin for HTTP handling, GORM for database operations, bcrypt for password hashing, and JWT for token generation.

## Overview

This project implements a RESTful API built on Clean Architecture principles, which divides the code into distinct layers:

- **Biz (Business Logic):** Contains core domain logic such as user registration, login, and password verification.
- **Storage (Data Access):** Handles database operations using GORM.
- **Transport (HTTP Handlers):** Manages HTTP endpoints with Gin.
- **Common:** Contains shared utilities, error handling, and token payload definitions.

The system also supports **Pub/Sub messaging** for asynchronous event handling and **Job & Job Groups** for task scheduling and execution

# Getting Started

## Prerequisites

Ensure you have the following installed:

- **Go** (>=1.18)
- **Docker** (for database setup)
- **MySQL** 

## Installation

### Clone the repository:

```sh
git clone https://github.com/hoangminhphuc/Task-Flow-System.git
cd Task-Flow-System
```

### Copy and configure environment variables:

```sh
cp .env.example .env
```

### Install dependencies:

```sh
go mod tidy
```

## Running the Application

### Run the application in development mode:

```sh
go build && ./Task-Flow-System
```

The server should start on [http://localhost:8080](http://localhost:8080).

# Running with Docker

## Setting up MySQL in Docker

Before running the system container, you need to start a MySQL container and ensure both containers are in the same network.

### 1. Create a Docker Network (if not already created):
```sh
docker network create your-network-name
```

### 2. Start the MySQL Container:
```sh
docker run -d \
  --name your-mysql-container-name \   # Replace with your MySQL container name
  --network your-network-name \        # Must match the network used by the system container
  -e MYSQL_ROOT_PASSWORD=your-password \  # Replace with your MySQL root password
  mysql:latest
```

## Building the Docker Image

From the project root directory, build your Docker image. Replace `task-management-system:1.0` with your desired image name and tag if needed:

```sh
docker build -t task-management-system:1.0 .
```

## Single Container Deployment

Use the command below to run your Docker container. **Ensure it uses the same network as MySQL. Replace the placeholder values with your own configuration details.**

```sh
docker run -d \
  --name your-container-name \   # Replace with your container name
  --network your-network-name \   # Replace with your Docker network name
  -p HOST_PORT:CONTAINER_PORT \   # Replace with your desired port mapping (e.g., 8080:3000)
  -e ITEM_SERVICE_URL="http://your-item-service:PORT" \  # Replace with your item service URL and port
  -e MYSQL_GORM_DB_TYPE=mysql \
  -e MYSQL_GORM_DB_URI="your-username:your-password@tcp(your-mysql-host:3306)/your-database?charset=utf8mb4&parseTime=True&loc=Local" \  # Update with your MySQL credentials and details
  -e SECRET="your-secret" \  # Replace with your secret key
  task-management-system:1.0
```


