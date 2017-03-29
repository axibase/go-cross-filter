/**
* Copyright 2015 Axibase Corporation or its affiliates. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License").
* You may not use this file except in compliance with the License.
* A copy of the License is located at
*
* https://www.axibase.com/atsd/axibase-apache-2.0.pdf
*
* or in the "license" file accompanying this file. This file is distributed
* on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
* express or implied. See the License for the specific language governing
* permissions and limitations under the License.
*
** Created by Gregory Kutuzov on 26/10/15.
*
 */

package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	neturl "net/url"
	"strings"

	"github.com/gorilla/mux"
)

type PortalController struct {
	atsdProxy *httputil.ReverseProxy
}

func NewPortalController(url neturl.URL) *PortalController {
	return &PortalController{
		atsdProxy: httputil.NewSingleHostReverseProxy(&url),
	}
}

func (self *PortalController) SetRouter(r *mux.Router) {
	r.HandleFunc("/api/v1/series", self.atsdProxyHandler)
	r.HandleFunc("/charts", self.portalHandler)
	r.HandleFunc("/chartsConfig", self.chartsConfigHandler)
}

func (self *PortalController) atsdProxyHandler(w http.ResponseWriter, r *http.Request) {
	if isValidRequest(r) {
		self.atsdProxy.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func isValidRequest(r *http.Request) bool {
	query := r.URL.Query()
	tableName := query.Get("table")
	table, ok := Instance().TableService.GetTable(tableName)
	if !ok {
		return false
	}
	entityColName := getEntityColName()
	if entityColName != "" {
		cols, err := getCols(r.URL.RawQuery)
		if err != nil {
			return false
		}
		entity := ""
		for _, col := range cols {
			if col.Name == entityColName {
				entity = col.Value
			}
		}
		for j, col := range table.Columns {
			if col.Label == entityColName {
				for _, row := range table.Rows {
					if row[j] == entity {
						return true
					}
				}
			}
		}
	}
	return false
}
func getEntityColName() string {
	entityColname := ""
	for _, tableConfig := range Instance().config.TableConfigs {
		for _, colConfig := range tableConfig.ColumnsConfig {
			colName := colConfig["name"].(string)
			isEntity, ok := colConfig["entity"].(bool)
			if ok && isEntity {
				entityColname = colName
			}
		}
	}
	return entityColname
}

func (self *PortalController) chartsConfigHandler(w http.ResponseWriter, r *http.Request) {
	if isValidRequest(r) {
		query := r.URL.Query()
		table := query.Get("table")
		cols, err := getCols(r.URL.RawQuery)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		ConfigView(table, cols, w)
	} else {
		http.NotFound(w, r)
	}

}

func getCols(rawQuery string) ([]*Col, error) {
	queries := strings.Split(rawQuery, "&")
	cols := []*Col{}
	for _, q := range queries {
		if strings.HasPrefix(q, "col=") {
			q, err := neturl.QueryUnescape(q)
			if err != nil {
				return nil, err
			}
			q = strings.Replace(q, "\\:", ":", -1)
			namevalue := strings.Split(strings.TrimPrefix(q, "col="), ":")
			name := namevalue[0]
			value := namevalue[1]
			cols = append(cols, &Col{Name: name, Value: value})
		}
	}
	return cols, nil
}

func (self *PortalController) portalHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(isValidRequest(r))
	if isValidRequest(r) {
		query := r.URL.Query()
		if len(query.Encode()) == 0 {
			http.NotFound(w, r)
			return
		}
		PortalView(query.Encode(), w)
	} else {
		http.NotFound(w, r)
	}
}

type TableController struct {
}

func (self *TableController) SetRouter(r *mux.Router) {
	r.Path("/").HandlerFunc(self.indexHandler)
	r.Path("/table").Queries("table", "{table}", "entityGroup", "{entityGroup}").HandlerFunc(self.tableHandler)
	r.Path("/tablesList").HandlerFunc(self.tableListHandler)
	r.Path("/entityGroupsList").Queries("table", "{table}").HandlerFunc(self.entityGroupsListHandler)
}

func (self *TableController) tableHandler(w http.ResponseWriter, r *http.Request) {
	reqTableName := r.URL.Query().Get("table")
	reqEntityGroup := r.URL.Query().Get("entityGroup")
	TableView(reqTableName, reqEntityGroup, w)
}

func (self *TableController) entityGroupsListHandler(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	reqTableName := v.Get("table")
	EntityGroupListView(reqTableName, w)
}

func (self *TableController) tableListHandler(w http.ResponseWriter, r *http.Request) {
	TableListView(w)
}

func (self *TableController) indexHandler(w http.ResponseWriter, r *http.Request) {
	IndexView(w)
}
