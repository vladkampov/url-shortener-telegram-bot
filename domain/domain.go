package domain

import (
	"context"
	log "github.com/sirupsen/logrus"
	pb "github.com/vladkampov/url-shortener/service"
	"google.golang.org/grpc"
	"os"
	"time"
)

var c pb.ShortenerClient

func SendUrl(url string) string {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Shorten(ctx, &pb.URLRequest{Url: url})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("URL was successfully shortened: %s", r.Url)
	return r.Url
}

func InitDomainGrpcSession() pb.ShortenerClient {
	domainServiceUrl := os.Getenv("SHORTENER_DOMAIN_PORT")

	if len(domainServiceUrl) == 0 {
		domainServiceUrl = "localhost:50051"
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(domainServiceUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c = pb.NewShortenerClient(conn)
	return c
}
