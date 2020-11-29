-- https://docs.timescale.com/latest/tutorials/tutorial-hello-timescale
-- https://docs.timescale.com/latest/tutorials/quickstart-go

CREATE DATABASE fcpa;
\c fcpa
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

-- psql -x "postgres://tsdbadmin:@localhost:5432/tsdb?sslmode=require"

-- \dx

CREATE TABLE IF NOT EXISTS "coins"(
    coin INTEGER,
    code TEXT,
    name TEXT
);
INSERT INTO coins(coin, code, name) VALUES
(1, 'BTC', 'Bitcoin'),
(2, 'ETH', 'Ethereum'),
(3, 'USD', 'United States Dollar'),
(4, 'TRY', 'Turkish Lira');

CREATE TABLE IF NOT EXISTS "exchanges"(
    exchange INTEGER,
    name TEXT
);
INSERT INTO exchanges(exchange, name) VALUES
(1, 'Kraken'),
(2, 'BTCTurk');

CREATE TABLE "price_queries"(
    coin INTEGER,
    exchange_a INTEGER,
    exchange_b INTEGER,
    exchange_a_coin INTEGER,
    exchange_b_coin INTEGER,
    exchange_a_query_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    exchange_b_query_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    exchange_a_tier_1_input DECIMAL,
    exchange_a_tier_1_order_count INTEGER,
    exchange_a_tier_1_vol DECIMAL,
    exchange_b_tier_1_input DECIMAL,
    exchange_b_tier_1_order_count INTEGER,
    exchange_b_tier_1_vol DECIMAL,
    exchange_a_tier_2_input DECIMAL,
    exchange_a_tier_2_order_count INTEGER,
    exchange_a_tier_2_vol DECIMAL,
    exchange_b_tier_2_input DECIMAL,
    exchange_b_tier_2_order_count INTEGER,
    exchange_b_tier_2_vol DECIMAL,
    exchange_a_tier_3_input DECIMAL,
    exchange_a_tier_3_order_count INTEGER,
    exchange_a_tier_3_vol DECIMAL,
    exchange_b_tier_3_input DECIMAL,
    exchange_b_tier_3_order_count INTEGER,
    exchange_b_tier_3_vol DECIMAL,
    kraken_total_book_length NUMERIC,
    kraken_total_book_crypto DECIMAL,
    kraken_total_book_fiat DECIMAL,
    btcturk_total_book_length NUMERIC,
    btcturk_total_book_crypto DECIMAL,
    btcturk_total_book_fiat DECIMAL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
