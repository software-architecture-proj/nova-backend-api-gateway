# Nova Backend API Gateway Documentation

This document describes all available endpoints in the Nova Backend API Gateway, including their request and response formats.

## Base URL
```
http://localhost:8080
```

## Authentication

### Login
Authenticate a user and get their session token.

**Endpoint:** `POST /login`

**Request Body:**
```json
{
    "email": "string",
    "password": "string"
}
```

**Response:**
```json
{
    "success": boolean,
    "message": "string",
    "email": "string"
}
```

## User Management

### Create User
Create a new user account.

**Endpoint:** `POST /users`

**Request Body:**
```json
{
    "email": "string",
    "username": "string",
    "code_id": "string",
    "phone": "string",
    "first_name": "string",
    "last_name": "string",
    "birthdate": "string"
}
```

**Response:**
```json
{
    "success": boolean,
    "message": "string",
    "user_id": "string"
}
```

### Get User
Retrieve user information.

**Endpoint:** `GET /users/{user_id}`

**Response:**
```json
{
    "user_id": "string",
    "email": "string",
    "username": "string",
    "phone": "string",
    "first_name": "string",
    "last_name": "string",
    "birthdate": "string"
}
```

### Update User
Update user information.

**Endpoint:** `PUT /users/{user_id}`

**Request Body:**
```json
{
    "email": "string",
    "username": "string",
    "phone": "string",
    "first_name": "string",
    "last_name": "string",
    "birthdate": "string"
}
```

**Response:**
```json
{
    "success": boolean,
    "message": "string"
}
```

### Delete User
Delete a user account.

**Endpoint:** `DELETE /users/{user_id}`

**Response:**
```json
{
    "success": boolean,
    "message": "string"
}
```

## Favorites Management

### Get User Favorites
Get all favorites for a user.

**Endpoint:** `GET /users/{user_id}/favorites`

**Response:**
```json
{
    "favorites": [
        {
            "id": "string",
            "favorite_user_id": "string",
            "alias": "string"
        }
    ]
}
```

### Create Favorite
Add a new favorite for a user.

**Endpoint:** `POST /users/{user_id}/favorites`

**Request Body:**
```json
{
    "favorite_user_id": "string",
    "alias": "string"
}
```

**Response:**
```json
{
    "success": boolean,
    "message": "string",
    "favorite_id": "string"
}
```

### Update Favorite
Update a favorite's information.

**Endpoint:** `PUT /users/{user_id}/favorites/{favorite_id}`

**Request Body:**
```json
{
    "alias": "string"
}
```

**Response:**
```json
{
    "success": boolean,
    "message": "string"
}
```

### Delete Favorite
Remove a favorite.

**Endpoint:** `DELETE /users/{user_id}/favorites/{favorite_id}`

**Response:**
```json
{
    "success": boolean,
    "message": "string"
}
```

## Pockets Management

### Get User Pockets
Get all pockets for a user.

**Endpoint:** `GET /users/{user_id}/pockets`

**Response:**
```json
{
    "pockets": [
        {
            "id": "string",
            "name": "string",
            "category": "string",
            "max_amount": number
        }
    ]
}
```

### Create Pocket
Create a new pocket for a user.

**Endpoint:** `POST /users/{user_id}/pockets`

**Request Body:**
```json
{
    "name": "string",
    "category": "string",
    "max_amount": number
}
```

**Response:**
```json
{
    "success": boolean,
    "message": "string",
    "pocket_id": "string"
}
```

### Update Pocket
Update a pocket's information.

**Endpoint:** `PUT /users/{user_id}/pockets/{pocket_id}`

**Request Body:**
```json
{
    "name": "string",
    "category": "string",
    "max_amount": number
}
```

**Response:**
```json
{
    "success": boolean,
    "message": "string"
}
```

### Delete Pocket
Delete a pocket.

**Endpoint:** `DELETE /users/{user_id}/pockets/{pocket_id}`

**Response:**
```json
{
    "success": boolean,
    "message": "string"
}
```

## Transactions

### Get Movements
Get user's transaction movements.

**Endpoint:** `GET /movements`

**Query Parameters:**
- `user_id`: string (required)
- `from_time`: number (unix timestamp)
- `to_time`: number (unix timestamp)
- `limit`: boolean

**Response:**
```json
{
    "movements": [
        {
            "id": "string",
            "from_user": "string",
            "to_user": "string",
            "amount": number,
            "timestamp": number
        }
    ]
}
```

### Get Balance
Get user's balance.

**Endpoint:** `GET /balance`

**Query Parameters:**
- `user_id`: string (required)
- `from_time`: number (unix timestamp)
- `to_time`: number (unix timestamp)

**Response:**
```json
{
    "balance": number,
    "currency": "string"
}
```

### Create Account
Create a new account.

**Endpoint:** `POST /account`

**Request Body:**
```json
{
    "username": "string",
    "bank": boolean,
    "user_id": "string"
}
```

**Response:**
```json
{
    "success": boolean,
    "message": "string",
    "account_id": "string"
}
```

### Transfer Funds
Make a transfer between accounts.

**Endpoint:** `POST /transfer`

**Request Body:**
```json
{
    "from_user": "string",
    "to_user": "string",
    "amount": number
}
```

**Response:**
```json
{
    "success": boolean,
    "message": "string",
    "transaction_id": "string"
}
```

## Utility Endpoints

### Get Country Codes
Get available country codes.

**Endpoint:** `GET /country-codes`

**Response:**
```json
{
    "codes": [
        {
            "id": "string",
            "name": "string",
            "code": "string"
        }
    ]
}
```

## Error Responses

All endpoints may return error responses in the following format:

```json
{
    "error": "string"
}
```

Common HTTP status codes:
- 200: Success
- 201: Created
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 500: Internal Server Error 