CREATE TABLE mailings (
    id BIGINT NOT NULL PRIMARY KEY,
    is_locked BOOLEAN DEFAULT FALSE NOT NULL
);

CREATE TABLE customers (
    id BIGSERIAL PRIMARY KEY,
    email varchar(64) NOT NULL,
    title text,
    content text,
    mailing_id BIGINT REFERENCES mailings(id),
    insertion_time TIMESTAMP
);

CREATE INDEX customers_mailing_id_idx ON customers(mailing_id);
