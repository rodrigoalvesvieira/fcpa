# OPERATOR DESK v0.0.2
Software to find cryptocurrency arbitrage opportunities within a network of cryptocurrency exchanges

## Technologies & Services

| Name                   | Version     |
| ---------------------- | -----------:|
| Go                     | 1.10.2      |
| KDB+                   | 3.6         |
| PostgreSQL             | 9.6.1.0     |
| Twilio                 |             |

## Strategy

Take two cryptocurrency exchanges A and B. A is selling coin `X` for _$5570_ and B is selling the same coin for _$5585_. That constitutes a $15 spread.

The idea is as follows:

1. Buy `X` on exchange A with initial operational capital **β**
2. Sell `Y` on exchange B for **β +** **_delta_**
3. **_delta_** is considered profit
4. Wait for opportunity to repeat operation from `Y` to `X`

## Ultimate Goal

1. The goal is to profit $15 per trade.
2. The more money to stake the more likely it is that a $15 profit can be made and the higher the frequency it may occur.
3. The more Cryptocurrencies and exchanges that are supported, the higher the frequency of possible profit
4. Operational Capital is always kept in fiat currency (USD or EUR)

| Profit Target   | Frequency _per day_ | Total profit |
| --------------- | -------------------:|-------------:|
| $15             | 10                  | $150         |
| $15             | 20                  | $300         |
| $15             | 50                  | $750         |
| $15             | 100                 | $1500        | 

## Up & Running

### Installing Dependencies

```shell
$ dep ensure
```

### How it Works

```shell
$ go run operatordesk.go
```

#### Integrated Exchanges

| Name                | Cryptocurrencies            | Status       | Package                |  
| ------------------- | ----------------------------|--------------|------------------------|
| Kraken              | BTC, ETH, LTC               | Pending      |                        |
| Bitstamp            | BTC, ETH                    | Pending      |                        |

## Authors

* Rodrigo Alves Vieira - <<rodrigoalvesvieira@protonmail.com>>

## Copyright

© 2018 Rodrigo Alves Vieira. All Rights Reserved.

[bitcointrade pkg]: https://github.com/rodrigoalvesvieira/operatordesk-bitcointrade
[negociecoins pkg]: https://github.com/rodrigoalvesvieira/operatordesk-negociecoins
