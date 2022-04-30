package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/twitchtv/twirp"

	"go.ntppool.org/monitor/api/pb"
	apitls "go.ntppool.org/monitor/api/tls"
	"go.ntppool.org/monitor/ntpdb"
	sctx "go.ntppool.org/monitor/server/context"
	"go.ntppool.org/monitor/server/metrics"
	twirpmetrics "go.ntppool.org/monitor/server/metrics/twirp"
	"go.ntppool.org/monitor/server/vault"
)

type Server struct {
	cfg    *Config
	tokens *vault.TokenManager
	m      *metrics.Metrics
	db     *ntpdb.Queries
	dbconn *sql.DB
}

type Config struct {
	DeploymentEnv string
	Listen        string
	CertProvider  apitls.CertificateProvider
}

func NewServer(cfg Config, dbconn *sql.DB) (*Server, error) {
	db := ntpdb.New(dbconn)

	tm, err := vault.New("monitor-tokens", cfg.DeploymentEnv)
	if err != nil {
		return nil, err
	}

	metrics := metrics.New()

	return &Server{
		cfg:    &cfg,
		db:     db,
		dbconn: dbconn,
		tokens: tm,
		m:      metrics,
	}, nil
}

func (srv *Server) getVerifiedCert(verifiedChains [][]*x509.Certificate) (*x509.Certificate, string) {
	for _, chain := range verifiedChains {
		for _, cert := range chain {
			for _, name := range cert.DNSNames {
				// log.Printf("verified cert, dns name: %+v", cert.DNSNames)
				// log.Printf("issuer: %s", cert.Issuer)
				if strings.HasSuffix(name, ".mon.ntppool.dev") {
					return cert, name
				}
			}
		}
	}
	return nil, ""
}

func (srv *Server) verifyClient(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	cert, _ := srv.getVerifiedCert(verifiedChains)
	if cert != nil {
		return nil
	}
	return errors.New("no valid certificate found")
}

func (srv *Server) Run() error {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	// logrusEntry := logrus.NewEntry(logger)

	capool, err := apitls.CAPool()
	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{
		MinVersion:            tls.VersionTLS12,
		ClientCAs:             capool,
		ClientAuth:            tls.RequireAndVerifyClientCert,
		GetCertificate:        srv.cfg.CertProvider.GetCertificate,
		VerifyPeerCertificate: srv.verifyClient,
	}

	twirpHandler := pb.NewMonitorServer(srv,
		twirp.WithServerPathPrefix("/api/v1"),
		twirp.WithServerHooks(NewLoggingServerHooks()),
		twirp.WithServerHooks(twirpmetrics.NewServerHooks(srv.m.Registry())),
	)

	mux := http.NewServeMux()
	mux.Handle(twirpHandler.PathPrefix(),
		srv.certificateMiddleware(
			WithUserAgent(
				twirpHandler,
			),
		),
	)

	logger.Infof("starting server")

	metricsServer := &http.Server{
		Addr:    ":9000",
		Handler: srv.m.Handler(),
	}

	go func() {
		log.Printf("metrics server: %s", metricsServer.ListenAndServe())
	}()

	server := &http.Server{
		Addr: ":8000",

		TLSConfig: tlsConfig,
		Handler:   mux,

		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       20 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       240 * time.Second,
	}

	return server.ListenAndServeTLS("", "")
}

func (srv *Server) certificateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, name := srv.getVerifiedCert(r.TLS.VerifiedChains)

		ctx := context.WithValue(r.Context(), sctx.CertificateKey, name)
		rctx := r.WithContext(ctx)
		next.ServeHTTP(w, rctx)
	})
}

func NewLoggingServerHooks() *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestRouted: func(ctx context.Context) (context.Context, error) {
			method, _ := twirp.MethodName(ctx)
			log.Println("Method: " + method)
			return ctx, nil
		},
		Error: func(ctx context.Context, twerr twirp.Error) context.Context {
			log.Println("Error: " + string(twerr.Code()))
			return ctx
		},
		// ResponseSent: func(ctx context.Context) {
		// 	log.Println("Response Sent (error or success)")
		// },
	}
}
