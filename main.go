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
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"code.google.com/p/goauth2/oauth"

	genomics "github.com/googlegenomics/api-client-go/v1beta"
	"github.com/spf13/cobra"
)

func obtainOauthCode(url string) string {
	fmt.Println("Please visit the below URL to obtain OAuth2 code.")
	fmt.Println()
	fmt.Println(url)
	fmt.Println()
	fmt.Println("Please enter the code here:")

	line, _, _ := bufio.NewReader(os.Stdin).ReadLine()

	return string(line)
}

func prepareClient() (*http.Client, error) {
	if oauthJsonFile == "" {
		return &http.Client{}, nil
	}

	jsonData, err := ioutil.ReadFile(oauthJsonFile)
	if err != nil {
		return nil, err
	}

	var data struct {
		Installed struct {
			Client_Id     string
			Client_Secret string
			Redirect_Uris []string
			Auth_Uri      string
			Token_Uri     string
		}
	}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}

	config := &oauth.Config{
		ClientId:     data.Installed.Client_Id,
		ClientSecret: data.Installed.Client_Secret,
		RedirectURL:  data.Installed.Redirect_Uris[0],
		Scope: strings.Join([]string{
			genomics.GenomicsScope,
			genomics.DevstorageRead_writeScope,
		}, " "),
		AuthURL:    data.Installed.Auth_Uri,
		TokenURL:   data.Installed.Token_Uri,
		TokenCache: oauth.CacheFile(".oauth2_cache.json"),
	}

	transport := &oauth.Transport{Config: config}
	token, err := config.TokenCache.Token()
	if err != nil {
		url := config.AuthCodeURL("")
		code := obtainOauthCode(url)
		token, err = transport.Exchange(code)
		if err != nil {
			return nil, err
		}
	}

	transport.Token = token
	client := transport.Client()

	return client, nil
}

func NewApiClient() *genomics.Service {
	client, err := prepareClient()
	if err != nil {
		log.Fatal(err)
	}
	baseApi, err := genomics.New(client)
	if err != nil {
		log.Fatal(err)
	}
	return baseApi
}

var oauthJsonFile string

func main() {
	mainCmd := &cobra.Command{
		Use:   os.Args[0],
		Short: "Google Genomics API Go client",
	}
	mainCmd.PersistentFlags().StringVarP(&oauthJsonFile, "use-oauth", "", "",
		"Path to client_secret.json")
	mainCmd.AddCommand(readsCmd)
	mainCmd.AddCommand(readsetsCmd)
	mainCmd.Execute()
	os.Exit(0)
}
