CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE TABLE category
(
    id   UUID DEFAULT uuid_generate_v4() NOT NULL,
    name VARCHAR(100)                    NOT NULL,
    slug VARCHAR(100)                    NOT NULL,
    CONSTRAINT category_pkey PRIMARY KEY (id),
    CONSTRAINT category_slug_key UNIQUE (slug)
);

CREATE TABLE product
(
    id          UUID                        DEFAULT uuid_generate_v4() NOT NULL,
    sku         VARCHAR(64)                                            NOT NULL,
    name        VARCHAR(255)                                           NOT NULL,
    description TEXT,
    created_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at  TIMESTAMP WITHOUT TIME ZONE,
    CONSTRAINT product_pkey PRIMARY KEY (id),
    CONSTRAINT product_sku_key UNIQUE (sku)
);

CREATE TABLE product_category
(
    product_id  UUID NOT NULL,
    category_id UUID NOT NULL,
    CONSTRAINT product_category_pkey PRIMARY KEY (product_id, category_id),
    CONSTRAINT product_category_product_id_fkey FOREIGN KEY (product_id) REFERENCES product (id) ON DELETE CASCADE,
    CONSTRAINT product_category_category_id_fkey FOREIGN KEY (category_id) REFERENCES category (id) ON DELETE CASCADE
);

CREATE TABLE product_image
(
    id         UUID                        DEFAULT uuid_generate_v4() NOT NULL,
    product_id UUID                                                   NOT NULL,
    s3_key     VARCHAR(512)                                           NOT NULL,
    url        VARCHAR(512)                                           NOT NULL,
    is_primary BOOLEAN                     DEFAULT FALSE              NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    CONSTRAINT product_image_pkey PRIMARY KEY (id),
    CONSTRAINT product_image_product_id_fkey FOREIGN KEY (product_id) REFERENCES product (id) ON DELETE CASCADE
);

CREATE INDEX idx_product_image_primary ON product_image (product_id) WHERE is_primary = true;

CREATE TABLE product_inventory
(
    product_id UUID PRIMARY KEY,
    quantity   INTEGER NOT NULL CHECK (quantity >= 0),
    CONSTRAINT product_inventory_product_id_fkey FOREIGN KEY (product_id) REFERENCES product (id) ON DELETE CASCADE
);

CREATE TABLE product_price
(
    id         UUID                        DEFAULT uuid_generate_v4() NOT NULL,
    product_id UUID                                                   NOT NULL,
    amount     NUMERIC(10, 2)                                         NOT NULL,
    currency   CHAR(3)                     DEFAULT 'KZT'              NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    CONSTRAINT product_price_pkey PRIMARY KEY (id),
    CONSTRAINT product_price_product_id_fkey FOREIGN KEY (product_id) REFERENCES product (id) ON DELETE CASCADE
);

CREATE INDEX idx_product_name_trgm ON product USING gin (name gin_trgm_ops);
