# Project Inventaris Golang
## Overview
This project is a Golang-based inventory management system that allows users to manage categories and items, track investments in items, and handle user authentication. The application uses PostgreSQL as the database and follows a structured architecture with separate layers for models, repositories, services, and handlers.

## Features
- User Registration and Login
- Category Management (Create, Update, Delete, Retrieve)
- Item Management (Create, Update, Delete, Retrieve)
- Investment Tracking for Items
- Session Management for User Authentication
- File Uploads for Item Photos
- Replacement Reminder for Items based on usage
  
## Technologies Used
- Go (Golang)
- PostgreSQL
- Chi Router for HTTP routing
- JSON for data interchange
- Bcrypt for password hashing

## Getting Started
### Prerequisites
- Go (version 1.16 or later)
- PostgreSQL (version 12 or later)
- Git

### Installation
1. Clone the repository:
```
git clone https://github.com/yourusername/project-app-inventaris-golang-safira.git
cd project-app-inventaris-golang-safira
```
2. Create a PostgreSQL database:
```
CREATE DATABASE inventaris;
```
3. Run the SQL script to set up the database tables:
```
psql -U postgres -d inventaris -f inventaris.sql
```
4. Update the database connection string in database/postgres.go if necessary.
5. Install dependencies:
```
go mod tidy
```
6. Run the application:
```
go run main.go
```
7. The server will start on http://localhost:8080.

## API Endpoints
### Authentication
- POST /api/auth/register: Register a new user.
  Request Body:
  ```
  {
    "username": "john_doe",
    "email": "john@example.com",
    "password": "securePassword123"
  }
  ```
- POST /api/auth/login: Login an existing user.
  Request Body:
  ```
  {
    "username": "john_doe",
    "password": "securePassword123"
  }
  ```
### Categories
- GET /api/categories: Retrieve all categories.
  _No request body is needed for this endpoint._
- GET /api/categories/{id}: Retrieve a category by ID.
  _No request body is needed for this endpoint; the ID is passed in the URL._

- POST /api/categories: Create a new category.
  Request Body:
  ```
  {
    "name": "Electronics",
    "description": "Devices and gadgets"
  }
  ```
- PUT /api/categories/{id}: Update an existing category.
  Request Body:
  ```
  {
    "name": "Home Appliances",
    "description": "Appliances used in the home"
  }
  ```
- DELETE /api/categories/{id}: Delete a category.
  _No request body is needed for this endpoint; the ID is passed in the URL._

### Items
- GET /api/items: Retrieve all items.
  _No request body is needed for this endpoint._
- GET /api/items/{id}: Retrieve an item by ID.
  _No request body is needed for this endpoint; the ID is passed in the URL._
- POST /api/items: Create a new item.
  Request Body:
  ```
  {
    "name": "Laptop",
    "category_id": 1,
    "price": 1500.00,
    "purchase_date": "2023-01-15",
    "depreciated_rate": 20,
    "photo": "data:image/jpeg;base64,..."
  }
  ```
  _Note: The photo field should contain the file data in base64 format. In a real application, this would typically be handled as a multipart form upload._
- PUT /api/items/{id}: Update an existing item.
  Request Body:
  ```
  {
    "name": "Gaming Laptop",
    "category_id": 1,
    "price": 2000.00,
    "purchase_date": "2023-01-15",
    "depreciated_rate": 15,
    "photo": "data:image/jpeg;base64,..."
  }
  ```
- DELETE /api/items/{id}: Delete an item.
  _No request body is needed for this endpoint; the ID is passed in the URL._
- GET /api/items/need-replacement: Retrieve items that need replacement.
  _No request body is needed for this endpoint ; the ID is passed in the URL._
### Investment Tracking
- GET /api/items/investment: Count all item investments.
  _No request body is needed for this endpoint._
- GET /api/items/investment/{id}: Get investment details for a specific item by ID.
  _No request body is needed for this endpoint; the ID is passed in the URL._

## Conclusion
This README provides an overview of the project, its features, and how to interact with the API. For further details, please refer to the codebase or reach out for assistance.
