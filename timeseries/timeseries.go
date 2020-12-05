package timeseries

import (
	"context"
	"fmt"
	"os"
	"github.com/jackc/pgx/v4"
)

func InsertEntry() {
	fmt.Println("INSERTING NEW TIMESERIES ENTRY")
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
