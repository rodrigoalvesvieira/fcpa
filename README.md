# Fiat-Crypto Parity Arbitrage Algorithm (FCPA) v0.0.2
> Software to find cryptocurrency arbitrage opportunities within a network of cryptocurrency exchanges

![The arbitrageurs](https://i.imgur.com/Pf5AZFX.jpg)

## Technologies & Services

| Name                   | Version     |
| ---------------------- | -----------:|
| Go                     | 1.15.3      |
| Timescale              | 1.7.4       |
| PostgreSQL             | 11.10       |
| Protocol Buffers       | 3.14.0      |

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

| Name              | Cryptocurrencies            | Status       | Package                |
| ----------------- | ----------------------------|--------------|------------------------|
| Kraken            | BTC                         | Pending      |                        |
| BTCTurk           | BTC                         | Pending      |                        |

## Authors

* Rodrigo Alves Vieira - <<rodrigoalvesvieira@protonmail.com>>

## RESEARCH

* https://data.bitcoinity.org/markets/volume/30d?c=e&t=b

* [Kraken: "How long do digital assets/cryptocurrency deposits take?"](https://support.kraken.com/hc/en-us/articles/203325283-How-long-do-digital-assets-cryptocurrency-deposits-take-)

## Copyright

© 2018-2020 Rodrigo Alves Vieira. All Rights Reserved.
