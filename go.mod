module operatordesk

go 1.15

require (
	github.com/ajph/bitstamp-go v0.0.0-20180531100631-da686669997a
	github.com/beldur/kraken-go-api-client v0.0.0-20180126083054-8d8ccfe4cc60
	github.com/fatih/color v1.7.0
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/jackc/pgtype v1.6.2
	github.com/jackc/pgx/v4 v4.10.0
	github.com/joho/godotenv v1.3.0
	golang.org/x/text v0.3.4
)

replace github.com/go-chi/chi => ./packages/chi
