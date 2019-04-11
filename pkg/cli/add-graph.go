package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/QMSTR/qmstr/pkg/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	yaml "gopkg.in/yaml.v2"
)

type missingPieces struct {
	File []*service.FileNode
}

var addGraphCmd = &cobra.Command{
	Use:   "add-graph [config_file]",
	Short: "Add partial graphs from a configuration file",
	Long: `Add missing files nodes and the connection between them.
Provide the extra information you want to include in the database,
through a configuration yaml file.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		setUpBuildService()
		if err := addMissingPieces(args); err != nil {
			Log.Fatalf("Describe failed: %v", err)
		}
		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(addGraphCmd)
}

func addMissingPieces(args []string) error {
	configFile := args[0]
	mp := &missingPieces{}
	if err := mp.readConfigFile(configFile); err != nil {
		return fmt.Errorf("Failed to read configuration from %s: %v", configFile, err)
	}

	stream, err := buildServiceClient.Build(context.Background())
	defer func() {
		res, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatalf("Failed to close the filenode stream: %v", err)
		}
		if !res.Success {
			log.Fatalln("Server filenode stream failed")
		}
	}()
	if err != nil {
		return err
	}
	for _, fnode := range mp.File {
		if err := stream.Send(fnode); err != nil {
			log.Fatalf("Failed to send filenode to server")
		}
		log.Printf("added filenode %s", fnode.Path)
	}
	return nil
}

func (mp *missingPieces) readConfigFile(configFile string) error {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("File %s not found", configFile)
	}

	cFile, err := os.Open(configFile)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(cFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, mp)
	if err != nil {
		return err
	}
	if mp.File == nil {
		return fmt.Errorf("no data found -- check indentation")
	}
	return nil
}
