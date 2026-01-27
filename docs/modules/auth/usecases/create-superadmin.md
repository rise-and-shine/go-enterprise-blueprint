# Create Superadmin

Creates the initial superadmin account for system bootstrap. This command should be run once during initial setup.

> **type**: manual_command

> **operation-id**: `create-superadmin`

> **usage**: `./app auth create-superadmin`

## Execute

- Hash the password

- Start UOW

- Create admin

- Create actor permission with superadmin permission

- Apply UOW
