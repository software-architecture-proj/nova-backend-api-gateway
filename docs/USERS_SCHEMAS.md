# API Gateway Schema Documentation for User Service

This document outlines the request and response schemas for the API gateway service.

## Authentication Endpoints

### `/login`

*   **Method:** `POST`
*   **Description:** Authenticates a user.
*   **Request Body:**
    ```json
    {
        "email": "user@example.com",
        "password": "securepassword123"
    }
    ```
*   **Response Body (Success):**
    ```json
    {
        "success": true,
        "message": "Login successful",
        "email": "user@example.com"
    }
    ```
*   **Response Body (Failure):**
    ```json
    {
        "success": false,
        "message": "Invalid credentials",
        "email": "user@example.com" 
    }
    ```

## User and Product Endpoints

All endpoints are prefixed with `/api`.

### Country Codes

#### `GET /country-codes`

*   **Method:** `GET`
*   **Description:** Retrieves a list of country codes.
*   **Request Body:** None
*   **Response Body (Success):**
    ```json
    {
        "success": true,
        "message": "Country codes retrieved successfully",
        "country_codes": [
            {
                "id": "some-uuid",
                "name": "United States",
                "code": "+1"
            },
            {
                "id": "another-uuid",
                "name": "Colombia",
                "code": "+57"
            }
        ]
    }
    ```

### Users

#### `POST /users`

*   **Method:** `POST`
*   **Description:** Creates a new user.
*   **Request Body:**
    ```json
    {
        "email": "user@example.com",
        "username": "newuser",
        "code_id": "country-code-uuid",
        "phone": "1234567890",
        "first_name": "John",
        "last_name": "Doe",
        "birthdate": "YYYY-MM-DD"
    }
    ```
*   **Response Body (Success):**
    ```json
    {
        "success": true,
        "message": "User created successfully",
        "user_id": "user-uuid"
    }
    ```

#### `GET /users/{user_id}`

*   **Method:** `GET`
*   **Description:** Retrieves a user by their ID.
*   **Path Parameters:**
    *   `user_id` (string, required): The ID of the user.
*   **Request Body:** None
*   **Response Body (Success):**
    ```json
    {
        "success": true,
        "message": "User retrieved successfully",
        "email": "user@example.com",
        "username": "existinguser",
        "phone": "+11234567890",
        "first_name": "Jane",
        "last_name": "Doe",
        "birthdate": "YYYY-MM-DD"
    }
    ```

#### `PUT /users/{user_id}`

*   **Method:** `PUT`
*   **Description:** Updates an existing user.
*   **Path Parameters:**
    *   `user_id` (string, required): The ID of the user to update.
*   **Request Body:**
    ```json
    {
        "email": "updateduser@example.com",
        "username": "updatedusername",
        "phone": "+10987654321",
        "first_name": "Johnathan",
        "last_name": "Doeboy",
        "birthdate": "YYYY-MM-DD"
    }
    ```
*   **Response Body (Success):**
    ```json
    {
        "success": true,
        "message": "User updated successfully",
        "email": "updateduser@example.com",
        "username": "updatedusername",
        "phone": "+10987654321",
        "first_name": "Johnathan",
        "last_name": "Doeboy",
        "birthdate": "YYYY-MM-DD" 
    }
    ```

#### `DELETE /users/{user_id}`

*   **Method:** `DELETE`
*   **Description:** Deletes a user by their ID.
*   **Path Parameters:**
    *   `user_id` (string, required): The ID of the user to delete.
*   **Request Body:** None
*   **Response Body (Success):**
    ```json
    {
        "success": true,
        "message": "User deleted successfully"
    }
    ```

### Favorites

#### `GET /users/{user_id}/favorites`

*   **Method:** `GET`
*   **Description:** Retrieves a user's list of favorite users.
*   **Path Parameters:**
    *   `user_id` (string, required): The ID of the user.
