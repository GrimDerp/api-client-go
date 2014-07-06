// Copyright 2014 The Google Genomics API Client in Go Authors.
// All rights reserved. Use of this source code is governed by
// an Apache license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/goauth2/oauth/jwt"

	genomics "code.google.com/p/google-api-go-client/genomics/v1beta"
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

type oauthConfig struct {
	Web struct {
		Client_Email string
		Client_Id    string
		Token_Uri    string
	}
	Installed struct {
		Client_Id     string
		Client_Secret string
		Redirect_Uris []string
		Auth_Uri      string
		Token_Uri     string
	}
}

func prepareNativeClient(data *oauthConfig) (*http.Client, error) {
	config := &oauth.Config{
		ClientId:     data.Installed.Client_Id,
		ClientSecret: data.Installed.Client_Secret,
		RedirectURL:  data.Installed.Redirect_Uris[0],
		Scope:        oauthScope,
		AuthURL:      data.Installed.Auth_Uri,
		TokenURL:     data.Installed.Token_Uri,
		TokenCache:   oauth.CacheFile(".oauth2_cache.json"),
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

func prepareServiceClient(data *oauthConfig, keyBytes []byte) (
	*http.Client, error) {
	// Craft the ClaimSet and JWT token.
	token := jwt.NewToken(data.Web.Client_Email, oauthScope, keyBytes)
	token.ClaimSet.Aud = data.Web.Token_Uri

	transport, err := jwt.NewTransport(token)
	if err != nil {
		log.Fatal("Assertion error:", err)
	}

	return transport.Client(), nil
}

func prepareClient() (*http.Client, error) {
	if oauthJsonFile == "" {
		return &http.Client{}, nil
	}

	jsonData, err := ioutil.ReadFile(oauthJsonFile)
	if err != nil {
		return nil, err
	}

	var data oauthConfig
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}

	if data.Web.Client_Id != "" {
		fullname := filepath.Base(oauthJsonFile)
		ext := filepath.Ext(fullname)
		pemFile := fullname[:len(fullname)-len(ext)] + ".pem"
		// Read the pem file bytes for the private key.
		keyBytes, err := ioutil.ReadFile(pemFile)
		if err != nil {
			log.Fatal("Error reading private key file:", err)
		}

		return prepareServiceClient(&data, keyBytes)
	} else if data.Installed.Client_Id != "" {
		return prepareNativeClient(&data)
	}

	return nil, fmt.Errorf("Not recognized OAuth data")
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

var (
	oauthJsonFile string
	oauthScope    = strings.Join([]string{
		genomics.GenomicsScope,
		genomics.DevstorageRead_writeScope,
	}, " ")
)

func main() {
	mainCmd := &cobra.Command{
		Use:   os.Args[0],
		Short: "Google Genomics API Go client",
	}
	mainCmd.PersistentFlags().StringVarP(&oauthJsonFile, "use-oauth", "", "",
		"Path to client_secret.json")
	mainCmd.AddCommand(datasetsCmd)
	mainCmd.AddCommand(readsCmd)
	mainCmd.AddCommand(readsetsCmd)
	mainCmd.Execute()
	os.Exit(0)
}
