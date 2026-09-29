ALTER TABLE authorizations
ADD COLUMN idempotency_key VARCHAR(255);

CREATE UNIQUE INDEX authorizations_merchant_id_idempotency_key_unique
ON authorizations(merchant_id, idempotency_key);