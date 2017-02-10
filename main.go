package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/deskr/grpc-health-checker/health"

	"google.golang.org/grpc"
)

const Name = "grpc-health-checker"
const Version = "0.1.0"

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

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Dial error:", err)
		os.Exit(1)
	}

	client := health.NewHealthClient(conn)

	res, err := client.Check(context.Background(), &health.HealthCheckRequest{
		Service: *service,
	})
	if err != nil {
		fmt.Println("Check error:", err)
		os.Exit(1)
	}

	fmt.Println("Status:", res.Status)

	if res.Status != health.HealthCheckResponse_SERVING {
		os.Exit(1)
	}
}
