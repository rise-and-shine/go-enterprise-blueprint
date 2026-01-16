# {Use Case Name}

{One-two sentence summary of what this subscriber does when the event occurs.}

> **type**: event_subscriber

> **operation-id**: `{operation-id}`

> **event**: `{EventName}`

## Event payload

```json
{
  "entity_id": "string", // required
  "action": "string", // required
  "timestamp": "2024-01-15T10:30:00Z", // required, RFC3339
  "metadata": {} // optional, additional context
}
```

## Handle

- Validate event payload
- Check idempotency using {dedup-key}
- Process {what}
- Produce `{OutputEventName}` event (if applicable)
- Log completion

## Idempotency

{Brief description of how duplicate events are handled - event ID tracking, deduplication key, or naturally idempotent operation.}
