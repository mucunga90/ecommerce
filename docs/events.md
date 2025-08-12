# Event Handling in MyService

## Supported Events

- `UserCreated`: Triggered when a new user registers.
- `OrderPlaced`: Triggered when an order is placed.

## Event Payload

Events should be sent as JSON POST requests to `/events` endpoint:

```json
{
  "type": "UserCreated",
  "payload": {
    "id": 123,
    "name": "John Doe"
  }
}
