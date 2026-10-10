# Fix: Frontend TypeError on /admin/admin/users.data

## Problem
The frontend was throwing a `TypeError: Cannot read properties of undefined (reading 'find')` when accessing the users management page at `/admin/admin/users`.

### Root Cause
The backend endpoint `/api/v1/headplane/users` was returning a bare JSON array:
```json
[
  {"id": 1, "username": "user1", "role": "admin", "createdAt": 1234567890}
]
```

But the frontend at `headplane/app/routes/admin/users/route.tsx:73` expected an object with a `users` property:
```typescript
const users = (data.users || []) as HeadplaneUserData[];
```

## Solution
Modified `hscontrol/headplane_users.go` (lines 222-239) to wrap the user array in an object:

### Before
```go
response := make([]HeadplaneUserResponse, len(users))
// ... populate response ...
json.NewEncoder(w).Encode(response)  // Returns [...]
```

### After
```go
userList := make([]HeadplaneUserResponse, len(users))
// ... populate userList ...
response := map[string]interface{}{
    "users": userList,
}
json.NewEncoder(w).Encode(response)  // Returns {"users": [...]}
```

## Files Changed
1. **`hscontrol/headplane_users.go`** - Modified `HandleListUsers()` to wrap response in object
2. **`hscontrol/headplane_users_response_test.go`** - Added tests to verify response format

## Testing
- Created `TestHandleListUsers_ResponseFormat` to verify the response structure
- Created `TestHandleListUsers_EmptyResponseFormat` to verify empty list handling
- Both tests pass successfully
- Build succeeds without errors

## Expected Result
The frontend will now receive:
```json
{
  "users": [
    {"id": 1, "username": "user1", "role": "admin", "createdAt": 1234567890}
  ]
}
```

And `data.users` will be defined, eliminating the TypeError.
