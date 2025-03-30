# Making REST Call POST Idempotent

## Run Once => Item Created

```bash
curl -X POST http://localhost:8080/order \
    -H "Content-Type: application/json" \
    -d '{
        "idempotency_key": "key-12345",
        "item_name": "item A",
        "amount": 100
    }'
```

## Run Again Without Changing The Key => Item Already Created

```bash
curl -X POST http://localhost:8080/order \
    -H "Content-Type: application/json" \
    -d '{
        "idempotency_key": "key-12345",
        "item_name": "item B",
        "amount": 200
    }'
