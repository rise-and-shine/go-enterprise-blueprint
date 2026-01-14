# Create Superadmin

Creates the initial superadmin account for system bootstrap. This command should be run once during initial setup.

> **type**: manual_command

> **operation-id**: `create-superadmin`

> **usage**: `app create-superadmin --username <username> --password <password>`

## Execute

- Ask username and password interactively
- Validate in a loop until input is valid
- Hash the password using bcrypt
- Create admin record with `is_superuser=true`
