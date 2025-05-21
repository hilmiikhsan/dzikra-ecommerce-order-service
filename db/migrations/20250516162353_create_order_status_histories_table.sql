-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_status_histories (
    id SERIAL PRIMARY KEY,
    order_id UUID NOT NULL,
    status VARCHAR(100) NOT NULL,
    description TEXT,
    changed_by VARCHAR(100) NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

ALTER TABLE order_status_histories ADD CONSTRAINT fk_order_status_histories_orders FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE RESTRICT ON UPDATE CASCADE;

CREATE INDEX idx_order_status_histories_active  ON order_status_histories(deleted_at) WHERE deleted_at IS NULL;

-- Fungsi untuk memperbarui kolom `updated_at`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW(); -- Simpan timestamp langsung dengan tipe TIMESTAMPTZ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk tabel `order_status_histories`
CREATE TRIGGER set_updated_at_order_status_histories
BEFORE UPDATE ON order_status_histories
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_order_status_histories_active;
DROP TRIGGER IF EXISTS set_updated_at_order_status_histories ON order_status_histories;
DROP TABLE IF EXISTS order_status_histories;
-- +goose StatementEnd
