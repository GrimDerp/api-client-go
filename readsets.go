// Copyright 2014 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	genomics "github.com/googlegenomics/api-client-go/v1beta"
	"github.com/spf13/cobra"
)

func printReadsetReferenceSequence(refSequence *genomics.ReferenceSequence,
	indentLevel int) {
	indent(indentLevel)
	fmt.Printf("Sequence Name: %v\tUri: %v\n",
		refSequence.Name, refSequence.Uri)
}

func printReadsetFileData(fileData *genomics.HeaderSection, indentLevel int) {
	indent(indentLevel)
	fmt.Printf("File: %v\n", fileData.FileUri)
	for _, sequence := range fileData.RefSequences {
		printReadsetReferenceSequence(sequence, indentLevel+1)
	}
}

func printReadset(readset *genomics.Readset, indentLevel int) {
	indent(indentLevel)
	fmt.Printf("ID: %v\tName: %v\tNumber of Files: %v\n",
		readset.Id, readset.Name, len(readset.FileData))
	for _, fileData := range readset.FileData {
		printReadsetFileData(fileData, indentLevel+1)
	}
}

type paramReadsetsSearch struct {
	datasetIds string
	pageToken  string
}

func (p *paramReadsetsSearch) command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for readsets in dataset(s)",
		Run:   p.run,
	}
	cmd.Flags().StringVarP(&p.datasetIds, "dataset-ids", "", "",
		"Dataset IDs (comma separated), 376902546192 for 1000 Genomes")
	cmd.Flags().StringVarP(&p.pageToken, "page-token", "", "",
		"Page token to continue with previous large result")
	return cmd
}

func (p *paramReadsetsSearch) run(cmd *cobra.Command, args []string) {
	datasetIds := strings.Split(p.datasetIds, ",")
	name := strings.Join(args, " ")

	client := NewApiClient()

	fmt.Fprintln(os.Stderr, "Searching for readsets...")
	res, err := client.Readsets.Search(&genomics.SearchReadsetsRequest{
		DatasetIds: datasetIds,
		PageToken:  p.pageToken,
		Name:       name,
	}).Do()
	if err != nil {
		log.Fatal(err)
	}

	if len(res.Readsets) == 0 {
		fmt.Fprintf(os.Stderr, "No match for %s.\n", name)
		return
	}

	for _, rs := range res.Readsets {
		printReadset(rs, 0)
	}

	if res.NextPageToken != "" {
		fmt.Printf("Next Page Token: %v\n", res.NextPageToken)
	}
}

type paramReadsetsGet struct {
	id string
}

func (p *paramReadsetsGet) command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a specified readset",
		Run:   p.run,
	}
	return cmd
}

func (p *paramReadsetsGet) run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Missing ID")
		return
	}

	id := args[0]

	client := NewApiClient()

	fmt.Fprintf(os.Stderr, "Retrieving information of readset %s\n", id)
	res, err := client.Readsets.Get(id).Do()
	if err != nil {
		log.Fatal(err)
	}

	printReadset(res, 0)
}

var readsetsCmd *cobra.Command

func init() {
	readsetsCmd = &cobra.Command{
		Use:   "readsets",
		Short: "Readsets functions",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	readsetsGet := &paramReadsetsGet{}
	readsetsCmd.AddCommand(readsetsGet.command())

	readsetsSearch := &paramReadsetsSearch{}
	readsetsCmd.AddCommand(readsetsSearch.command())
}
