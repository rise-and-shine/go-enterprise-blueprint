# Create Superuser

**UseCase type:** manual_command

**Module:** auth

**Operation ID:** auth.create-superuser

## Summary

Creates a superuser account with full system privileges. This command is run by operators during initial system setup or when a new administrator needs to be provisioned outside of normal workflows.

## Usage

```bash
app auth create-superuser --username <username> --email <email>
```

### Flags

| Flag | Short | Type | Required | Description |
|------|-------|------|----------|-------------|
| --username | -u | string | yes | Username for the superuser account |
| --email | -e | string | yes | Email address for password setup |
| --pin | -p | string | no | Person PIN (creates Person if provided) |
| --force | -f | bool | no | Skip confirmation prompt |

### Arguments

No positional arguments required.

## Examples

```bash
# Example 1: Basic superuser creation
app auth create-superuser --username admin --email admin@example.com

# Example 2: With Person PIN
app auth create-superuser --username admin --email admin@example.com --pin 12345678901234

# Example 3: Non-interactive mode
app auth create-superuser --username admin --email admin@example.com --force
```

## Flow

### Validate

1. Validate username format (3-50 chars, alphanumeric with underscore)
2. Validate email format
3. Check username is not already taken
4. If PIN provided, validate format (14 digits)

### Execute

1. Generate a secure random temporary password
2. If PIN provided, create or find Person record
3. Create User record with is_active=true
4. Assign "superuser" role to the new user
5. Grant all system permissions
6. Send password setup email to provided address
7. Output success message with username

### Side Effects

- [ ] Update database: Create User, Person (if PIN provided), ActorRole records
- [ ] Emit event: `SuperUserCreated`
- [ ] Create audit log
- [ ] Send email: Password setup link

## Preconditions

- [ ] Database connection is available
- [ ] Email service is configured (for password setup email)
- [ ] "superuser" role exists in the system

## Output

### Success

```
Superuser created successfully!
Username: admin
Email: admin@example.com
Status: Active

A password setup link has been sent to admin@example.com
The link expires in 24 hours.
```

### Failure

```
Error: Username 'admin' already exists
Use a different username or reset the existing account password.
```

## Rollback

- **Reversible:** yes
- **Rollback procedure:** Delete the created User record and associated ActorRole entries. If Person was created by this command, it can be kept or deleted based on policy.

## Permissions

- **Required access:** CLI access to the server, database connectivity
- **Audit logging:** yes

## Safety

- [ ] **Dry-run support:** yes (--dry-run flag shows what would be created)
- [ ] **Confirmation prompt:** yes (unless --force is used)
- [ ] **Idempotent:** no (will fail if username exists)

## Notes

- This command should only be used for initial setup or emergency access recovery
- The temporary password expires after 24 hours
- All superuser actions are logged for audit purposes
- Consider implementing MFA requirement for superuser accounts

### When to Use

- Initial system deployment - creating the first admin account
- Emergency access recovery when all admin accounts are locked
- Provisioning new system administrators

### When NOT to Use

- Regular user creation - use the create-user API instead
- Promoting existing users - use assign-role-to-actor instead
- Automated provisioning - use the API with proper service account
