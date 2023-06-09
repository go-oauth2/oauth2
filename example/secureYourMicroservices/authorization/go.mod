module server

go 1.20

replace github.com/go-oauth2/oauth2/v4 => ../../../

require (
	github.com/go-oauth2/oauth2/v4 v4.4.3
	github.com/jackc/pgx/v4 v4.18.1
	github.com/vgarvardt/go-oauth2-pg/v4 v4.4.3
	github.com/vgarvardt/go-pg-adapter v1.0.0
)

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.2 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/jmoiron/sqlx v1.3.4 // indirect
	github.com/vgarvardt/pgx-helpers/v4 v4.0.0-20200225100150-876aee3d1a22 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/text v0.7.0 // indirect
)
