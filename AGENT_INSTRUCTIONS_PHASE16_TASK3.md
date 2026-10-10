# Phase 16 Task 3: Secure API Keys Retrieval (SECURITY CRITICAL)

**Branch:** `feature/phase16-api-proxy`  
**Worktree:** `/home/denny/Project/headscale-project/headscale-phase16-api-proxy`  
**Base Commit:** `813c46df` (dev branch)

---

## Objective

Proxy all Headscale API key operations through Headplane backend to eliminate direct browser-to-Headscale requests.

---

## Problem Statement (SECURITY CRITICAL)

In the `/admin/admin` page, the browser attempts to retrieve API keys directly from the `headscale` service. This exposes internal service endpoints to the client and violates the principle of backend-for-frontend (BFF) architecture.

**Current (insecure) flow:**
```
Browser → headscale:8080/api/v1/apikeys (direct)
```

**Expected (secure) flow:**
```
Browser → headplane:3000/api/admin/apikeys → headscale:8080/api/v1/apikeys
```

**Security implications:**
- Exposes internal headscale API endpoints to clients
- May bypass headplane authentication/authorization checks
- Reveals internal network topology
- Potential CORS issues and security policy violations

---

## Requirements

### Backend (Headplane)
1. Create proxy endpoints for all API key operations
2. Verify requesting user is admin (session token auth)
3. Add rate limiting to prevent abuse
4. Log all operations for audit trail
5. Proper error handling and HTTP status codes

### Frontend (Headplane)
1. Update all API key operations to use Headplane proxy endpoints
2. Remove all direct calls to `headscale:8080`
3. Update API client configuration
4. Ensure error handling remains functional

---

## Implementation Steps

### 1. Create Headplane Proxy Endpoints

Create new API route in Headplane:
```typescript
// headplane/app/server/api/admin/apikeys.ts

import { headscaleClient } from '~/server/headscale/client';
import { requireAuth, requireAdmin } from '~/server/auth/middleware';

export async function getApiKeys(request: Request) {
  const session = await requireAuth(request);
  await requireAdmin(session);
  
  const response = await headscaleClient.get('/api/v1/apikeys');
  return Response.json(response.data);
}

export async function createApiKey(request: Request) {
  const session = await requireAuth(request);
  await requireAdmin(session);
  
  const body = await request.json();
  const response = await headscaleClient.post('/api/v1/apikeys', body);
  
  // Audit log
  await logAudit({
    action: 'apikey.create',
    actor: session.username,
    details: { keyId: response.data.id }
  });
  
  return Response.json(response.data);
}
```

### 2. Update Frontend API Client

Locate and update the API client:
```bash
cd headplane
find app -name "*.ts" -o -name "*.tsx" | xargs grep -l "headscale:8080/api/v1/apikeys"
```

Update to use Headplane proxy:
```typescript
// Before:
// const response = await fetch('http://headscale:8080/api/v1/apikeys');

// After:
const response = await fetch('/api/admin/apikeys');
```

### 3. Test in Browser DevTools

Verify NO direct requests to headscale:
1. Open browser DevTools → Network tab
2. Navigate to `/admin/admin`
3. Verify all requests go to `headplane:3000/*`
4. Verify NO requests to `headscale:8080/*`

---

## Files Expected to Change

- `headplane/app/server/api/admin/apikeys.ts` (new)
- `headplane/app/routes/admin/admin/page.tsx` (UI)
- `headplane/app/lib/api-client.ts` or similar (API client)

---

## Acceptance Criteria

- [x] Browser NEVER makes direct requests to headscale:8080
- [x] All API key operations work through Headplane proxy
- [x] Proper authentication/authorization on proxy endpoints
- [x] Rate limiting prevents abuse
- [x] Audit logging for all operations

---

## Commit Message

```
headplane: proxy API key operations through backend (SECURITY)

Replace direct browser-to-Headscale API key requests with
secure proxy through Headplane backend.

Fixes Phase 16 Issue 3 (SECURITY CRITICAL).
Refs: ROADMAP.md Phase 16
```
