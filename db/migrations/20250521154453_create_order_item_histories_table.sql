-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_item_histories (
    id SERIAL PRIMARY KEY,
    order_id UUID NOT NULL,
    product_id INT NOT NULL,
    product_name VARCHAR(100) NOT NULL,
    product_variant VARCHAR(100),
    product_discount BIGINT DEFAULT 0,
    quantity INTEGER NOT NULL DEFAULT 0,
    product_weight DOUBLE PRECISION NOT NULL DEFAULT 0,
    product_price NUMERIC(10,0) NOT NULL DEFAULT 0,
    product_discount_price NUMERIC(10,0) DEFAULT 0,
    total_amount NUMERIC(10,0) NOT NULL DEFAULT 0,
    product_variant_id INT,
    product_capital_price NUMERIC(10,0) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

ALTER TABLE order_item_histories ADD CONSTRAINT fk_order_item_histories_orders FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE RESTRICT ON UPDATE CASCADE;

CREATE INDEX idx_order_item_histories_active  ON order_item_histories(deleted_at) WHERE deleted_at IS NULL;

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `order_item_histories`
CREATE TRIGGER set_updated_at_order_item_histories
BEFORE UPDATE ON order_item_histories
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_order_item_histories_active;
DROP TRIGGER IF EXISTS set_updated_at_order_item_histories ON order_item_histories;
DROP TABLE IF EXISTS order_item_histories;
-- +goose StatementEnd
