/*Copyright (C) 2020-2030 Supernet Technologies India Private Limited. All rights reserved.
All material appearing on this repository, folder, sub folder and file(s) (“Content/Code etc”) is protected by copyright laws and is the property of Supernet Technologies India Pvt Ltd. You may not copy, reproduce, distribute, publish, display, perform, modify, create derivative works, transmit, or in any way exploit any such content, nor may you distribute any part of this content over any network, including a local area network, sell, offer it for sale, or use such content to construct any kind of database, you may not alter or remove any copyright or other notice from copies of contents on this file.
Copying or storing any contents except as provided by Supernet Technologies India Pvt Ltd is expressly prohibited without prior written permission of the Supernet Technologies.
For permission to use the content, please contact legal@supernet-india.com.
 */
package filter

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.WithFields(log.Fields{
			"body":    r.Body,
			"headers": r.Header,
			"url":     r.RequestURI}).Trace("Request")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
