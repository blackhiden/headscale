# Phase 13b Backend Implementation - Complete

## Summary

Successfully implemented backend settings storage and REST endpoints for Headplane settings management.

## Files Created

1. **hscontrol/headplane_settings.go** (~12KB)
   - `HeadplaneSettings` model with GORM annotations
   - AES-256-GCM encryption/decryption functions with PBKDF2 key derivation
   - CRUD functions: `GetSettings()`, `UpdateSettings()`, `InitHeadplaneSettings()`
   - HTTP handlers: `HandleGetSettings()`, `HandleUpdateSettings()`, `HandleChangePassword()`
   - Table initialization with AutoMigrate

2. **hscontrol/headplane_settings_test.go** (~13KB)
   - 13 comprehensive unit tests covering:
     - Encryption/decryption roundtrip
     - Empty token validation
     - Wrong token decryption failure
     - Settings CRUD operations
     - HTTP endpoint authorization
     - GET/POST settings endpoints
     - Partial settings updates
     - Password change (success, invalid current, empty new)

## Files Modified

1. **hscontrol/app.go**
   - Added `InitHeadplaneSettings()` call in `NewHeadscale()` (line 139)
   - Registered 3 new endpoints in `createRouter()` (lines 524-526):
     - `GET /api/v1/headplane/settings`
     - `POST /api/v1/headplane/settings`
     - `POST /api/v1/headplane/change-password`

2. **go.mod**
   - Added `gorm.io/driver/sqlite` dependency for test database

## Implementation Details

### Database Schema

```sql
CREATE TABLE headplane_settings (
  id INTEGER PRIMARY KEY CHECK (id = 1),  -- Single row constraint
  api_key_encrypted TEXT,                 -- AES-256-GCM encrypted
  api_key_nonce TEXT,                     -- Encryption nonce (base64)
  api_key_salt TEXT,                      -- PBKDF2 salt (base64)
  theme TEXT DEFAULT 'light',             -- 'light' or 'dark'
  profile_name TEXT,                      -- Optional display name
  updated_at TIMESTAMP
);
```

### Encryption Details

- **Algorithm:** AES-256-GCM (authenticated encryption)
- **Key Derivation:** PBKDF2 with SHA256 (100,000 iterations, 32-byte key)
- **Salt:** 32 random bytes per encryption
- **Nonce:** 12 bytes (GCM standard)
- **Encoding:** Base64 for storage

### API Endpoints

#### GET /api/v1/headplane/settings
- **Auth:** Valid session token (header or cookie)
- **Response:** `{apiKey: string, theme: string, profileName: string}`
- **Behavior:** Decrypts API key before returning; returns defaults if no settings exist

#### POST /api/v1/headplane/settings
- **Auth:** Valid session token (header or cookie)
- **Body:** `{apiKey?: string, theme?: string, profileName?: string}`
- **Behavior:** Partial updates supported; encrypts API key before storing

#### POST /api/v1/headplane/change-password
- **Auth:** Valid session token (header or cookie)
- **Body:** `{currentPassword: string, newPassword: string}`
- **Behavior:** Validates current password (constant-time), updates config in-memory

### Test Coverage

- ✅ Encryption roundtrip with correct token
- ✅ Encryption with empty token (error)
- ✅ Decryption with wrong token (error)
- ✅ Get non-existent settings (error)
- ✅ Create settings (first time)
- ✅ Update settings (single row maintained)
- ✅ GET endpoint without auth (401)
- ✅ GET endpoint with auth, no settings (returns defaults)
- ✅ GET endpoint with auth, stored settings (decrypts API key)
- ✅ POST endpoint creates and encrypts settings
- ✅ POST endpoint partial update (preserves unchanged fields)
- ✅ Change password with correct current password
- ✅ Change password with incorrect current password (401)
- ✅ Change password with empty new password (400)

## Security Features

1. **API Key Encryption at Rest:** AES-256-GCM with session-derived key
2. **Constant-Time Password Comparison:** Prevents timing attacks
3. **Session-Based Authentication:** All endpoints require valid session
4. **Single-Row Constraint:** Enforces single-user scenario (id=1)

## Next Steps

1. Frontend implementation will consume these endpoints
2. Integration testing after both backend and frontend complete
3. Documentation updates for API reference

## Notes

- Password change updates in-memory config only (not persisted to file)
- Session token used as encryption key base (via PBKDF2)
- Supports both Authorization header and cookie-based auth
- Compatible with Phase 13a (existing login functionality)
