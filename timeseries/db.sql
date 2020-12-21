-- https://docs.timescale.com/latest/tutorials/tutorial-hello-timescale
-- https://docs.timescale.com/latest/tutorials/quickstart-go

CREATE DATABASE fcpa;
\c fcpa
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

-- psql -x "postgres://tsdbadmin:@localhost:5432/tsdb?sslmode=require"

-- \dx

DROP TABLE coins;
DROP TABLE tiers;
DROP TABLE exchanges;
DROP TABLE price_queries;

CREATE TABLE IF NOT EXISTS "coins"(
    coin INTEGER,
    code TEXT,
    name TEXT
);
INSERT INTO coins(coin, code, name) VALUES
(1, 'BTC', 'Bitcoin'),
(2, 'ETH', 'Ethereum'),
(3, 'USD', 'United States Dollar'),
(4, 'TRY', 'Turkish Lira'),
(5, 'GBP', 'British Pound'),
(6, 'JPY', 'Japanese Yen'),
(7, 'AED', 'United Arab Emirates dirham');

CREATE TABLE IF NOT EXISTS "tiers"(
    id INTEGER,
    coin INTEGER,
    amount DECIMAL,
    active BOOLEAN DEFAULT FALSE
);
INSERT INTO tiers(id, coin, amount, active) VALUES
(1, 3, 10000, true),
(2, 3, 50000, true),
(3, 3, 100000, true),
(4, 3, 250000, true),
(5, 3, 500000, false),
(6, 4, 80000, true),
(7, 4, 400000, true),
(8, 4, 800000, true),
(9, 4, 1000000, true),
(10, 4, 4000000, false);

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
    exchange_buy INTEGER NOT NULL,
    exchange_sell INTEGER NOT NULL,
    exchange_buy_coin INTEGER NOT NULL,
    exchange_sell_coin INTEGER NOT NULL,
    exchange_buy_timestamp TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    exchange_sell_timestamp TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    exchange_buy_query_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    exchange_sell_query_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
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
    exchange_rate DECIMAL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
);

-- SELECT * FROM tiers INNER JOIN coins ON coins.coin = tiers.coin WHERE coins.code = 'USD';

-- SELECT buy_input, buy_vol, buy_order_count, query_time, created_at FROM price_queries WHERE tier = 3 AND exchange = 2 ORDER BY created_at DESC LIMIT 1;

-- SELECT buy_input, buy_vol, buy_order_count, query_time, created_at FROM price_queries WHERE tier = 3 AND exchange = 1 ORDER BY created_at DESC LIMIT 1;

-- SELECT buy_input, buy_vol, buy_order_count, buy_total_book_length, buy_total_book_crypto, buy_total_book_fiat, sell_total_book_length, sell_total_book_crypto, sell_total_book_fiat, query_time, created_at FROM price_queries WHERE tier = 2 AND exchange = 1 ORDER BY created_at DESC LIMIT 1;

-- select pg_size_pretty(pg_relation_size('price_queries'));

-- SELECT 1116 * 150;
