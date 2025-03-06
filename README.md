# Task-Flow-System

A sample Go project demonstrating Clean Architecture with user authentication, registration, and API documentation. This project uses Gin for HTTP handling, GORM for database operations, bcrypt for password hashing, and JWT for token generation.

## Overview

This project implements a RESTful API built on Clean Architecture principles, which divides the code into distinct layers:

- **Biz (Business Logic):** Contains core domain logic such as user registration, login, and password verification.
- **Storage (Data Access):** Handles database operations using GORM.
- **Transport (HTTP Handlers):** Manages HTTP endpoints with Gin.
- **Common:** Contains shared utilities, error handling, and token payload definitions.

## Features

- **User Registration & Login:** Secure authentication with password hashing (bcrypt) and JWT token generation.
- **Clean Architecture:** Separation of concerns between business logic, data access, and presentation layers.
- **Database Interaction:** Uses GORM to simplify CRUD operations on your SQL database.
- **Robust Error Handling:** Custom error types and logging to simplify debugging and maintenance.

## Architecture

