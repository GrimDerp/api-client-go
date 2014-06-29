// Copyright 2014 The Google Genomics API Client in Go Authors.
// All rights reserved. Use of this source code is governed by
// an Apache license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"

	genomics "github.com/googlegenomics/api-client-go/v1beta"
	"github.com/spf13/cobra"
)

func printDataset(ds *genomics.Dataset, indentLevel int) {
	indent(indentLevel)
	fmt.Printf("ID: %v, Project ID: %v, Is Public: %v\n",
		ds.Id, ds.ProjectId, ds.IsPublic)
}

type paramDatasetsList struct {
	projectId  int64
	pageToken  string
	maxResults uint64
}

func (p *paramDatasetsList) command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list for datasets in a project",
		Run:   p.run,
	}
	cmd.Flags().Int64VarP(&p.projectId, "project-id", "", 376902546192,
		"Google Cloud console project number")
	cmd.Flags().StringVarP(&p.pageToken, "page-token", "", "",
		"Page token to continue with previous large result")
	cmd.Flags().Uint64VarP(&p.maxResults, "max-results", "", ^uint64(0),
		"Maximum number of results returned by this request.")
	return cmd
}

func (p *paramDatasetsList) run(cmd *cobra.Command, args []string) {
	client := NewApiClient()

	fmt.Fprintf(os.Stderr, "Listing datasets in project %v...\n",
		p.projectId)

	query := client.Datasets.List()
	query.MaxResults(p.maxResults)
	query.ProjectId(p.projectId)
	if p.pageToken != "" {
		query.PageToken(p.pageToken)
	}
	res, err := query.Do()
	if err != nil {
		log.Fatal(err)
	}

	if len(res.Datasets) == 0 {
		fmt.Fprintf(os.Stderr, "No match for %s.\n", p.projectId)
		return
	}

	for _, ds := range res.Datasets {
		printDataset(ds, 0)
	}

	if res.NextPageToken != "" {
		fmt.Printf("Next Page Token: %v\n", res.NextPageToken)
	}
}

type paramDatasetsGet struct {
	datasetId string
}

func (p *paramDatasetsGet) command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get a dataset by ID",
		Run:   p.run,
	}
	cmd.Flags().StringVarP(&p.datasetId, "dataset-id", "", "376902546192",
		"The dataset ID to retrieve.")
	return cmd
}

func (p *paramDatasetsGet) run(cmd *cobra.Command, args []string) {
	client := NewApiClient()

	query := client.Datasets.Get(p.datasetId)
	res, err := query.Do()
	if err != nil {
		log.Fatal(err)
	}

	printDataset(res, 0)
}

var datasetsCmd *cobra.Command

func init() {
	datasetsCmd = &cobra.Command{
		Use:   "datasets",
		Short: "Datasets functions",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	datasetsGet := &paramDatasetsGet{}
	datasetsCmd.AddCommand(datasetsGet.command())

	datasetsList := &paramDatasetsList{}
	datasetsCmd.AddCommand(datasetsList.command())
}
