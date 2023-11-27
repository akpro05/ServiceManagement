/*Copyright (C) 2020-2030 Supernet Technologies India Private Limited. All rights reserved.
All material appearing on this repository, folder, sub folder and file(s) (“Content/Code etc”) is protected by copyright laws and is the property of Supernet Technologies India Pvt Ltd. You may not copy, reproduce, distribute, publish, display, perform, modify, create derivative works, transmit, or in any way exploit any such content, nor may you distribute any part of this content over any network, including a local area network, sell, offer it for sale, or use such content to construct any kind of database, you may not alter or remove any copyright or other notice from copies of contents on this file.
Copying or storing any contents except as provided by Supernet Technologies India Pvt Ltd is expressly prohibited without prior written permission of the Supernet Technologies.
For permission to use the content, please contact legal@supernet-india.com.
*/
package server

import (
	"ServiceManagement/controller"
	"ServiceManagement/utils/filter"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"ServiceManagement/config/config"
)

func Init() (err error) {
	r := mux.NewRouter()
	r.HandleFunc("/ping", controller.Ping).Methods("GET")

	r.HandleFunc("/getservices", controller.GetServiceApi).Methods("POST")

	//CNPS PGW
	r.HandleFunc("/cnps/payment", controller.CnpsPgwPaymentHandler).Methods("POST")
	r.HandleFunc("/cnps/status", controller.CnpsPgwStatusHandler).Methods("POST")
	//Individual PGW
	r.HandleFunc("/individual/payment", controller.IndividualPgwPaymentHandler).Methods("POST")
	r.HandleFunc("/individual/status", controller.IndividualPgwStatusHandler).Methods("POST")
	//Superpay PGW
	r.HandleFunc("/superpay/payment", controller.SuperpayPgwPaymentHandler).Methods("POST")
	r.HandleFunc("/superpay/status", controller.SuperpayPgwStatusHandler).Methods("POST")

	r.Use(filter.LoggingMiddleware)
	http.Handle("/", r)
	srv := &http.Server{
		Handler: r,
		Addr:    config.CTX.ServerURL,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.WithFields(log.Fields{
		"URL": config.CTX.ServerURL,
	}).Info("http server started....")

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
