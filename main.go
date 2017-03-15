package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"

	health "github.com/deskr/grpc-health-checker/grpc_health_v1"
)

const Name = "grpc-health-checker"
const Version = "0.1.3"

var addr = flag.String("address", "", "Address of GRPC service")
var service = flag.String("service", "", "Name of service")
var ver = flag.Bool("version", false, "Prints version")

func main() {
	flag.Parse()

	if *ver {
		fmt.Println("Version:", Version)
		return
	}

	if *addr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	conn, err := grpc.Dial(*addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second*3),
		grpc.FailOnNonTempDialError(true))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	client := health.NewHealthClient(conn)

	res, err := client.Check(context.Background(), &health.HealthCheckRequest{
		Service: *service,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if res.Status != health.HealthCheckResponse_SERVING {
		fmt.Fprintf(os.Stderr, "Status: %v\n", res.Status)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "Status: %v\n", res.Status)
}
