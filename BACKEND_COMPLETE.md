# Backend Password Authentication - Complete

**Branch:** feature/password-backend
**Commit:** 8e2131bf
**Status:** ✅ Ready for review

## Implementation

✅ Added HeadplaneConfig to types.Config with Password field
✅ Environment variable: HEADSCALE_HEADPLANE_PASSWORD
✅ Endpoint: POST /api/v1/headplane/login
✅ Security: constant-time comparison, rate limiting (5/min), 256-bit tokens
✅ Session management: 24-hour expiry, automatic cleanup
✅ Tests: 5 unit tests, all passing
✅ Documentation: config-example.yaml updated

## Files Changed

- hscontrol/types/config.go (+11)
- hscontrol/app.go (+5)
- config-example.yaml (+13)
- hscontrol/headplane_auth.go (new, 187 lines)
- hscontrol/headplane_auth_test.go (new, 115 lines)

## API

POST /api/v1/headplane/login
{"password": "..."}
→ {"token": "...", "expires_at": 1234567890}

## Next Steps

1. Frontend: implement login form and session management
2. Documentation: update authentication docs
3. Integration: test end-to-end flow
4. Merge: backend → frontend → docs
