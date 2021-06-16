package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	gamepb "pb"
	"sync"
	"time"
)

const (
	address     = "localhost:50051"
	port        = ":50051"
	defaultName = "world"
)

type server struct {
	gamepb.UnimplementedGreeterServer
	mu         sync.Mutex
	routeNotes map[string][]*gamepb.RouteNote
}

func (s *server) SayHello(ctx context.Context, in *gamepb.HelloRequest) (*gamepb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &gamepb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) RouteChat(stream gamepb.Greeter_RouteChatServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		key := fmt.Sprintf("%d %d", in.Location.X, in.Location.Y)

		s.mu.Lock()
		s.routeNotes[key] = append(s.routeNotes[key], in)
		// Note: this copy prevents blocking other clients while serving this one.
		// We don't need to do a deep copy, because elements in the slice are
		// insert-only and never modified.
		rn := make([]*gamepb.RouteNote, len(s.routeNotes[key]))
		copy(rn, s.routeNotes[key])
		s.mu.Unlock()

		for _, note := range rn {
			if err := stream.Send(note); err != nil {
				return err
			}
		}
	}
}

func main() {
	listen()
	time.Sleep(time.Millisecond * 1000)

	connect()
}

func listen() {
	lis, err := net.Listen("tcp", port)
	if err != nil {

	}
	s := grpc.NewServer()
	gamepb.RegisterGreeterServer(s, &server{})

	go func() {
		if err := s.Serve(lis); err != nil {

		}
	}()
}

func connect() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {

	}
	defer conn.Close()
	c := gamepb.NewGreeterClient(conn)

	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &gamepb.HelloRequest{Name: name})
	if err != nil {

	}
	log.Printf("Greeting: %s", r.GetMessage())

	// stream
	notes := []*gamepb.RouteNote{
		{Location: &gamepb.Point{X: 0, Y: 1}, Message: "First message"},
		{Location: &gamepb.Point{X: 0, Y: 2}, Message: "Second message"},
		{Location: &gamepb.Point{X: 0, Y: 3}, Message: "Third message"},
		{Location: &gamepb.Point{X: 0, Y: 1}, Message: "Fourth message"},
		{Location: &gamepb.Point{X: 0, Y: 2}, Message: "Fifth message"},
		{Location: &gamepb.Point{X: 0, Y: 3}, Message: "Sixth message"},
	}
	stream, err := c.RouteChat(ctx)
	if err != nil {

	}
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				break
			}
			log.Printf("Got message %s at point(%d, %d)", in.Message, in.Location.X, in.Location.Y)
		}
	}()

	for _, note := range notes {
		if err := stream.Send(note); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}
	}
}
