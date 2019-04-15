package commands

import (
	"github.com/ashkarin/go-grpc-http-rest-microservice/internal/cmd"
	"github.com/spf13/cobra"
)

var (
	cfg cmd.Config

	cmdServer = &cobra.Command{
		Use:   "serve",
		Short: "Start a server",
		RunE:  serve,
	}
)

func init() {
	cmdServer.Flags().StringVarP(&cfg.Interface, "interface", "", "localhost", "interface to which the server will bind")
	cmdServer.Flags().StringVarP(&cfg.Port, "port", "", "8080", "TCP port to listen by HTTP/REST gateway")
	cmdServer.Flags().StringVarP(&cfg.GrpcPort, "grpc-port", "", "9090", "TCP port to listen by gRPC server")
	cmdServer.Flags().StringVarP(&cfg.CACertFile, "cert-file", "", "", "certificate used for SSL/TLS connections to the server.")
	cmdServer.Flags().StringVarP(&cfg.CAKeyFile, "key-file", "", "", "key for the certificate.")
	cmdServer.Flags().StringVarP(&cfg.DbDriver, "db-driver", "", "postgres", "SQL driver used to connect to the database")
	cmdServer.Flags().StringVarP(&cfg.DbHost, "db-host", "", "127.0.0.1", "Hostname of the database server")
	cmdServer.Flags().StringVarP(&cfg.DbPort, "db-port", "", "5432", "Port")
	cmdServer.Flags().StringVarP(&cfg.DbDatabase, "db-name", "", "", "Name of the database")
	cmdServer.Flags().StringVarP(&cfg.DbUser, "db-user", "", "", "Username")
	cmdServer.Flags().StringVarP(&cfg.DbPassword, "db-password", "", "", "Password")

}

func serve(ccmd *cobra.Command, args []string) error {
	return cmd.Serve(&cfg)
}
