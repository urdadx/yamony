# Authentication API Documentation

This document describes the authentication endpoints for the Yamony API.

## Authentication Endpoints

### Register a New User

**POST** `/auth/register`

Creates a new user account.

**Request Body:**
```json
{
  "firstName": "John",
  "lastName": "Doe",
  "email": "john.doe@example.com",
  "password": "securePassword123"
}
```

**Success Response (201 Created):**
```json
{
  "message": "user registered successfully",
  "user": {
    "id": 1,
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@example.com",
    "emailVerified": false,
    "image": ""
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body or validation errors
- `409 Conflict`: User with this email already exists

---

### Login

**POST** `/auth/login`

Authenticates a user and creates a session.

**Request Body:**
```json
{
  "email": "john.doe@example.com",
  "password": "securePassword123"
}
```

**Success Response (200 OK):**
```json
{
  "message": "login successful",
  "user": {
    "id": 1,
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@example.com",
    "emailVerified": false,
    "image": ""
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Invalid email or password

**Note:** A session cookie named `yamony_session` will be set upon successful login.

---

### Logout

**POST** `/auth/logout`

Logs out the current user and destroys their session.

**Success Response (200 OK):**
```json
{
  "message": "logout successful"
}
```

---

### Get Current User

**GET** `/api/me`

Returns the currently authenticated user's information.

**Authentication Required:** Yes

**Success Response (200 OK):**
```json
{
  "user": {
    "id": 1,
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@example.com",
    "emailVerified": false,
    "image": ""
  }
}
```

**Error Response:**
- `401 Unauthorized`: Not authenticated or session expired

---

## Authentication Middleware

### Protected Routes

All routes under `/api/*` require authentication. The middleware checks for a valid session cookie and validates it against the database.

If the session is invalid or expired, the middleware will return a `401 Unauthorized` response.

### Session Details

- **Cookie Name:** `yamony_session`
- **Duration:** 7 days
- **HttpOnly:** true
- **Secure:** false (set to true in production with HTTPS)
- **SameSite:** Lax

---

## Testing with cURL

### Register
```bash
curl -X POST http://localhost:3000/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@example.com",
    "password": "securePassword123"
  }'
```

### Login
```bash
curl -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{
    "email": "john.doe@example.com",
    "password": "securePassword123"
  }'
```

### Get Current User (Protected)
```bash
curl -X GET http://localhost:3000/api/me \
  -b cookies.txt
```

### Logout
```bash
curl -X POST http://localhost:3000/auth/logout \
  -b cookies.txt
```

---

## Database Schema

### Users Table
```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  email_verified BOOLEAN NOT NULL DEFAULT FALSE,
  image VARCHAR(500) NOT NULL DEFAULT '',
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### Sessions Table
```sql
CREATE TABLE sessions (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  session_token VARCHAR(255) NOT NULL UNIQUE,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

---

## Security Notes

1. **Password Hashing:** Passwords are hashed using bcrypt with default cost (10)
2. **Session Tokens:** Generated using cryptographically secure random bytes (32 bytes, base64 encoded)
3. **Session Expiration:** Sessions expire after 7 days
4. **CORS:** Configured to allow requests from `http://localhost:3001` and `https://yamony.com`
5. **Production:** Remember to:
   - Change the session secret key (`secret-key-change-this-in-production`)
   - Set `Secure: true` for session cookies when using HTTPS
   - Use environment variables for sensitive configuration
