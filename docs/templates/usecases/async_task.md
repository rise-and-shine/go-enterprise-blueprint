# {Use Case Name}

{One-two sentence summary of what this job does and why it runs on a schedule.}

> **type**: async_task

> **operation-id**: `{operation-id}`

## Task payload

```json
{
  "entity_id": "string", // required
  "action": "string", // required
  "timestamp": "2024-01-15T10:30:00Z" // required, RFC3339
}
```

## Handle

- Validate task payload
- Process {what}

## Idempotency

{Brief description of how the job handles reruns - watermark/timestamp tracking, idempotent operations, overlap handling.}
