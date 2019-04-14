package protocol

import (
	"context"
	"crypto/tls"
	"crypto/x509"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewGRPCServer(serveAddress string, certPool *x509.CertPool) *grpc.Server {
	var opts []grpc.ServerOption
	if certPool != nil {
		creds := grpc.Creds(credentials.NewClientTLSFromCert(certPool, serveAddress))
		opts = []grpc.ServerOption{creds}
	}
	server := grpc.NewServer(opts...)
	return server
}

func NewGatewayMux(ctx context.Context, serverName string, certPool *x509.CertPool) (*runtime.ServeMux, []grpc.DialOption) {
	var mux *runtime.ServeMux
	var opts []grpc.DialOption

	if certPool == nil {
		mux = runtime.NewServeMux()
		opts = []grpc.DialOption{grpc.WithInsecure()}
	} else {
		muxopts := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true})
		mux = runtime.NewServeMux(muxopts)

		// Ensure that ServerName corresponds one in cerificate
		creds := credentials.NewTLS(&tls.Config{
			ServerName: serverName,
			RootCAs:    certPool,
		})
		opts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	}
	return mux, opts
}
