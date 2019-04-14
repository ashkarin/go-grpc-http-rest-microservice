package protocol

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// HandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func HandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

// Serve GRPC Server and it's Gateway
func Serve(
	ctx context.Context,
	grpcServer *grpc.Server,
	grpcGateway *runtime.ServeMux,
	iface string,
	port string,
	grpcPort string,
	keypair *tls.Certificate,
) error {
	if grpcPort == "" && keypair == nil {
		return errors.New("Either the 'grpcPort' or 'keypair' must be given.")
	}
	if keypair != nil {
		return ServeMultiplexed(ctx, grpcServer, grpcGateway, iface, port, keypair)
	}

	// Serve the GPRC Server and the Gateway on different ports
	go func() {
		_ = ServeGatewayServer(ctx, grpcGateway, port)
	}()
	return ServeGRPCServer(ctx, grpcServer, grpcPort)
}

func ServeGatewayServer(ctx context.Context, mux *runtime.ServeMux, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	server := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down HTTP/REST gateway...")
		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = server.Shutdown(ctx)
	}()

	// start HTTP/REST gateway server
	log.Println("starting HTTP/REST gateway...")
	return server.ListenAndServe()
}

func ServeGRPCServer(ctx context.Context, grpcServer *grpc.Server, grpcPort string) error {
	// Creat listener
	listen, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Panic(err)
	}

	// Run GRPC Server with a graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")

			grpcServer.GracefulStop()

			<-ctx.Done()
		}
	}()

	log.Println("starting gRPC server...")
	return grpcServer.Serve(listen)
}

func ServeMultiplexed(
	ctx context.Context,
	grpcServer *grpc.Server,
	grpcGateway *runtime.ServeMux,
	iface string,
	port string,
	keypair *tls.Certificate,
) error {
	address := fmt.Sprintf("%v:%v", iface, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	log.Info("LISTEN ", address)
	// Create a multiplexed server
	server := &http.Server{
		Addr:    address,
		Handler: HandlerFunc(grpcServer, grpcGateway),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*keypair},
			NextProtos:   []string{"h2"},
		},
	}

	// Run the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down...")
			_ = server.Shutdown(ctx)
		}
	}()

	log.Info("starting gRPC server and HTTP gateway on: ", address)
	return server.Serve(tls.NewListener(listener, server.TLSConfig))
}
