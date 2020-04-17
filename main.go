package main

import (
	"log"
	"os"
	"simulator/messenger/eventlog"
	"simulator/messenger/passhash"
	"simulator/messenger/sessid"
	authsim "simulator/service/auth/simulator"
	gateway "simulator/service/gateway"
	messagingsim "simulator/service/messaging/simulator"
	userssim "simulator/service/users/simulator"
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

	usersService, err := userssim.New(l, passHashComparer)
	if err != nil {
		log.Fatalf("initializing users service: %s", err)
	}

	sessIDGen, err := sessid.NewGenerator(128)
	if err != nil {
		log.Fatalf("initializing session id generator: %s", err)
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

	log.Printf("listening on http://localhost:%s", port)
	gatewayServer.Addr = "localhost:" + port
	if err := gatewayServer.ListenAndServe(); err != nil {
		log.Fatalf("listening: %s", err)
	}

	// log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))
}
