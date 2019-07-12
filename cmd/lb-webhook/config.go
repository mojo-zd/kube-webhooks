package main

import (
	"crypto/tls"
	"flag"

	"k8s.io/klog"
)

// Config contains the server (the webhook) cert and key.
type Config struct {
	CertFile string
	KeyFile  string
}

var (
	memory, cpu string
)

func (c *Config) addFlags() {
	flag.StringVar(&c.CertFile, "tls-cert-file", c.CertFile, ""+
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
		"after server cert).")
	flag.StringVar(&c.KeyFile, "tls-private-key-file", c.KeyFile, ""+
		"File containing the default x509 private key matching --tls-cert-file.")
	flag.StringVar(&memory, "memory", "200m",
		"default memory size of resource quota")
	flag.StringVar(&cpu, "cpu", "100m",
		"default cpu size of resource quota")
}

func configTLS(config Config) *tls.Config {
	sCert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		klog.Fatal(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{sCert},
	}
}
