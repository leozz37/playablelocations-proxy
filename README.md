# Playable Locations Proxy Server

A proxy implementation for the playable locations service converting lat lng bounds to s2 cell id.

## Build and Run

```
docker build . -t proxy
docker run proxy:latest
```
## Usage

### Sample Request via Curl

```sh
curl -X POST \
  localhost:8080/sampleplayablelocations \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/javascript' \
  -H 'x-goog-api-key: YOUR_API_KEY' \
  -d '{
  "criteria": [{
  	"game_object_type": 1234,
  	"fields_to_return": { "paths": ["name", "types", "snapped_point"] },
  	"filter": {
  		"max_location_count": 2,
  		"included_types": ["retail"]
  	}
  }],
  "area_filter_lat_lng_bounds": {
    "northeast": { "lat": 37.761419193645686, "lng": -122.41189956665039 },
    "southwest": { "lat": 37.75714420786591, "lng": -122.41790771484375 }
  }
}
'

```
### JSON Response
```js
{
    "locations_per_game_object_type": {
        "1234": {
            "locations": [
                {
                    "name": "curatedPlayableLocations/ChIJC2juwih-j4ARvLq4f6Oyeuc",
                    "LocationId": {
                        "PlaceId": "ChIJC2juwih-j4ARvLq4f6Oyeuc"
                    },
                    "types": [
                        "retail"
                    ],
                    "center_point": {
                        "latitude": 37.77033960000001,
                        "longitude": -122.4119288
                    },
                    "snapped_point": {
                        "latitude": 37.7703263,
                        "longitude": -122.41195379999999
                    }
                },
                {
                    "name": "curatedPlayableLocations/ChIJkyP7ByZ-j4ARmLwtGpDTgjk",
                    "LocationId": {
                        "PlaceId": "ChIJkyP7ByZ-j4ARmLwtGpDTgjk"
                    },
                    "types": [
                        "retail"
                    ],
                    "center_point": {
                        "latitude": 37.7686504,
                        "longitude": -122.41297919999998
                    },
                    "snapped_point": {
                        "latitude": 37.7686515,
                        "longitude": -122.41295099999999
                    }
                }
            ]
        }
    },
    "ttl": {
        "seconds": 86400
    }
}
```