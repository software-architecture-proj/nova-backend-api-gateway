# Nova Backend API Gateway

The API Gateway service for the Nova Backend microservices architecture. This service acts as an intermediary between the frontend and the various backend microservices.

## Services

The API Gateway interfaces with the following microservices:

1. User Product Service (`:50051`)
   - User management
   - Favorites management
   - Pockets management
   - Verifications
   - Country codes

2. Auth Service (`:50052`)
   - User authentication

3. Transaction Service (`:50053`)
   - Account management
   - Transfers
   - Balance queries
   - Movement history

## Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and adjust values as needed
3. Install dependencies: `go mod download`
4. Run the service: `go run main.go`

## Environment Variables

- `API_GATEWAY_PORT`: Port for the API Gateway (default: 8080)
- `USER_PRODUCT_SERVICE_GRPC_HOST`: User Product Service gRPC endpoint
- `AUTH_SERVICE_GRPC_HOST`: Auth Service gRPC endpoint
- `TRANSACTION_SERVICE_GRPC_HOST`: Transaction Service gRPC endpoint

## API Endpoints

### User Management
- `GET /api/country-codes` - Get list of country codes
- `POST /api/users` - Create new user
- `GET /api/users/{user_id}` - Get user details
- `PUT /api/users/{user_id}` - Update user
- `DELETE /api/users/{user_id}` - Delete user

### Favorites
- `GET /api/users/{user_id}/favorites` - Get user's favorites
- `POST /api/users/{user_id}/favorites` - Add favorite
- `PUT /api/users/{user_id}/favorites/{favorite_id}` - Update favorite
- `DELETE /api/users/{user_id}/favorites/{favorite_id}` - Delete favorite

### Pockets
- `GET /api/users/{user_id}/pockets` - Get user's pockets
- `POST /api/users/{user_id}/pockets` - Create pocket
- `PUT /api/users/{user_id}/pockets/{pocket_id}` - Update pocket
- `DELETE /api/users/{user_id}/pockets/{pocket_id}` - Delete pocket

### Verifications
- `GET /api/users/{user_id}/verifications` - Get user's verifications
- `PUT /api/users/{user_id}/verifications` - Update verification status

### Authentication
- `POST /api/login` - User login

### Transactions
- `POST /api/accounts` - Create account
- `POST /api/transfers` - Make transfer
- `GET /api/balance` - Get balance
- `GET /api/movements` - Get transaction history