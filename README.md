# Fiat-Crypto Parity Arbitrage Algorithm (FCPA) v0.0.3
> Software to find cryptocurrency arbitrage opportunities within a network of cryptocurrency exchanges

![The arbitrageurs](https://i.imgur.com/Pf5AZFX.jpg)

## Technologies

| Name                   | Version     |
| ---------------------- | -----------:|
| Go                     | 1.15.6      |
| Timescale              | 1.7.4       |
| PostgreSQL             | 11.10       |
| Protocol Buffers       | 3.14.0      |

## Up & Running

```

### How it Works

```shell
$ go run fcpa
```

#### Integrated Exchanges

| Name              | Cryptocurrencies            | Status       | Package                |
| ----------------- | ----------------------------|--------------|------------------------|
| Kraken            | BTC                         | Pending      |                        |
| BTCTurk           | BTC                         | Pending      |                        |

## More

Using [VisiData]:

```sql

\copy (SELECT * FROM price_queries ORDER BY created_at DESC) to '/Users/rodrigo/src/src/github.com/rodrigoalvesvieira/fcpa/extract.csv' with csv header;
```

```sh
$ vd extract.csv
```

## Authors

* Rodrigo Alves Vieira - <<rodrigoalvesvieira@protonmail.com>>

## Research

* https://data.bitcoinity.org/markets/volume/30d?c=e&t=b
* http://www.erights.org/elib/capability/ode/index.html
* https://medium.com/@FidelityDigitalAssets/the-omnibus-model-for-custody-96b69710f92d

## Copyright

© 2018-2020 Rodrigo Alves Vieira. All Rights Reserved.

[VisiData]: https://www.visidata.org/
