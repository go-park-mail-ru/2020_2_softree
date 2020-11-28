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
	"server/canal/pkg/infrastructure/config"
	"server/profile/pkg/infrastructure/persistence"
	profile "server/profile/pkg/profile/gen"
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
	lis, err := net.Listen("tcp", ":8082")
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

	profile.RegisterProfileServiceServer(server, persistence.NewUserDBManager(db))

	fmt.Println("starting server at :8082")
	if err := server.Serve(lis); err != nil {
		log.Fatalln("cant listen port", err)
	}
}