# Headplane Users API Reference

This document provides technical API documentation for the Headplane Users management endpoints introduced in Phase 13c. These endpoints allow administrators to manage Headplane dashboard users, including creating, listing, updating, and deleting users.

## Overview

All user management endpoints require administrative authentication. These endpoints are accessible via:
- **Headplane password authentication** — Admin users logged in with Headplane credentials
- **API key authentication** — Admin API keys with appropriate permissions (Phase 13c Known Gap fix)

**Base URL**: `/api/v1/headplane/users`

**Authentication**: All endpoints require admin authentication via session token or API key in the `Authorization` header:
```
Authorization: Bearer hp_<session_token>
```
or
```
Authorization: Bearer <api_key>
```

**Content-Type**: All requests and responses use `application/json`.

---

## Endpoints

### GET /api/v1/headplane/users

List all Headplane dashboard users in the system.

#### Authentication

Required. Must be an administrator (authenticated via Headplane password or admin API key).

#### Request

**Method**: `GET`

**Headers**:
```
Authorization: Bearer hp_AbCdEfGhIjKlMnOpQrStUvWxYz0123456789AbCdEfGhIjKlMnOpQrStUvWxYz
```

**Parameters**: None

#### Response

**Success (200 OK)**:

Returns an object containing all Headplane users in its `users` array.

**Schema**:

```yaml
type: object
required: [users]
properties:
  users:
    type: array
    items:
      type: object
      required: [id, username, role, createdAt]
      properties:
        id:
          type: integer
          description: Unique user identifier
          example: 1
        username:
          type: string
          description: User's login username
          example: "admin"
        role:
          type: string
          description: User's role
          example: "admin"
        createdAt:
          type: integer
          format: int64
          description: Unix timestamp of user creation
          example: 1790856000
```

**Example Response**:

```json
{
  "users": [
    {
      "id": 1,
      "username": "admin",
      "role": "admin",
      "createdAt": 1790856000
    },
    {
      "id": 2,
      "username": "operator",
      "role": "user",
      "createdAt": 1790929800
    }
  ]
}
```

#### Error Responses

**401 Unauthorized**:

Session token or API key is invalid, expired, or missing.

```json
{
  "error": "Unauthorized"
}
```

**403 Forbidden**:

The authenticated user does not have admin privileges.

```json
{
  "error": "Forbidden: admin access required"
}
```

**500 Internal Server Error**:

Server error (database failure, etc.).

```json
{
  "error": "Failed to list users"
}
```

#### Example Usage

**With session token:**
```bash
curl -X GET https://headscale.example.com/api/v1/headplane/users \
  -H "Authorization: Bearer hp_AbCdEfGhIjKlMnOpQrStUvWxYz0123456789..."
```

**With API key:**
```bash
curl -X GET https://headscale.example.com/api/v1/headplane/users \
  -H "Authorization: Bearer hs_1234567890abcdefghijklmnopqrstuvwxyz"
```



### POST /api/v1/headplane/users

Create a new Headplane dashboard user.

#### Authentication

Required. Must be an administrator.

#### Request

**Method**: `POST`

**Headers**:
```
Authorization: Bearer hp_AbCdEfGhIjKlMnOpQrStUvWxYz0123456789AbCdEfGhIjKlMnOpQrStUvWxYz
Content-Type: application/json
```

**Body Schema**:

```yaml
type: object
properties:
  username:
    type: string
    minLength: 3
    maxLength: 50
    pattern: ^[a-zA-Z0-9_-]+$
    description: Username for the new user (alphanumeric, underscore, hyphen only)
    example: "newuser"
  password:
    type: string
    minLength: 8
    description: Password for the new user
    example: "secure-password-123"
  role:
    type: string
    enum: [admin, user]
    description: Role for the new user
    default: user
    example: "user"
required:
  - username
  - password
```

**Example Request**:

```json
{
  "username": "operator",
  "password": "secure-password-123",
  "role": "user"
}
```

#### Response

**Success (201 Created)**:

Returns the newly created user (without password).

```json
{
  "id": 3,
  "username": "operator",
  "role": "user",
  "created_at": "2026-10-04T11:11:06.436Z",
  "updated_at": "2026-10-04T11:11:06.436Z"
}
```

#### Error Responses

**400 Bad Request**: Invalid input.

```json
{
  "error": "Username is required"
}
```

**401 Unauthorized**: Session token or API key is invalid.

```json
{
  "error": "Unauthorized"
}
```

**403 Forbidden**: Not an administrator.

