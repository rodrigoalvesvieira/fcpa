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

CREATE TABLE IF NOT EXISTS "tiers"(
    id INTEGER,
    coin INTEGER,
    amount DECIMAL
);
INSERT INTO tiers(id, coin, amount) VALUES
(1, 3, 10000),
(2, 3, 100000),
(3, 3, 500000);

CREATE TABLE IF NOT EXISTS "exchanges"(
    exchange INTEGER,
    name TEXT
);
INSERT INTO exchanges(exchange, name) VALUES
(1, 'Kraken'),
(2, 'BTCTurk');

CREATE TABLE IF NOT EXISTS "price_queries"(
    coin INTEGER NOT NULL,
    tier INTEGER NOT NULL,
    crypto_input DECIMAL,
    exchange INTEGER NOT NULL,
    exchange_coin INTEGER NOT NULL,
    query_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    buy_input DECIMAL,
    buy_order_count INTEGER,
    buy_vol DECIMAL,
    sell_input DECIMAL,
    sell_order_count INTEGER,
    sell_vol DECIMAL,
    buy_total_book_length NUMERIC,
    buy_total_book_crypto DECIMAL,
    buy_total_book_fiat DECIMAL,
    sell_total_book_length NUMERIC,
    sell_total_book_crypto DECIMAL,
    sell_total_book_fiat DECIMAL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

-- SELECT * FROM tiers INNER JOIN coins ON coins.coin = tiers.coin WHERE coins.code = 'USD';

-- SELECT buy_input, buy_vol, buy_order_count, query_time, created_at FROM price_queries WHERE tier = 3 AND exchange = 2 ORDER BY created_at DESC LIMIT 1;

-- SELECT buy_input, buy_vol, buy_order_count, query_time, created_at FROM price_queries WHERE tier = 3 AND exchange = 1 ORDER BY created_at DESC LIMIT 1;
