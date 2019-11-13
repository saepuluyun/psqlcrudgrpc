package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"

	//postgres driver
	_ "github.com/lib/pq"

	"aniqma/aniqma/crudgrpc/pkg/protocol/grpc"
	v1 "aniqma/aniqma/crudgrpc/pkg/service/v1"
)

// Config is configuration for Server MYSQL
// type Config struct {
// 	// gRPC server start parameters section
// 	// gRPC is TCP port to listen by gRPC server
// 	GRPCPort string

// 	// DB Datastore parameters section
// 	// DatastoreDBHost is host of database
// 	DatastoreDBHost string
// 	// DatastoreDBUser is username to connect to database
// 	DatastoreDBUser string
// 	// DatastoreDBPassword password to connect to database
// 	DatastoreDBPassword string
// 	// DatastoreDBSchema is schema of database
// 	DatastoreDBSchema string
// }

// Config is configuration for Server POSTGRES
type Config struct {
	// gRPC server start parameters section
	// gRPC is TCP port to listen by gRPC server
	GRPCPort string

	// DB Datastore parameters section
	// DatastoreDBName is database name
	DatastoreDBName string
	// DatastoreDBUser is username to connect to database
	DatastoreDBUser string
	// DatastoreDBPassword password to connect to database
	DatastoreDBPassword string
	// DatastoreDBHost is host of database
	DatastoreDBHost string
	// DatastoreDBSsl is SSL of database
	DatastoreDBSslmode string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get configuration
	var cfg Config

	//mysql
	// flag.StringVar(&cfg.GRPCPort, "grpc-port", "9090", "gRPC port to bind")
	// flag.StringVar(&cfg.DatastoreDBHost, "db-host", "localhost", "Database host")
	// flag.StringVar(&cfg.DatastoreDBUser, "db-user", "root", "Database user")
	// flag.StringVar(&cfg.DatastoreDBPassword, "db-password", "", "Database password")
	// flag.StringVar(&cfg.DatastoreDBSchema, "db-schema", "example", "Database schema")
	// flag.Parse()

	flag.StringVar(&cfg.GRPCPort, "grpc-port", "9090", "gRPC port to bind")
	flag.StringVar(&cfg.DatastoreDBName, "db-name", "example", "Database name")
	flag.StringVar(&cfg.DatastoreDBUser, "db-user", "postgres", "Database user")
	flag.StringVar(&cfg.DatastoreDBPassword, "db-password", "localhost", "Database password")
	flag.StringVar(&cfg.DatastoreDBHost, "db-host", "localhost", "Database host")
	flag.StringVar(&cfg.DatastoreDBSslmode, "db-sslmode", "disable", "Database sslmode")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	// add MySQL driver specific parameter to parse date/time
	// Drop it for another database
	// param := "parseTime=true"

	//mysql
	// dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
	// 	cfg.DatastoreDBUser,
	// 	cfg.DatastoreDBPassword,
	// 	cfg.DatastoreDBHost,
	// 	cfg.DatastoreDBSchema,
	// 	param)

	//postgres
	dsn := fmt.Sprintf("dbname=%s user=%s password=%s host=%s sslmode=%s",
		cfg.DatastoreDBName,
		cfg.DatastoreDBUser,
		cfg.DatastoreDBPassword,
		cfg.DatastoreDBHost,
		cfg.DatastoreDBSslmode)
	// db, err := sql.Open("mysql", dsn) //mysql
	db, err := sql.Open("postgres", dsn) //postgres
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	v1API := v1.NewUsersServiceServer(db)

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
