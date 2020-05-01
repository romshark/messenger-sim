package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/romshark/messenger-sim/messenger/eventlog"
	"github.com/romshark/messenger-sim/messenger/passhash"
	"github.com/romshark/messenger-sim/messenger/sessid"
	authsim "github.com/romshark/messenger-sim/service/auth/simulator"
	gateway "github.com/romshark/messenger-sim/service/gateway"
	messagingsim "github.com/romshark/messenger-sim/service/messaging/simulator"
	userssim "github.com/romshark/messenger-sim/service/users/simulator"
)

func main() {
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	l := eventlog.New()

	passHashComparer := passhash.NewBcrypt()

	sessIDGen, err := sessid.NewGenerator(128)
	if err != nil {
		log.Fatalf("initializing session id generator: %s", err)
	}

	usersService, err := userssim.New(l, passHashComparer)
	if err != nil {
		log.Fatalf("initializing users service: %s", err)
	}

	authService, err := authsim.New(l, sessIDGen, passHashComparer)
	if err != nil {
		log.Fatalf("initializing sessions service: %s", err)
	}

	messagingService, err := messagingsim.New(l)
	if err != nil {
		log.Fatalf("initializing messaging service: %s", err)
	}

	gatewayServer, err := gateway.NewServer(
		usersService,
		authService,
		messagingService,
	)
	if err != nil {
		log.Fatalf("initializing gateway server")
	}

	httpsSrv := &http.Server{
		Addr:    *flagHostAddr + ":" + port,
		Handler: gatewayServer,
	}

	log.Printf("listening on https://%s:%s", *flagHostAddr, port)
	if err := httpsSrv.ListenAndServeTLS(
		*flagCertFilePath,
		*flagKeyFilePath,
	); err != nil {
		log.Fatalf("listening: %s", err)
	}
}

// DefaultPort defines the default fallback server port
const DefaultPort = "443"

var (
	flagCertFilePath = flag.String(
		"crt",
		"ssl/dev.messenger.org.crt",
		"SSL public certificate file path",
	)
	flagKeyFilePath = flag.String(
		"pkey",
		"ssl/dev.messenger.org.key",
		"SSL private key file path",
	)
	flagHostAddr = flag.String(
		"host",
		"dev.messenger.org",
		"host address",
	)
)
