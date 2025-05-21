-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_histories (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    status VARCHAR(100) NOT NULL,
    shipping_name VARCHAR(100) NOT NULL,
    shipping_address TEXT NOT NULL,
    shipping_phone VARCHAR(25) NOT NULL,
    shipping_number VARCHAR(100) NOT NULL,
    shipping_type VARCHAR(100) NOT NULL,
    total_quantity INT NOT NULL,
    total_weight DOUBLE PRECISION NOT NULL DEFAULT 0,
    voucher_discount BIGINT NOT NULL DEFAULT 0,
    address_id INT NOT NULL,
    cost_name VARCHAR(100),
    cost_service VARCHAR(100),
    voucher_id INT,
    total_product_amount NUMERIC(10,0) NOT NULL DEFAULT 0,
    total_shipping_cost NUMERIC(10,0) NOT NULL DEFAULT 0,
    total_shipping_amount NUMERIC(10,0) NOT NULL DEFAULT 0,
    total_amount NUMERIC(10,0) NOT NULL DEFAULT 0,
    notes TEXT,
    order_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_order_histories_user_id    ON order_histories(user_id);
CREATE INDEX idx_order_histories_address_id ON order_histories(address_id);
CREATE INDEX idx_order_histories_voucher_id ON order_histories(voucher_id);
CREATE INDEX idx_order_histories_active  ON order_histories(deleted_at) WHERE deleted_at IS NULL;

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `order_histories`
CREATE TRIGGER set_updated_at_order_histories
BEFORE UPDATE ON order_histories
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_order_histories_active;
DROP INDEX IF EXISTS idx_order_histories_voucher_id;
DROP INDEX IF EXISTS idx_order_histories_address_id;
DROP INDEX IF EXISTS idx_order_histories_user_id;
DROP TRIGGER IF EXISTS set_updated_at_order_histories ON order_histories;
DROP TABLE IF EXISTS order_histories;
-- +goose StatementEnd
