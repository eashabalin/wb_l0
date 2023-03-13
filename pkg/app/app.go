package app

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
	"log"
	"os"
	"os/signal"
	"time"
	"wb_l0/pkg/handler/nats_streaming"
	"wb_l0/pkg/service"
)

const (
	subject     = "foo"
	natsURL     = ""
	clusterName = "test-cluster"
	clientName  = "myID"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run(port string) {
	services := service.NewService()
	handlers := nats_streaming.NewHandler(services)

	a.runServer()
	a.runStanSub(natsURL, handlers.CreateOrder)
}

func (a *App) runServer() {
}

func (a *App) runStanSub(natsServerURL string, handler func(msg *stan.Msg)) {
	var (
		clusterID, clientID = clusterName, clientName
		URL                 string
		userCreds           string
		showTime            bool
		qgroup              string
		unsubscribe         bool
		startSeq            uint64
		startDelta          string
		deliverAll          bool
		newOnly             bool
		deliverLast         bool
		durable             string
	)

	if natsServerURL == "" {
		URL = stan.DefaultNatsURL
	} else {
		URL = natsServerURL
	}

	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Streaming Example Subscriber")}
	// Use UserCredentials
	if userCreds != "" {
		opts = append(opts, nats.UserCredentials(userCreds))
	}

	// Connect to NATS
	nc, err := nats.Connect(URL, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, URL)
	}
	log.Printf("Connected to %s clusterID: [%s] clientID: [%s]\n", URL, clusterID, clientID)

	// Process Subscriber Options.
	startOpt := stan.StartAt(pb.StartPosition_NewOnly)
	if startSeq != 0 {
		startOpt = stan.StartAtSequence(startSeq)
	} else if deliverLast {
		startOpt = stan.StartWithLastReceived()
	} else if deliverAll && !newOnly {
		startOpt = stan.DeliverAllAvailable()
	} else if startDelta != "" {
		ago, err := time.ParseDuration(startDelta)
		if err != nil {
			sc.Close()
			log.Fatal(err)
		}
		startOpt = stan.StartAtTimeDelta(ago)
	}

	subj, i := subject, 0
	mcb := func(msg *stan.Msg) {
		i++
		handler(msg)
	}

	sub, err := sc.QueueSubscribe(subj, qgroup, mcb, startOpt, stan.DurableName(durable))
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}

	log.Printf("Listening on [%s], clientID=[%s], qgroup=[%s] durable=[%s]\n", subj, clientID, qgroup, durable)

	if showTime {
		log.SetFlags(log.LstdFlags)
	}

	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			// Do not unsubscribe a durable on exit, except if asked to.
			if durable == "" || unsubscribe {
				sub.Unsubscribe()
			}
			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
