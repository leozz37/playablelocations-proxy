// Copyright 2020 Google Inc. All Rights Reserved.
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
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net/http"

	playablelocations "google.golang.org/genproto/googleapis/maps/playablelocations/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	c playablelocations.PlayableLocationsClient
)

func main() {
	var (
		port int
	)
	flag.IntVar(&port, "port", 8080, "Port: defaults to 8080")
	flag.Parse()

	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("did not get cert pool: %s", err)
	}

	creds := credentials.NewClientTLSFromCert(pool, "")
	if err != nil {
		log.Fatalf("did not get service account: %s", err)
	}

	conn, err := grpc.Dial(
		"playablelocations.googleapis.com:443",
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	defer conn.Close()

	c = playablelocations.NewPlayableLocationsClient(conn)

	http.HandleFunc("/sampleplayablelocations", SamplePlayableLocations)

	fmt.Printf("Playable Locations Proxy Server Running on port :%d\n", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("failed to start server: %s", err)
	}

}
