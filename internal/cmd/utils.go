package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

func loadCA(keyFilepath, certFilepath string) (*tls.Certificate, *x509.CertPool) {
	serverCrt, err := ioutil.ReadFile(certFilepath)
	if err != nil {
		log.Fatal(err)
	}
	serverKey, err := ioutil.ReadFile(keyFilepath)
	if err != nil {
		log.Fatal(err)
	}

	pair, err := tls.X509KeyPair(serverCrt, serverKey)
	if err != nil {
		log.Fatal(err)
	}

	keyPair := &pair
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(serverCrt); !ok {
		log.Fatal("bad certs")
	}

	return keyPair, certPool
}