```json
{
  "error": "Forbidden: admin access required"
}
```

**409 Conflict**: Username already exists.

```json
{
  "error": "Username already exists"
}
```

#### Example Usage

```bash
curl -X POST https://headscale.example.com/api/v1/headplane/users \
  -H "Authorization: Bearer hp_AbCdEfGhIjKlMnOpQrStUvWxYz0123456789..." \
  -H "Content-Type: application/json" \
  -d '{
    "username": "operator",
    "password": "secure-password-123",
    "role": "user"
  }'
```

---

### PUT /api/v1/headplane/users/:id

Update an existing Headplane dashboard user.

#### Authentication

Required. Must be an administrator.

#### Request

**Method**: `PUT`

**Path Parameters**: `id` (integer, required) - The unique identifier of the user to update

**Body Schema** (all fields optional):

```yaml
type: object
properties:
  username:
    type: string
    minLength: 3
    maxLength: 50
  password:
    type: string
    minLength: 8
  role:
    type: string
    enum: [admin, user]
```

**Example Request**:

```json
{
  "role": "admin"
}
```

#### Response

**Success (200 OK)**: Returns the updated user.

**Error Responses**: 400 Bad Request, 401 Unauthorized, 403 Forbidden, 404 Not Found, 409 Conflict

#### Example Usage

```bash
curl -X PUT https://headscale.example.com/api/v1/headplane/users/2 \
  -H "Authorization: Bearer hp_AbCdEfGhIjKlMnOpQrStUvWxYz0123456789..." \
  -H "Content-Type: application/json" \
  -d '{"role": "admin"}'
```

---

### DELETE /api/v1/headplane/users/:id

Delete a Headplane dashboard user.

#### Authentication

Required. Must be an administrator.

#### Request

**Method**: `DELETE`

**Path Parameters**: `id` (integer, required) - The unique identifier of the user to delete

#### Response

**Success (200 OK)**:

```json
{
  "message": "User deleted successfully"
}
```

**Error Responses**: 401 Unauthorized, 403 Forbidden (including "Cannot delete the last admin user"), 404 Not Found

#### Example Usage

```bash
curl -X DELETE https://headscale.example.com/api/v1/headplane/users/3 \
  -H "Authorization: Bearer hp_AbCdEfGhIjKlMnOpQrStUvWxYz0123456789..."


## Authorization

All user management endpoints require **admin privileges**. This is enforced for both:
- **Headplane password authentication**: Users with `role: admin`
- **API key authentication**: API keys with admin permissions (Phase 13c Known Gap fix)

---

## Security Considerations

### Password Handling

- **Storage**: Passwords are hashed using bcrypt before storage
- **Transmission**: Passwords only transmitted during user creation/update (over HTTPS)
- **Response**: Passwords never returned in API responses

### Transport Security

- **HTTPS Required**: Use HTTPS in production to protect credentials in transit
- **Rate Limiting**: Recommended 10 requests per minute per authenticated session

---

## Phase 13c Known Gap Fix

### What Changed

Prior to the Known Gap fix, user management UI was only accessible when logged in with Headplane password credentials. API key authentication did not provide access to the admin user management interface.

**After the Known Gap fix**:
- Admin API keys can now access all user management endpoints
- The `/admin/users` page shows "Add Headplane User" and full user management UI when accessed via admin API key
- Authorization checks properly validate admin privileges for both session tokens and API keys

### Screenshots

![Admin Users Page - API Key Session](../../assets/screenshots/admin-users-api-key.png)
*Figure 1: Admin users page showing "Add Headplane User" button when authenticated via API key*

![Add Headplane User Dialog](../../assets/screenshots/add-headplane-user-dialog.png)
*Figure 2: Dialog for creating a new Headplane dashboard user*

**Note**: Screenshots show the unified admin UI available to both password-authenticated and API key-authenticated admin sessions. These are placeholder references - actual screenshots should be captured from the implementation.

---

## Related Documentation

- [Settings User Guide](../../usage/settings.md) — End-user documentation for the Settings page
- [Authentication Guide](../../usage/authentication.md) — Password and API key authentication setup
- [Headplane Settings API](./headplane-settings.md) — Settings management endpoints
- [ROADMAP: Phase 13c](../../../ROADMAP.md#phase-13c--multi-user-headplane-implemented) — Feature overview and Known Gap fix

---

## Changelog

- **2026-10-04**: Phase 13c Known Gap documentation — Comprehensive API reference created for user management endpoints, documenting the Known Gap fix that enables admin user management via API key authentication

```

---

---
