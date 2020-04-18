package main

import (
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

// DefaultPort defines the default fallback server port
const DefaultPort = "8080"

func main() {
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

	httpSrv := &http.Server{
		Addr:    "localhost:" + port,
		Handler: gatewayServer,
	}

	log.Printf("listening on http://localhost:%s", port)
	if err := httpSrv.ListenAndServe(); err != nil {
		log.Fatalf("listening: %s", err)
	}
}
