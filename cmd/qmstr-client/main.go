//go:generate protoc -I ../../pkg/buildservice --go_out=plugins=grpc:../../pkg/buildservice ../../pkg/buildservice/buildservice.proto

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"text/scanner"

	"golang.org/x/net/context"

	pb "github.com/QMSTR/qmstr/pkg/buildservice"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

var (
	buildServiceClient pb.BuildServiceClient
	conn               *grpc.ClientConn
	currentCommand     string
)

func main() {
	var pipe bool
	pflag.BoolVarP(&pipe, "pipe", "p", false, "read commands from piped stdin")
	pflag.Parse()

	// Set up server connection
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to master: %v", err)
	}
	defer conn.Close()
	buildServiceClient = pb.NewBuildServiceClient(conn)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			handleCommand("help")
		}
	}()

	if pipe {
		runPipe()
	} else {
		cli()
	}
}

func runPipe() {
	reader := bufio.NewScanner(os.Stdin)
	reader.Scan()
	commands := reader.Text()

	var commandScanner scanner.Scanner
	commandScanner.Init(strings.NewReader(commands))
	commandScanner.Filename = "example"
	for tok := commandScanner.Scan(); tok != scanner.EOF; tok = commandScanner.Scan() {
		handleCommand(commandScanner.TokenText())
	}
}

func cli() {
	fmt.Println("QMSTR cli")
	for currentCommand != "exit" {
		fmt.Scanln(&currentCommand)
		handleCommand(currentCommand)
	}

}

func handleCommand(command string) {
	switch command {
	case "exit":
		fmt.Println("Good bye")
		os.Exit(0)
	case "quit":
		buildServiceClient.Quit(context.Background(), &pb.QuitMessage{Kill: false})
	case "forcequit":
		buildServiceClient.Quit(context.Background(), &pb.QuitMessage{Kill: true})
	case "", "help":
		fmt.Println("Those are the supported commands: exit, quit, forcequit, help")
	default:
		fmt.Printf("unknown command %s\n", command)
	}

}
