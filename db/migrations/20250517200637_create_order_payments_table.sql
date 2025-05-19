-- +goose Up
-- +goose StatementBegin
CREATE TYPE application_id_enum AS ENUM (
  'POS',
  'E_COMMERCE'
);

CREATE TABLE IF NOT EXISTS order_payments (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL,
    payment_method VARCHAR(100) NOT NULL,
    payment_status VARCHAR(50) NOT NULL,
    payment_type VARCHAR(50),
    transaction_id UUID NOT NULL,
    gross_amount NUMERIC(10,0) NOT NULL,
    transaction_status VARCHAR(50) NOT NULL,
    payment_code VARCHAR(100),
    signature_key VARCHAR(255),
    midtrans_response   jsonb NOT NULL DEFAULT '{}'::jsonb,
    callback_response   jsonb NOT NULL DEFAULT '{}'::jsonb,
    transaction_request jsonb NOT NULL DEFAULT '{}'::jsonb,
    transaction_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP + INTERVAL '24 hours'),
    application_id application_id_enum NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

ALTER TABLE order_payments ADD CONSTRAINT fk_order_payments_orders FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE RESTRICT ON UPDATE CASCADE;

CREATE INDEX idx_order_payments_order_id ON order_payments(order_id);
CREATE UNIQUE INDEX idx_order_payments_transaction_id ON order_payments(transaction_id);
CREATE INDEX idx_order_payments_status ON order_payments(payment_status, transaction_status);
CREATE INDEX idx_order_payments_active  ON order_payments(deleted_at) WHERE deleted_at IS NULL;

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `order_payments`
CREATE TRIGGER set_updated_at_order_payments
BEFORE UPDATE ON order_payments
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS application_id_enum CASCADE;
DROP INDEX IF EXISTS idx_order_payments_active;
DROP INDEX IF EXISTS idx_order_payments_status;
DROP INDEX IF EXISTS idx_order_payments_transaction_id;
DROP INDEX IF EXISTS idx_order_payments_order_id;
DROP TRIGGER IF EXISTS set_updated_at_order_payments ON order_payments;
DROP TABLE IF EXISTS order_payments;
-- +goose StatementEnd
