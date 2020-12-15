-- Get the most expensive exchange rate
-- for a given tier
-- on a given day
SELECT MAX(exchange_rate)
FROM price_queries
WHERE exchange_sell = 2
AND tier = 1
AND DATE(created_at) = '2020-12-15';

-- Get the queries that were below a certain margin from an exchange rate
-- for a given tier
-- on a given day
SELECT DISTINCT
SUM(buy_input) AS total_usd,
SUM(sell_vol) AS total_try,
SUM(sell_input) AS total_crypto
FROM price_queries
WHERE exchange_sell = 2
AND tier = 1
AND exchange_rate <= (7.83 - (7.83 * 0.01))
AND DATE(created_at) = '2020-12-15';

-- https://www.bloomberg.com/quote/USDTRY:CUR?sref=RSaTyALF
