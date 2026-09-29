# SwipeLab

SwipeLab is an educational payment processing system built to explore how a card payment travels from a merchant to an issuing bank.

The project is inspired by *Anatomy of the Swipe* and focuses on implementing payment infrastructure concepts rather than integrating with a real card network.

## Architecture

```text
POS
 │
 ▼
Acquirer Processor
 │
 ▼
Card Network
 │
 ▼
Issuer Processor
 │
 ▼
Cardholder Account
```

Each component is implemented as an independent service and communicates over network protocols rather than importing code from other services.

## Services

### Acquirer Processor

**Language:** Go  
**Database:** PostgreSQL

The acquirer receives payment authorization requests from merchants and is responsible for processing and forwarding them through the payment system.

Currently implemented:

- `POST /authorizations`
- Authorization request validation
- Merchant validation
- Authorization transaction creation
- PostgreSQL persistence
- UUID transaction identifiers
- Authorization status tracking

Planned:

- Idempotency
- Concurrent request handling
- Card network communication
- Timeouts and retries
- Authorization responses
- Capture
- Reversals
- Clearing and settlement
- Reconciliation

### Card Network

**Language:** Go

Planned.

The card network will route authorization messages between the acquirer and issuer.

### Issuer Processor

**Language:** Java / Spring Boot

Planned.

The issuer will simulate the cardholder's bank and make authorization decisions based on card and account state.

## Authorization Flow

An authorization represents a request to approve a payment.

For example:

```text
Merchant
   │
   │ 2500 HUF
   ▼
Acquirer
   │
   ▼
Card Network
   │
   ▼
Issuer
   │
   ├── APPROVED
   │
   └── DECLINED
```

An authorization is initially stored with the status:

```text
PENDING
```

Later, the issuer will determine whether it becomes `APPROVED` or `DECLINED`.

## Running the Acquirer

Start PostgreSQL:

```bash
docker compose up -d acquirer-db
```

Then run the acquirer:

```bash
cd acquirer
go run ./cmd/acquirer
```

The service runs on:

```text
http://localhost:8081
```

## Example Authorization

```bash
curl -X POST http://localhost:8081/authorizations \
  -H "Content-Type: application/json" \
  -d '{
    "merchant_id": "MERCHANT_001",
    "amount": 2500,
    "currency": "HUF"
  }'
```

Example response:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "merchant_id": "MERCHANT_001",
  "amount": 2500,
  "currency": "HUF",
  "status": "PENDING",
  "created_at": "2026-09-29T18:00:00Z"
}
```

## Money Representation

Amounts are stored as integers rather than floating-point numbers.

For currencies with two decimal minor units:

```text
€25.99 → 2599
$25.99 → 2599
```

Currency-specific minor-unit rules will be handled explicitly by the payment system.

## Purpose

SwipeLab is not intended to process real card payments.

It is a learning project for exploring concepts such as:

- Payment authorization
- Acquiring and issuing
- Card-network routing
- Idempotency
- Concurrency
- Database transactions
- Timeouts and retries
- Reversals
- Clearing
- Settlement
- Reconciliation

All payment data used by SwipeLab is simulated.