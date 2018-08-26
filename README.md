# OPERATOR DESK
Software to find cryptocurrency arbitrage opportunities within a network of Brazilian exchanges

## Strategy

Take two cryptocurrency exchanges A and B. A is selling coin `X` for R$20k and B is selling the same coin for R$20.3k. The idea is as follows:

1. Buy `X` on exchange A
2. Sell `Y` on exchange B
3. Buy R$300 worth of `BTC` and keep it!
4. Rinse and repeat!

## Ultimate Goal

2 BTC per month!

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
| Foxbit              | BTC                         | Pending      |                        |
| MercadoBitcoin      | BTC, ETH                    | Pending      |                        |
| BitcoinTrade        | BTC                         | Pending      | [bitcointrade pkg]     |
| OmniTrade           | BTC, ETH, LTC, BCH, DASH    | Pending      | [negociecoins pkg]     |
| Braziliex           | BTC, ETH, LTC               | Pending      |                        |
| Walltime            | BTC                         | Pending      |                        |
| NegocieCoins        | BTC                         | Pending      |                        |

## Authors

* Rodrigo Alves Vieira - <<rodrigoalvesvieira@protonmail.com>>

## Copyright

© 2018 Rodrigo Alves Vieira. All Rights Reserved.

[bitcointrade pkg]: https://github.com/rodrigoalvesvieira/operatordesk-bitcointrade
[negociecoins pkg]: https://github.com/rodrigoalvesvieira/operatordesk-negociecoins
