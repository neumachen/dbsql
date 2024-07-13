-- Migration goes here.
CREATE TABLE testing_datatypes (
    testing_datatype_id SERIAL NOT NULL,
    testing_datatype_uuid UUID NOT NULL,
    word VARCHAR NOT NULL,
    paragraph TEXT NOT NULL,
    metadata JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE customers (
    customer_id SERIAL PRIMARY KEY NOT NULL,
    last_name VARCHAR NOT NULL,
    first_name VARCHAR NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE
    DEFAULT (now() AT TIME ZONE 'UTC') NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE
    DEFAULT (now() AT TIME ZONE 'UTC') NOT NULL
);

CREATE TABLE email_addresses (
    email_address_id SERIAL PRIMARY KEY NOT NULL,
    customer_id INTEGER NOT NULL,
    email_address VARCHAR NOT NULL,
    CONSTRAINT customer_fk FOREIGN KEY (customer_id)
    REFERENCES customers (customer_id) ON DELETE CASCADE,
    CONSTRAINT uniq_email_address UNIQUE (email_address),
    CONSTRAINT chk_email_address_case CHECK (
        email_address = lower(email_address)
    )
);

CREATE UNIQUE INDEX idx_uniq_email_address
ON email_addresses
USING btree (lower(email_address));

CREATE INDEX idx_email_address
ON email_addresses
USING btree (customer_id);

CREATE TABLE addresses (
    address_id SERIAL NOT NULL,
    street_number VARCHAR NOT NULL,
    route VARCHAR NOT NULL,
    unit_number VARCHAR,
    locality VARCHAR NOT NULL,
    administrative_area_level_1 VARCHAR NOT NULL,
    country VARCHAR NOT NULL,
    postal_code VARCHAR NOT NULL,
    latitude NUMERIC,
    longitude NUMERIC,
    geodata JSONB,
    geom GEOMETRY NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE
    DEFAULT (now() AT TIME ZONE 'utc') NOT NULL,
    CONSTRAINT addresses_pk PRIMARY KEY (address_id)
);

CREATE TABLE customer_addresses (
    customer_id INTEGER NOT NULL,
    address_id INTEGER NOT NULL,
    CONSTRAINT customer_address_id PRIMARY KEY (customer_id, address_id),
    CONSTRAINT customer_fk FOREIGN KEY (customer_id)
    REFERENCES customers (customer_id) ON DELETE CASCADE,
    CONSTRAINT address_fk FOREIGN KEY (address_id)
    REFERENCES addresses (address_id)
    ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_uniq_customer_address
ON customer_addresses
USING btree (customer_id, address_id);
