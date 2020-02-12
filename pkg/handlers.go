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

// Handlers for the server endpoints and conversion method for request to match gRPC SamplePlayableLocationsRequest
package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	playablelocations "google.golang.org/genproto/googleapis/maps/playablelocations/v3"
	sample "google.golang.org/genproto/googleapis/maps/playablelocations/v3/sample"
	"google.golang.org/grpc/metadata"
)

// LatLngSamplePlayableLocationsRequest is similar to the SamplePlayableLocationsRequest
// except for the custom AreaFilterLatLngBounds
type LatLngSamplePlayableLocationsRequest struct {
	playablelocations.SamplePlayableLocationsRequest
	AreaFilterLatLngBounds *LatLngBounds `json:"area_filter_lat_lng_bounds,omitempty"`
}

// ToSamplePlayableLocationsRequest convert request to SamplePlayableLocationsRequest
// by converting the AreaFilterLatLngBounds field to an AreaFilter containing a S2CellId
func (request LatLngSamplePlayableLocationsRequest) ToSamplePlayableLocationsRequest() (*playablelocations.SamplePlayableLocationsRequest, error) {
	if request.AreaFilterLatLngBounds != nil {
		s2CellID, err := request.AreaFilterLatLngBounds.ToS2Cell()
		if err != nil {
			return nil, err
		}
		// instantiate struct if necessary
		if request.AreaFilter == nil {
			request.AreaFilter = &sample.AreaFilter{}
		}

		request.AreaFilter.S2CellId = uint64(s2CellID)
	}

	return &playablelocations.SamplePlayableLocationsRequest{
		AreaFilter: request.AreaFilter,
		Criteria:   request.Criteria,
	}, nil
}

// SamplePlayableLocations proxies the JSON payload to the gRPC endpoint
func SamplePlayableLocations(w http.ResponseWriter, r *http.Request) {
	var (
		request                  LatLngSamplePlayableLocationsRequest
		playableLocationsRequest *playablelocations.SamplePlayableLocationsRequest
		err                      error
		requestBody              []byte
		responseBody             []byte
	)

	requestBody, err = ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		responseBody = []byte(err.Error())
		w.Write(responseBody)
		return
	}

	err = json.Unmarshal(requestBody, &request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		responseBody = []byte(err.Error())
		w.Write(responseBody)
		return
	}

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-goog-api-key", r.Header.Get("x-goog-api-key"))

	playableLocationsRequest, err = request.ToSamplePlayableLocationsRequest()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		responseBody = []byte(err.Error())
		w.Write(responseBody)
		return
	}

	response, err := c.SamplePlayableLocations(ctx, playableLocationsRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		responseBody = []byte(err.Error())
		w.Write(responseBody)
		return
	}

	responseBody, err = json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		responseBody = []byte(err.Error())
		w.Write(responseBody)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}
