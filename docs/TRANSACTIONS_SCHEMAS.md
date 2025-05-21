# API Gateway Schema Documentation for Transactions Service

This document outlines the request and response schemas for the Transactions service, accessed via the API gateway.

All endpoints are assumed to be prefixed with a base path for the transactions service (e.g., `/api/transactions`).

## Account Management

### `POST /accounts` (Assumed Endpoint: `PostAccount`)

*   **Method:** `POST`
*   **Description:** Creates a new account for a user. (Assuming based on `PostAccount` key)
*   **Request Body:**
    ```json
    {
        "user_id": "string", 
        "username": "string", 
        "bank": "bool"
    }
    ```
*   **Response Body (Success):**
    ```json
    {
        "success": "bool",
        "message": "string",
        "user_id": "string", 
        "timestamp": "string"
    }
    ```

## Transfers

### `POST /transfers` (Assumed Endpoint: `PostTransfer`)

*   **Method:** `POST`
*   **Description:** Creates a new transfer between users. (Assuming based on `PostTransfer` key)
*   **Request Body:**
    ```json
    {
        "from_user_id": "string", 
        "to_user_id": "string", 
        "amount": "uint64"
    }
    ```
*   **Response Body (Success):**
    ```json
    {
        "success": "bool",
        "message": "string",
        "transfer_id": "string", 
        "timestamp": "string"
    }
    ```

## Balance Inquiry

### `GET /balance` (Assumed Endpoint: `GetBalance`) 

*   **Method:** `GET` (Common for retrieving data)
*   **Description:** Retrieves the account balance for a user within a specified time range. (Assuming based on `GetBalance` key)
*   **Request Body:**
    ```json
    {
        "user_id": "string",
        "from_time": "uint64", 
        "to_time": "uint64"
    }
    ```
*   **Response Body (Success):**
    ```json
    {
        "success": "bool",
        "message": "string",
        "timestamp": "string",
        "current": "string",
        "balances": [
            {
                "income": "string",
                "outcome": "string"
            }
        ]
    }
    ```

## Movement History

### `GET /movements` (Assumed Endpoint: `GetMovements`)

*   **Method:** `GET` (Common for retrieving data)
*   **Description:** Retrieves the transaction movements for a user. (Assuming based on `GetMovements` key)
*   **Request Body:**
    ```json
    {
        "user_id": "string",
        "from_time": "uint64",
        "to_time": "uint64",
        "limit": "bool"
    }
    ```
*   **Response Body (Success):**
    ```json
    {
        "success": "bool",
        "message": "string",
        "movements": [
            {
                "transferId": "string",
                "fromUsername": "string",
                "toUsername": "string",
                "amount": "string",
                "timestamp": "string"
            }
        ]
    }
    ```
