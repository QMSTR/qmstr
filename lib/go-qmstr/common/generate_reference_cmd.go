package common

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// CreateGenerateReferenceCmd creates the command to create the command line reference
func CreateGenerateReferenceCmd(theRootCmd *cobra.Command) *cobra.Command {

	// create a closure function that refers to theRootCmd:
	generateReference := func(cmd *cobra.Command, args []string) {
		print("Generating command line reference...\n")
		outputPath := args[0]
		if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
			log.Fatal(err)
		}

		const fmTemplate = `---
date: %s
title: "%s"
slug: "%s"
---
`

		filePrepender := func(filename string) string {
			now := time.Now().Format(time.RFC3339)
			name := filepath.Base(filename)
			base := strings.TrimSuffix(name, path.Ext(name))
			return fmt.Sprintf(fmTemplate, now, strings.Replace(base, "_", " ", -1), base)
		}

		linkHandler := func(name string) string {
			base := strings.TrimSuffix(name, path.Ext(name))
			return "../" + strings.ToLower(base) + "/"
		}
		if err := doc.GenMarkdownTreeCustom(theRootCmd, outputPath, filePrepender, linkHandler); err != nil {
			log.Fatal(err)
		}
	}

	var devGenerateReferenceCmd = &cobra.Command{
		Use:   "dev-generate-reference [path]",
		Short: "Generate command line reference in the specified path",
		Long:  `Generate the command line reference in Markdown format. The pages will contain metadata to integrate into a Hugo web site.`,
		Args:  cobra.ExactArgs(1),
		Run:   generateReference,
	}
	return devGenerateReferenceCmd
}