*   **Request Body:** None
*   **Response Body (Success):**
    ```json
    {
        "success": true,
        "message": "Favorites retrieved successfully",
        "favorites": [
            {
                "id": "favorite-uuid",
                "user_id": "user-uuid",
                "favorite_user_id": "another-user-uuid",
                "favorite_username": "favuser1",
                "alias": "My Best Friend"
            }
        ]
    }
    ```

#### `POST /users/{user_id}/favorites`

*   **Method:** `POST`
*   **Description:** Adds a user to the favorites list.
*   **Path Parameters:**
    *   `user_id` (string, required): The ID of the user.
*   **Request Body:**
    ```json
    {
        "favorite_user_id": "user-to-favorite-uuid",
        "alias": "Work Colleague"
    }
    ```
*   **Response Body (Success):**
    ```json
    {
        "success": true,
        "message": "Favorite added successfully",
        "favorite_id": "new-favorite-uuid"
    }
    ```

#### `PUT /users/{user_id}/favorites/{favorite_id}`

*   **Method:** `PUT`
*   **Description:** Updates the alias of a favorite user.
*   **Path Parameters:**
    *   `user_id` (string, required): The ID of the user.
    *   `favorite_id` (string, required): The ID of the favorite entry.
*   **Request Body:**
    ```json
    {
        "alias": "Close Friend"
    }
    ```
*   **Response Body (Success):**
    ```json
    {
        "success": true,
        "message": "Favorite updated successfully",
        "new_alias": "Close Friend"
    }
    ```

#### `DELETE /users/{user_id}/favorites/{favorite_id}`

*   **Method:** `DELETE`
*   **Description:** Removes a user from the favorites list.
*   **Path Parameters:**
    *   `user_id` (string, required): The ID of the user.
    *   `favorite_id` (string, required): The ID of the favorite entry.
*   **Request Body:** None
*   **Response Body (Success):**
    ```json
    {
        "success": true,
        "message": "Favorite deleted successfully"
    }
    ```

### Pockets

#### `GET /users/{user_id}/pockets`

*   **Method:** `GET`
*   **Description:** Retrieves a user's list of pockets.
*   **Path Parameters:**
    *   `user_id` (string, required): The ID of the user.
*   **Request Body:** None
*   **Response Body (Success):**
    ```json
    {
        "success": true,
        "message": "Pockets retrieved successfully",
        "pockets": [
            {
                "id": "pocket-uuid",
                "user_id": "user-uuid",
                "name": "Vacation Fund",
                "category": "Savings",
                "max_amount": 5000
            }
        ]
    }
    ```

#### `POST /users/{user_id}/pockets`

*   **Method:** `POST`
*   **Description:** Creates a new pocket for a user.
*   **Path Parameters:**
    *   `user_id` (string, required): The ID of the user.
*   **Request Body:**
    ```json
    {
        "name": "Emergency Fund",
        "category": "Savings",
        "max_amount": 10000
    }
    ```
*   **Response Body (Success):**
    ```json
    {
        "success": true,
        "message": "Pocket created successfully",
        "pocket_id": "new-pocket-uuid"
    }
    ```

#### `PUT /users/{user_id}/pockets/{pocket_id}`

*   **Method:** `PUT`
*   **Description:** Updates an existing pocket.
*   **Path Parameters:**
    *   `user_id` (string, required): The ID of the user.
    *   `pocket_id` (string, required): The ID of the pocket to update.
*   **Request Body:**
    ```json
    {
        "name": "New Car Fund",
        "category": "Savings",
        "max_amount": 20000
    }
    ```
*   **Response Body (Success):**
    ```json
    {
        "success": true,
        "message": "Pocket updated successfully",
        "name": "New Car Fund",
        "category": "Savings",
        "max_amount": 20000
    }
    ```

#### `DELETE /users/{user_id}/pockets/{pocket_id}`

*   **Method:** `DELETE`
*   **Description:** Deletes a pocket by its ID.
*   **Path Parameters:**
    *   `user_id` (string, required): The ID of the user.
    *   `pocket_id` (string, required): The ID of the pocket to delete.
*   **Request Body:** None
*   **Response Body (Success):**
    ```json
    {
        "success": true,
        "message": "Pocket deleted successfully"
    }
    ```
