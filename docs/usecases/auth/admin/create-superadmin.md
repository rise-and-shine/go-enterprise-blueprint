# Create Superadmin

Creates the initial superadmin account for system bootstrap. This command should be run once during initial setup.

> **type**: manual_command

> **operation-id**: `create-superadmin`

> **usage**: `app create-superadmin --username <username> --password <password>`

## Execute

- Hash the password using bcrypt

- Create admin record with `is_superadmin=true`
