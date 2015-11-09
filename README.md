#Go Cross-filter

This is a golang application that implements the [cross-filter concept](https://square.github.io/crossfilter/) on top of data stored in [Axibase Time Series Database](https://axibase.com/products/axibase-time-series-database/). It provides a capability to apply graphical filters to entity tags and time series retrieved from the database, without reloading the dataset on the client. Datasets are loaded from ATSD using [SQL queries](https://axibase.com/atsd/api/#sql) and refreshed on schedule. [Cross-filter.js](https://square.github.io/crossfilter/) and [datatable](https://www.datatables.net/) are used to build and maintain indices for fast filterting on the client. We use pie and histogram charts from the ATSD widget library to display the filters.

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
          "filter": true            #filter flag to use column for a visual filtering.
        },
        {
          "name": "loc_area",
          "filter": true
        },
        {
          "name": "app",
          "filter": true
        },
        {
          "name": "FS id",
          "filter": true
        },
        {
          "name": "Cpu Busy, %",
          "filter": true,
          "formatter": {            #specify formatter for numeric columns to render column value based on formula (round(multiplier * value)) 
            "round": 1              #you can specify "round" and "multiplier" independently. 
          }                         
        },
        {
          "name": "Memory Used, %",
          "filter": true,
          "formatter": {
            "round": 1
          }
        },
        {
          "name": "Memory Free, Mb",
          "filter": true,
          "formatter": {
            "round": 0,
            "multiplier": 0.00097656  #converting bytes to Mb
          }
        },
        {
          "name": "FS used, %",
          "filter": true,
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

