-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_items (
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
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

ALTER TABLE order_items ADD CONSTRAINT fk_order_items_orders FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE RESTRICT ON UPDATE CASCADE;

CREATE INDEX idx_order_items_active  ON order_items(deleted_at) WHERE deleted_at IS NULL;

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `order_items`
CREATE TRIGGER set_updated_at_order_items
BEFORE UPDATE ON order_items
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_order_items_active;
DROP TRIGGER IF EXISTS set_updated_at_order_items ON order_items;
DROP TABLE IF EXISTS order_items;
-- +goose StatementEnd
