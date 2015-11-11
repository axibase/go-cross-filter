#Go Cross-filter

This is a golang application that implements the [cross-filter concept](https://square.github.io/crossfilter/) on top of data stored in [Axibase Time Series Database](https://axibase.com/products/axibase-time-series-database/). It provides a capability to apply graphical filters to entity tags and time series retrieved from the database, without reloading the dataset on the client. Datasets are loaded from ATSD using [SQL queries](https://axibase.com/atsd/api/#sql) and refreshed on schedule. [Cross-filter.js](https://square.github.io/crossfilter/) and [datatables](https://www.datatables.net/) are used to build and maintain indices for fast filterting on the client. We use pie and histogram charts from the ATSD widget library to display the filters.

Try it [here](http://apps.axibase.com/cross-filter).
## Dependencies

Client:

1. bootstrap v3.3.5
2. jquery v2.1.4
3. d3 v3.5.6
4. datatables v1.10.9
5. datatables scroller v1.3.0
6. crossfilter v1.3.12

Server:

1. ATSD rev. >=11130
2. Go dependencies from Godeps.json  


## Getting Started

### Requirements

1. golang build tools
2. godep

### Install

Navigate to your destination folder and execute the following commands:

Load last version from github repository:
```bash
# Load last version from github repository
go get github.com/axibase/go-capacity-screener

# Compile the project:
godep go build
```

### Usage

This is a default program structure.

```text
.
|-- config.json                   - Capacity Screener configuration file
|-- go-capacity-screener          - Capacity Screener binary
|-- templates                     - Web page templates
|   |-- index.tpl
|   `-- view.tpl
`-- portals                       - User defined portals folder
    |-- itcam-rtt.conf        
    |-- nmon.conf
    `-- oracle_databases.conf
```

#### Edit configuration file

Define ATSD url, username and password:
```bash
{
  "url": "http://atsd_server:8088",   #define ATSD access credentials
  "user": "username",                 
  "password": "password",             
  "port": 8000,                       #application web view port
  "updatePeriod": "1m",               #interval used to update all defined tables 
  "tables": [
  ...                                 #define tables here
  ]                        
}
```

Define tables:
```bash
...
"tables": [
	{
      "name": "Linux Performance",                                 #table name                               
      "entityGroups": ["nmon-linux", "nurswg-dc1", "nurswg-dc2"],  #which entity groups can be used to filter the table(dataset) 
      "sqlQuery": ATSD_SQL_QUERY,                                  #SQL query to load data from ATSD
      "portalConfigPath": "portals/linux_performance.config",      #portal configuration file
      "columns": [                  #all column names should match aliases in sqlQuery.
        {
          "name": "entity",                  
          "entity": true            #entity column flag (only one per table) to view the portal
        },
        {
          "name": "os",             
          "filter": {}            #filter object to use column for a visual filtering.
        },
        {
          "name": "loc_area",
          "filter": {}
        },
        {
          "name": "app",
          "filter": {}
        },
        {
          "name": "FS id",
          "filter": {}
        },
        {
          "name": "Cpu Busy, %",
          "filter": {
            range: [0, 100]         #range option to set histogram range for numeric columns
          },
          "formatter": {            #specify formatter for numeric columns to render column value based on formula (round(multiplier * value)) 
            "round": 1              #you can specify "round" and "multiplier" independently. 
          }                         
        },
        {
          "name": "Memory Used, %",
          "filter": {
            range: [0, 100]
          },
          "formatter": {
            "round": 1
          }
        },
        {
          "name": "Memory Free, Mb",
          "filter": {},
          "formatter": {
            "round": 0,
            "multiplier": 0.00097656  #converting bytes to Mb
          }
        },
        {
          "name": "FS used, %",
          "filter": {
            range: [0, 100]
          },
          "formatter": {
            "round": 1
          }
        }
      ]
	},
	...
]
```

#### Start the Application

```bash
./go-cross-filter
```

