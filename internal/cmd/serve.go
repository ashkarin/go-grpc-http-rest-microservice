package cmd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"

	"github.com/ashkarin/go-grpc-http-rest-microservice/internal/protocol"
	"github.com/ashkarin/go-grpc-http-rest-microservice/internal/usecases"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"

	samplesStorage "github.com/ashkarin/go-grpc-http-rest-microservice/internal/services/samples_storage"
)

type Config struct {
	Interface  string
	Port       string
	GrpcPort   string
	CACertFile string
	CAKeyFile  string

	DbDriver   string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbDatabase string
}

func Serve(cfg *Config) error {
	// Load CAKey and CACert if given
	var (
		keyPair  *tls.Certificate
		certPool *x509.CertPool
	)
	if cfg.CACertFile != "" && cfg.CAKeyFile != "" {
		keyPair, certPool = loadCA(cfg.CAKeyFile, cfg.CACertFile)

		if cfg.GrpcPort != "" {
			log.Info("The gRPC port will be ignored")
		}
		cfg.GrpcPort = cfg.Port
	}

	// Create a context
	ctx := context.Background()

	// Connect to DB
	conf := fmt.Sprintf("host=%v port=%v dbname=%v user=%v password='%v' sslmode=disable",
		cfg.DbHost, cfg.DbPort, cfg.DbDatabase, cfg.DbUser, cfg.DbPassword)

	db, err := sql.Open(cfg.DbDriver, conf)
	if err != nil {
		return err
	}

	// Create the usecase and inject it to the service
	sampleStoreUsecase, err := usecases.CreateSamplesStorageUsecase(db)
	if err != nil {
		return err
	}
	samplesStoreService := samplesStorage.NewServiceServer(sampleStoreUsecase)

	// Create GRPC server and gateway
	grpcServerAddress := fmt.Sprintf("%v:%v", cfg.Interface, cfg.GrpcPort)
	grpcServer := protocol.NewGRPCServer(grpcServerAddress, certPool)
	grpcGateway, opts := protocol.NewGatewayMux(ctx, grpcServerAddress, certPool)

	// Register
	samplesStorage.RegisterService(grpcServer, samplesStoreService)
	if err := samplesStorage.RegisterGateway(ctx, grpcGateway, grpcServerAddress, opts); err != nil {
		return err
	}

	// Serve
	return protocol.Serve(ctx, grpcServer, grpcGateway, cfg.Interface, cfg.Port, cfg.GrpcPort, keyPair)
}
