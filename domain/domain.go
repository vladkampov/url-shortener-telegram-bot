package domain

import (
	"context"
	log "github.com/sirupsen/logrus"
	pb "github.com/vladkampov/url-shortener/service"
	"google.golang.org/grpc"
	"os"
	"strconv"
	"time"
)

var c pb.ShortenerClient

func GetUrls(userId int) (*pb.ArrayURLsReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	urls, err := c.GetMyUrls(ctx, &pb.UserIdRequest{UserId: strconv.FormatInt(int64(userId), 10)})
	if err != nil {
		return nil, err
	}
	log.Printf("URLs was successfully executed for user: %d", userId)
	return urls, nil
}

func SendUrl(url string, userId int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Shorten(ctx, &pb.URLRequest{Url: url, UserId: strconv.FormatInt(int64(userId), 10) })
	if err != nil {
		return "", err
	}
	log.Printf("URL was successfully shortened: %s", r.Url)
	return r.Url, nil
}

func RunDomainGrpcSession() (pb.ShortenerClient, error) {
	domainServiceUrl := os.Getenv("SHORTENER_DOMAIN_URL")

	if len(domainServiceUrl) == 0 {
		domainServiceUrl = "localhost:50051"
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(domainServiceUrl, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	c = pb.NewShortenerClient(conn)
	return c, nil
}
