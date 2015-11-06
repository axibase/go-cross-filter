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
** Created by Gregory Kutuzov on 19/10/15.
*
 */

package application

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"net"
)

type Web struct {
	router *mux.Router
	server *http.Server

	portListener net.Listener
	isListening bool
}

func NewWeb() *Web {
	r := mux.NewRouter()
	server := &http.Server{Handler: r}
	return &Web{router: r, server: server}
}

func (self *Web) Serve(port uint) error {
	if self.isListening {
		panic("Error server is already listenning on "+ self.server.Addr)
	}
	ln, err := net.Listen("tcp", ":"+strconv.FormatUint(uint64(port), 10))
	if err!=nil {
		return err
	}
	self.portListener = ln
	self.server.Serve(self.portListener)

	return nil
}
func (self *Web) Stop() error {
	return self.portListener.Close()
}
func (self *Web) ResetHandlers() {
	self.router = mux.NewRouter()
	self.server.Handler = self.router
}
func (self *Web) IsListening() bool {
	return self.isListening
}

func (self *Web) Register(path string, controller Controller) {
	controller.SetRouter(self.router.PathPrefix(path).Subrouter())
}

type Controller interface {
	SetRouter(r *mux.Router)
}
