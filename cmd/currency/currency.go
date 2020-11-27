package main

import (
	"database/sql"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"server/src/canal/pkg/infrastructure/config"
	"server/src/currency/pkg/infrastructure/financial"
	"server/src/currency/pkg/infrastructure/persistence"
	currency "server/src/currency/pkg/currency/gen"
)

func init() {
	pflag.StringP("viper", "c", "", "path to viper file")
	pflag.BoolP("help", "h", false, "usage info")

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	if viper.GetBool("help") {
		pflag.Usage()
		os.Exit(0)
	}

	if viper.GetString("viper") == "" {
		_, _ = fmt.Fprintln(os.Stderr, "There is must explicitly specify the viper file")
		pflag.Usage()
		os.Exit(1)
	}

	if err := config.ParseConfig(
		viper.GetString("viper"),
		map[string]interface{}{
			"postgres": map[string]interface{}{
				"host":     "127.0.0.1",
				"port":     5432,
				"db":       "db",
				"user":     "user",
				"password": "",
			},
		}); err != nil {
		log.Fatalln("Error during parse defaults", err)
	}
}

func main() {
	lis, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	server := grpc.NewServer()

	db, err := sql.Open("postgres", viper.GetString("postgres.URL"))
	if err != nil {
		log.Fatalln("cant listen port", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln("cant listen port", err)
	}
	db.SetMaxOpenConns(10)
	currency.RegisterCurrencyServiceServer(server, persistence.NewRateDBManager(db, financial.NewForexAPI()))

	fmt.Println("starting server at :8083")
	server.Serve(lis)
}