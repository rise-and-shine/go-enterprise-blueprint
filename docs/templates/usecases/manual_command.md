# {Use Case Name}

{One-two sentence summary of what this command does and when an operator would run it.}

> **type**: manual_command

> **operation-id**: `{operation-id}`

> **usage**: `{app-name} {command} [flags] [arguments]`

## Input

Flags:

- `--flag-name`, `-f`: string, required, description

Arguments:

- `arg_name`: string, required, description

## Execute

- Validate input arguments and flags
- Verify operator permissions
- Check preconditions: {what must be true}
- Process {what}
- Produce audit log
- Return result
