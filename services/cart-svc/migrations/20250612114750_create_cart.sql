-- +goose Up
CREATE TYPE cart_status AS ENUM ('OPEN', 'CHECKOUT', 'ABANDONED');

CREATE TABLE cart (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID         NULL,
    status     cart_status  NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE cart_item (
    cart_id    UUID         NOT NULL REFERENCES cart(id) ON DELETE CASCADE,
    product_id UUID         NOT NULL,
    price      NUMERIC(10,2) NOT NULL,
    quantity   INT          NOT NULL CHECK (quantity >= 0),
    PRIMARY KEY (cart_id, product_id)
);

CREATE INDEX cart_item_product_idx ON cart_item (product_id, cart_id);

-- +goose Down
DROP INDEX IF EXISTS cart_item_product_idx;
DROP TABLE IF EXISTS cart_item;
DROP TABLE IF EXISTS cart;
DROP TYPE IF EXISTS cart_status;
