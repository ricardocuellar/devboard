module github.com/ricardocuellar/devboard/services/notification

go 1.26.4

require (
	github.com/jackc/pgx/v5 v5.11.0
	github.com/ricardocuellar/devboard/contracts v0.0.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	golang.org/x/text v0.40.0 // indirect
)

replace github.com/ricardocuellar/devboard/contracts => ../../contracts
