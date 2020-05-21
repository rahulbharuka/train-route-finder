# Train Route Finder
This project creates a route finder service for an urban rail network. It returns shortest route between passed source and destination.

### Features
- Supports `simple route` i.e. least number of stops to the destination. The returned routes are ranked in decreasing order of stop count.
- Also support `realtime route` i.e. shortest time route based on travel time/day and number of rail line interchange. The returned routes are ranked in decreasing order of travel time.
- Returns configurable number of `Top K` shortest routes by using Yen's algorithm (https://en.wikipedia.org/wiki/Yen%27s_algorithm).
---

### How to run ?
* Download the complete source code.
* cd to the directory and build the Golang binary for the platform (darwin, linux, windows etc) and architecture (amd64, 386 etc)
```
    env GOOS=<plafrom> GOARCH=<architecture> go build -o <binary-name>

    e.g. for amd64 linux
    env GOOS=linux GOARCH=amd64 go build -o train-route-finder
```
* Export following environment variables.
```
    export PORT=<port-number>
    export STATION_MAP_FILE=<station-map-file-path>
    export TRAINLINE_COST_FILE=<trainline-cost-file-path>
    export INTERCHANGE_COST_FILE=<interchange-cost-file-path>
    export MAX_ROUTES=<max-routes-to-return>

    e.g.
    export PORT=8080
    export STATION_MAP_FILE=./StationMap.csv
    export TRAINLINE_COST_FILE=./trainline_cost.csv
    export INTERCHANGE_COST_FILE=./interchange_cost.csv
    export MAX_ROUTES=3
```
* Now run
```    
    ./train-route-finder
```
---

### Assumptions
- All travel time cost are postive integers.
- No two consecutive stations share more than one rail line.
- Rail network data is provided in CSV file with format <stationCode,station-name,date-of-opening>
    * First two character of **_stationCode_** are used to determine train line.
    * **_stationCode_** is used to determine order of stations on a train line.

---

### APIs
`GET /routes`
  * Usage: To get route(s) from source to destination.

```
    Query parameters:
    src - source station name (required)
    dst - destination station name (required)
    journeyTime - expected start time of journey in YYYY-MM-DDTHH:MM format (optional)

    HTTP Response:
    200 - if one are more routes are found
    400 - if request format is not correct
    404 - if no route exist between source and destination
    500 - if unknown error occured while finding route(s).
``` 
  * Sample `Simple route` request/response:
```
Request:
        curl --location --request GET 'http://localhost:8080/routes?src=Holland%20Village&dst=Bugis'
Response:
        [
            {
                "heading": "Number of stops to destination: 7",
                "steps": "Take CC line from Holland Village to Botanic Gardens. Change from CC line to DT line. Take DT line from Botanic Gardens to Bugis."
            },
            {
                "heading": "Number of stops to destination: 8",
                "steps": "Take CC line from Holland Village to Botanic Gardens. Change from CC line to DT line. Take DT line from Botanic Gardens to Little India. Change from DT line to NE line. Take NE line from Little India to Dhoby Ghaut. Change from NE line to NS line. Take NS line from Dhoby Ghaut to City Hall. Change from NS line to EW line. Take EW line from City Hall to Bugis."
            },
            {
                "heading": "Number of stops to destination: 9",
                "steps": "Take CC line from Holland Village to Caldecott. Change from CC line to TE line. Take TE line from Caldecott to Stevens. Change from TE line to DT line. Take DT line from Stevens to Bugis."
            }
        ]
```

  * Sample `Realtime route` equest/response:
```
Request:
        curl --location --request GET 'http://localhost:8080/routes?src=Boon%20Lay&dst=Little%20India&journeyTime=2019-01-31T19:00'
Response:
        [
            {
                "heading": "Expected Travel time: 150",
                "steps": "Take EW line from Boon Lay to Buona Vista. Change from EW line to CC line. Take CC line from Buona Vista to Botanic Gardens. Change from CC line to DT line. Take DT line from Botanic Gardens to Little India."
            },
            {
                "heading": "Expected Travel time: 173",
                "steps": "Take EW line from Boon Lay to Outram Park. Change from EW line to NE line. Take NE line from Outram Park to Little India."
            },
            {
                "heading": "Expected Travel time: 185",
                "steps": "Take EW line from Boon Lay to Buona Vista. Change from EW line to CC line. Take CC line from Buona Vista to Caldecott. Change from CC line to TE line. Take TE line from Caldecott to Stevens. Change from TE line to DT line. Take DT line from Stevens to Little India."
            }
        ]
```

---

### External Dependencies
- Uses `Gin` web framework for routing.
---

### Testing
- Added unit tests for repository layer.
---

### References
This project uses following implementation of Yen's algorithm as a base and modifies it for current use case.
https://github.com/starwander/goraph 
