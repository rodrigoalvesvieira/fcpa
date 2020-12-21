package timeseries

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"fcpa/envars"
)

type PriceQuery struct {
	Coin                    int64
	Tier                    int64
	Exchange_buy            int64
	Exchange_sell           int64
	Exchange_buy_coin       int64
	Exchange_sell_coin      int64
	Exchange_sell_timestamp int64
	Buy_query_time          int64
	Sell_query_time         int64
	Buy_input               *big.Float
	Buy_count               int
	Buy_vol                 *big.Float
	Sell_input              *big.Float
	Sell_count              int
	Sell_vol                *big.Float
	Buy_total_book_length   int
	Buy_total_book_crypto   *big.Float
	Buy_total_book_fiat     *big.Float
	Sell_total_book_length  int
	Sell_total_book_crypto  *big.Float
	Sell_total_book_fiat    *big.Float
	Exchange_rate           *big.Float
}

func InsertPriceQuery(entry *PriceQuery) {
	ctx := context.Background()
	connStr := envars.GetDBUrl()
	dbpool, dbpoolerr := pgxpool.Connect(ctx, connStr)
	if dbpoolerr != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", dbpoolerr)
		os.Exit(1)
	}
	defer dbpool.Close()

	//insert := `INSERT INTO price_queries (coin, tier, exchange_buy, exchange_sell, exchange_buy_coin,  exchange_sell_coin, exchange_sell_timestamp, exchange_buy_query_time, exchange_sell_query_time, buy_input, buy_order_count, buy_vol, sell_input, sell_order_count, sell_vol, buy_total_book_length, buy_total_book_crypto, buy_total_book_fiat, sell_total_book_length, sell_total_book_crypto, sell_total_book_fiat, exchange_rate) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22);`
	insert := `INSERT INTO price_queries (coin, tier, exchange_buy, exchange_sell, exchange_buy_coin,  exchange_sell_coin, exchange_buy_query_time, exchange_sell_query_time, buy_input, buy_order_count, buy_vol, sell_input, sell_order_count, sell_vol, buy_total_book_length, buy_total_book_crypto, buy_total_book_fiat, sell_total_book_length, sell_total_book_crypto, sell_total_book_fiat, exchange_rate) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21);`

	buyquerytime := &pgtype.Timestamp{Time: time.Unix(0, entry.Buy_query_time).UTC(), Status: pgtype.Present}
	sellquerytime := &pgtype.Timestamp{Time: time.Unix(0, entry.Sell_query_time).UTC(), Status: pgtype.Present}
	//selltimestamp := &pgtype.Timestamp{Time: time.Unix(0, entry.Exchange_sell_timestamp).UTC(), Status: pgtype.Present}

	// Execute INSERT command
	_, err := dbpool.Exec(
		ctx,
		insert,
		entry.Coin,
		entry.Tier,
		entry.Exchange_buy,
		entry.Exchange_sell,
		entry.Exchange_buy_coin,
		entry.Exchange_sell_coin,
		buyquerytime,
		sellquerytime,
		entry.Buy_input.String(),
		entry.Buy_count,
		entry.Buy_vol.String(),
		entry.Sell_input.String(),
		entry.Sell_count,
		entry.Sell_vol.String(),
		entry.Buy_total_book_length,
		entry.Buy_total_book_crypto.String(),
		entry.Buy_total_book_fiat.String(),
		entry.Sell_total_book_length,
		entry.Sell_total_book_crypto.String(),
		entry.Sell_total_book_fiat.String(),
		entry.Exchange_rate.String(),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert data into database: %v\n", err)
		os.Exit(1)
	}
}

type Result struct {
	Ids    []int64
	Values []*big.Float
}

func GetTiers(coin int64) (Result, error) {
	connStr := envars.GetDBUrl()

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var id int64
	var amount float64
	var ids []int64
	var result []*big.Float

	rows, rowsErr := conn.Query(context.Background(), "SELECT id, amount FROM tiers WHERE active = TRUE AND coin = $1", &coin)
	if rowsErr != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", rowsErr)
		os.Exit(1)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &amount)
		if err != nil {
			return Result{}, err
		}
		ids = append(ids, id)
		result = append(result, big.NewFloat(amount))
	}

	return Result{Ids: ids, Values: result}, nil
}
