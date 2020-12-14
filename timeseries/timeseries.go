package timeseries

import (
	"context"
	"fmt"
	"os"
	"math/big"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgtype"
	"time"
)

func InsertEntry(coin int64, tier int64, exchange int64, exchange_coin int64, query_time int64, count int, input *big.Float, vol *big.Float, b_book_len int, b_total_book_crypto *big.Float, b_total_book_fiat *big.Float, s_book_len int, s_total_book_crypto *big.Float, s_total_book_fiat *big.Float, s_count int, s_vol *big.Float) {
	ctx := context.Background()
	connStr := "postgres://rodrigoalvesvieira:@localhost:5432/fcpa"
	dbpool, dbpoolerr := pgxpool.Connect(ctx, connStr)
	if dbpoolerr != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", dbpoolerr)
		os.Exit(1)
	}
	defer dbpool.Close()

	s := &pgtype.Timestamp{Time: time.Unix(0, query_time).UTC(), Status: pgtype.Present}

	insert := `INSERT INTO price_queries (coin, tier, exchange, exchange_coin, query_time, buy_order_count, buy_input, buy_vol, buy_total_book_length, buy_total_book_crypto, buy_total_book_fiat, sell_total_book_length, sell_total_book_crypto, sell_total_book_fiat, sell_order_count, sell_vol) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);`

	//Execute INSERT command
	_, err := dbpool.Exec(ctx, insert, coin, tier, exchange, exchange_coin, s, count, input.String(), vol.String(), b_book_len, b_total_book_crypto.String(), b_total_book_fiat.String(), s_book_len, s_total_book_crypto.String(), s_total_book_fiat.String(), s_count, s_vol.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert data into database: %v\n", err)
		os.Exit(1)
	}
}

func Connect() {
	connStr := "postgres://rodrigoalvesvieira:@localhost:5432/fcpa"

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var coin int64
	var code string
	err = conn.QueryRow(context.Background(), "SELECT coin, code FROM coins WHERE coin=$1", 1).Scan(&coin, &code)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(coin, code)
}

type Result struct {
	Ids    []int64
	Values []*big.Float
}

func GetTiers(coin int64) (Result, error) {
	connStr := "postgres://rodrigoalvesvieira:@localhost:5432/fcpa"

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
