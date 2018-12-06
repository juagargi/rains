package integration

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	log "github.com/inconshreveable/log15"
	"github.com/netsec-ethz/rains/internal/pkg/connection"
	"github.com/netsec-ethz/rains/internal/pkg/libresolve"
	"github.com/netsec-ethz/rains/internal/pkg/message"
	"github.com/netsec-ethz/rains/internal/pkg/publisher"
	"github.com/netsec-ethz/rains/internal/pkg/query"
	"github.com/netsec-ethz/rains/internal/pkg/rainsd"
	"github.com/netsec-ethz/rains/internal/pkg/section"
	"github.com/netsec-ethz/rains/internal/pkg/token"
	"github.com/netsec-ethz/rains/internal/pkg/util"
	"github.com/netsec-ethz/rains/internal/pkg/zonefile"
)

func TestFullCoverage(t *testing.T) {
	h := log.CallerFileHandler(log.StdoutHandler)
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, h))
	//Start authoritative Servers and publish zonefiles to them
	rootServer := startAuthServer(t, "Root", nil)
	chServer := startAuthServer(t, "ch", []connection.Info{rootServer.Addr()})
	ethzChServer := startAuthServer(t, "ethz.ch", []connection.Info{rootServer.Addr()})
	log.Info("all authoritative servers successfully started")
	//Start client resolver
	cachingResolver, err := rainsd.New("testdata/conf/resolver.conf", "resolver")
	if err != nil {
		t.Fatalf("Was not able to create client resolver: %v", err)
	}
	cachingResolver.SetResolver(libresolve.New([]connection.Info{rootServer.Addr()}, nil,
		libresolve.Recursive, cachingResolver.Addr()))
	go cachingResolver.Start(false)
	time.Sleep(1000 * time.Millisecond)
	log.Info("caching server successfully started")

	//Send queries to client resolver and observe the recursive lookup results.
	queries := loadQueries(t)
	answers := loadAnswers(t)
	log.Info("begin sending queries which require recursive lookup")
	for i, query := range queries {
		sendQueryVerifyResponse(t, *query, cachingResolver.Addr(), answers[i])
	}

	//Shut down authoritative servers
	rootServer.Shutdown()
	chServer.Shutdown()
	ethzChServer.Shutdown()
	time.Sleep(500 * time.Millisecond)
	log.Info("begin sending queries which should be cached")
	//Send queries to client resolver and observe the cached results.
	for i, query := range queries {
		sendQueryVerifyResponse(t, *query, cachingResolver.Addr(), answers[i])
	}
}

func startAuthServer(t *testing.T, name string, rootServers []connection.Info) *rainsd.Server {
	server, err := rainsd.New("testdata/conf/namingServer"+name+".conf", "nameServerRoot")
	if err != nil {
		t.Fatal(fmt.Sprintf("Was not able to create %s server: ", name), err)
	}
	server.SetResolver(libresolve.New(rootServers, nil, libresolve.Recursive, server.Addr()))
	go server.Start(false)
	config, err := publisher.LoadConfig("testdata/conf/publisher" + name + ".conf")
	if err != nil {
		t.Fatal(fmt.Sprintf("Was not able to load %s publisher config: ", name), err)
	}
	pubServer := publisher.New(config)
	pubServer.Publish()
	time.Sleep(1000 * time.Millisecond)
	return server
}

func loadQueries(t *testing.T) []*query.Name {
	encoding, err := ioutil.ReadFile("testdata/messages/queries.txt")
	if err != nil {
		t.Fatal("Was not able to open queries.txt file: ", err)
	}
	zfParser := zonefile.Parser{}
	queries := zfParser.DecodeNameQueriesUnsafe(encoding)
	for _, q := range queries {
		q.Expiration = time.Now().Add(time.Hour).Unix()
	}
	return queries
}

func loadAnswers(t *testing.T) []section.WithSigForward {
	encoding, err := ioutil.ReadFile("testdata/messages/answers.txt")
	if err != nil {
		t.Fatal("Was not able to open answers.txt file: ", err)
	}
	zfParser := zonefile.Parser{}
	sections, err := zfParser.Decode(encoding)
	if err != nil {
		t.Fatal("Was not able to parse answers.txt file: ", err)
	}
	return sections
}

func sendQueryVerifyResponse(t *testing.T, query query.Name, connInfo connection.Info,
	answer section.Section) {
	msg := message.Message{Token: token.New(), Content: []section.Section{&query}}
	answerMsg, err := util.SendQuery(msg, connInfo, time.Second)
	if err != nil {
		t.Fatalf("could not send query or receive answer. query=%v err=%v", msg.Content, err)
	}
	if len(answerMsg.Content) != 1 {
		t.Fatalf("Got not exactly one answer for the query. msg=%v", answerMsg)
	}
	correctAnswer := false
	switch s := answerMsg.Content[0].(type) {
	case *section.Assertion:
		if a, ok := answer.(*section.Assertion); ok {
			correctAnswer = s.CompareTo(a) == 0
		}
	case *section.Shard:
		if a, ok := answer.(*section.Shard); ok {
			correctAnswer = s.CompareTo(a) == 0
		}
	case *section.Pshard:
		if a, ok := answer.(*section.Pshard); ok {
			correctAnswer = s.CompareTo(a) == 0
		}
	case *section.Zone:
		if a, ok := answer.(*section.Zone); ok {
			correctAnswer = s.CompareTo(a) == 0
		}
	default:
		t.Fatalf("Not yet implemented! So far only assertion, shard, pshard and zones are supported")
	}
	if !correctAnswer {
		t.Fatalf("Answer does not match expected result. actual=%v expected=%v",
			answerMsg.Content[0], answer)
	}
}
