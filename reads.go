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

func printRead(read *genomics.Read, indentLevel int) {
	indent(indentLevel)
	fmt.Printf("ID: %s, Name: %s\n", read.Id, read.Name)
}

type paramReadsSearch struct {
	readsetIds    string
	pageToken     string
	sequenceStart uint64
	sequenceEnd   uint64
}

func (p *paramReadsSearch) command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for reads in readset(s)",
		Run:   p.run,
	}
	cmd.Flags().StringVarP(&p.readsetIds, "readset-ids", "", "",
		"Readset IDs (comma separated)")
	cmd.Flags().StringVarP(&p.pageToken, "page-token", "", "",
		"Page token to continue with previous large result")
	cmd.Flags().Uint64VarP(&p.sequenceStart, "start", "", 0,
		"Starting position")
	cmd.Flags().Uint64VarP(&p.sequenceEnd, "end", "", ^uint64(0),
		"Starting position")
	return cmd
}

func (p *paramReadsSearch) run(cmd *cobra.Command, args []string) {
	readsetIds := strings.Split(p.readsetIds, ",")
	name := strings.Join(args, " ")

	client := NewApiClient()

	fmt.Fprintln(os.Stderr, "Searching for reads...")

	query := &genomics.SearchReadsRequest{
		ReadsetIds: readsetIds,
		PageToken:  p.pageToken,
	}
	if name != "" {
		query.SequenceName = name
		query.SequenceStart = p.sequenceStart
		query.SequenceEnd = p.sequenceEnd
	}

	res, err := client.Reads.Search(query).Do()
	if err != nil {
		log.Fatal(err)
	}

	if len(res.Reads) == 0 {
		fmt.Fprintf(os.Stderr, "No match for %s.\n", name)
		return
	}

	for _, read := range res.Reads {
		printRead(read, 0)
	}

	if res.NextPageToken != "" {
		fmt.Printf("Next Page Token: %v\n", res.NextPageToken)
	}
}

var readsCmd *cobra.Command

func init() {
	readsCmd = &cobra.Command{
		Use:   "reads",
		Short: "Reads functions",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	readsSearch := &paramReadsSearch{}
	readsCmd.AddCommand(readsSearch.command())
}
