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

// Helper methods and structs for working with S2 Geometry
package main

import (
	"errors"

	"github.com/golang/geo/s2"
)

// LatLng struct
type LatLng struct {
	Lat float64 `json:"latitude"`
	Lng float64 `json:"longitude"`
}

// LatLngBounds contains two LatLngs
type LatLngBounds struct {
	Northeast LatLng `json:"northeast"`
	Southwest LatLng `json:"southwest"`
}

// ToS2Cell Converts LatLng to a single S2 Cell and raises and error if more than
// a single cell is found
func (bounds LatLngBounds) ToS2Cell() (s2.CellID, error) {
	// This is an approximation. No guarantee to have a single cell
	rc := &s2.RegionCoverer{MinLevel: 11, MaxLevel: 14, MaxCells: 1}

	r := s2.EmptyRect()
	r = r.AddPoint(s2.LatLngFromDegrees(bounds.Northeast.Lat, bounds.Northeast.Lng))
	r = r.AddPoint(s2.LatLngFromDegrees(bounds.Southwest.Lat, bounds.Southwest.Lng))

	cells := rc.FastCovering(r)
	if len(cells) > 1 {
		return 0, errors.New("area_filter_lat_lng_bounds spans multiple S2 cells")
	}

	return cells[0], nil
}
