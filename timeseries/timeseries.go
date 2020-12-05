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

func InsertEntry(coin int64, tier int64, exchange int64, exchange_coin int64, query_time int64, count int, input *big.Float, vol *big.Float) {
	ctx := context.Background()
	connStr := "postgres://rodrigoalvesvieira:@localhost:5432/fcpa"
	dbpool, dbpoolerr := pgxpool.Connect(ctx, connStr)
	if dbpoolerr != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", dbpoolerr)
		os.Exit(1)
	}
	defer dbpool.Close()

	s := &pgtype.Timestamp{Time: time.Unix(0, query_time).UTC(), Status: pgtype.Present}

	insert := `INSERT INTO price_queries (coin, tier, exchange, exchange_coin, query_time, buy_order_count, buy_input, buy_vol) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

	//Execute INSERT command
	_, err := dbpool.Exec(ctx, insert, coin, tier, exchange, exchange_coin, s, count, input.String(), vol.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert data into database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Inserted price query into database")
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
