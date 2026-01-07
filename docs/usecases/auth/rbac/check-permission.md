# Check Permission

**UseCase type:** user_action

**Module:** auth

**Operation ID:** auth.check-permission

## Summary

Checks if an actor has a specific permission. Used by authorization middleware and services to enforce access control. This is a high-frequency, performance-critical operation.

## Actor

- **Primary:** System, Service
- **Authorization:** `auth:permission:check` or checking own permissions

## Input

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| actor_type | string | yes | enum: user, group, service_acc |
| actor_id | string | yes | must exist |
| permission | string | yes | permission to check |
| resource_id | string | no | specific resource ID for resource-level permissions |

## Output

| Field | Type | Description |
|-------|------|-------------|
| allowed | bool | True if actor has the permission |
| source | string | How permission was granted (role/direct/inherited) |
| role_name | string | Role that granted permission (if applicable) |
| cached | bool | Whether result came from cache |

## Flow

### Validate

1. Validate input format
2. Verify actor_type is valid
3. Validate permission string format

### Execute

1. Check authorization cache for this actor+permission combination
2. If cache hit, return cached result
3. If cache miss:
   - Check direct permissions (ActorPermission)
   - Check role-based permissions (ActorRole -> RolePermission)
   - For groups: check inherited permissions from parent groups
   - For resource_id: check resource-specific grants
4. Cache the result with appropriate TTL
5. Return permission decision

### Side Effects

- [ ] Update cache with result
- [ ] Optionally log permission checks for audit (configurable)

## Error Scenarios

| Error | Code | When |
|-------|------|------|
| ErrInvalidInput | 400 | Invalid actor_type or permission format |
| ErrActorNotFound | 404 | Actor doesn't exist |

## Business Rules

1. Permission is granted if ANY of these conditions are true:
   - Actor has direct permission grant
   - Actor has a role with this permission
   - Actor belongs to a group with this permission
   - Actor has wildcard permission (e.g., `auth:user:*`)
2. Superuser role has all permissions implicitly
3. Inactive users always return denied
4. Results should be cached aggressively
5. Cache TTL should balance security vs performance

## Notes

- This is performance-critical - optimize for speed
- Use Redis or in-memory cache with short TTL
- Consider bloom filters for negative caching
- Implement permission hierarchy (wildcards)
- Log denied attempts for security monitoring
- Consider batch checking for multiple permissions

### Performance Considerations

- Target response time: < 5ms (cached), < 50ms (uncached)
- Cache TTL: 5-15 minutes recommended
- Implement cache warming for frequently accessed actors
- Use connection pooling for database queries

### Permission Format

```
module:resource:action

Examples:
- auth:user:create
- auth:user:*           (all user actions)
- auth:*:read           (read any auth resource)
- *:*:*                 (superuser - all permissions)
```
