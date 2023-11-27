/*Copyright (C) 2020-2030 Supernet Technologies India Private Limited. All rights reserved.
All material appearing on this repository, folder, sub folder and file(s) (“Content/Code etc”) is protected by copyright laws and is the property of Supernet Technologies India Pvt Ltd. You may not copy, reproduce, distribute, publish, display, perform, modify, create derivative works, transmit, or in any way exploit any such content, nor may you distribute any part of this content over any network, including a local area network, sell, offer it for sale, or use such content to construct any kind of database, you may not alter or remove any copyright or other notice from copies of contents on this file.
Copying or storing any contents except as provided by Supernet Technologies India Pvt Ltd is expressly prohibited without prior written permission of the Supernet Technologies.
For permission to use the content, please contact legal@supernet-india.com.
*/
package controller

import (
	"ServiceManagement/config/db"
	// "ServiceManagement/entity"

	"ServiceManagement/models"
	// "context"
	"ServiceManagement/config/config"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	// "github.com/google/uuid"
	//"github.com/google/uuid"
	//	"github.com/sirupsen/logrus"
	//	log "github.com/sirupsen/logrus"
	"fmt"
	//	"gorm.io/gorm/clause"
)

type SuperpayStatusRequest struct {
	AccessCode                 string `json:"access_code"`
	SuperesbProducerAccessCode string `json:"superesb_producer_access_code"`
	Channel                    string `json:"channel"`
	TxnNumber                  string `json:"txn_number"`
	Mobile                     string `json:"mobile"`
	RequestID                  string `json:"request_id"`
	Language                   string `json:"language"`
	UserType                   string `json:"user_type"`
	ClientURLUsed              string `json:"client_url_used"`
	EsbTxnId                   string `json:"esb_txn_id"`
}

type SuperpayStatusResponse struct {
	Code             string    `json:"code"`
	Message          string    `json:"message"`
	Timestamp        time.Time `json:"timestamp"`
	Mobile           string    `json:"mobile"`
	RequestID        string    `json:"request_id"`
	TxnNumber        string    `json:"txn_number"`
	PgTxnNumber      string    `json:"pg_txn_number"`
	BillerTxnNumber  string    `json:"biller_txn_number"`
	Operator         string    `json:"operator"`
	AccNumber        string    `json:"acc_number"`
	Amount           string    `json:"amount"`
	Status           string    `json:"status"`
	ServiceName      string    `json:"service_name"`
	UserType         string    `json:"user_type"`
	CustomerMetadata string    `json:"customer_metadata"`
}

func SuperpayPgwStatusHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("---------------- Request at SuperpayPgwStatusHandler --------------")

	var resp SuperpayStatusResponse
	var paybillresponse SuperpayStatusResponse

	var getservicerequest SuperpayStatusRequest
	json.NewDecoder(r.Body).Decode(&getservicerequest)
	request, err := json.Marshal(&getservicerequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("SuperpayStatusHandler request", string(request))

	var producerID string
	language := getservicerequest.Language
	// Perform the GORM query
	result := db.Db.Model(&models.Producer{}).Select("id").Where("access_code = ?", getservicerequest.SuperesbProducerAccessCode).Take(&producerID)

	if result.Error != nil {
		// Handle the error if the query fails
		fmt.Println("Error while fetching producer ID:", result.Error)
		// Handle the error appropriately (e.g., return an error response)
		return
	}

	// The subscriberID variable now contains the ID of the subscriber with the specified access code
	fmt.Println("Producer ID:", producerID)

	var count int64
	result2 := db.Db.Model(&models.ProducerToConsumer{}).
		Where("producer_id = ?", producerID).
		Count(&count)

	if result2.Error != nil {
		// Handle the error if the query fails
		fmt.Println("Error while counting records:", result2.Error)
		// Handle the error appropriately (e.g., return an error response)
		return
	}

	if count <= 0 {
		// No records found for the given producerID
		fmt.Println("No records found for producerID:", producerID)
		resp.Code = "0001"
		resp.Status = "Declined"
		if language == "fr" {
			resp.Message = config.L.FR_NoMappedConsumer + getservicerequest.SuperesbProducerAccessCode
		} else {
			resp.Message = config.L.EN_NoMappedConsumer + getservicerequest.SuperesbProducerAccessCode
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	client_url_used := getservicerequest.ClientURLUsed
	esb_txn_id := getservicerequest.EsbTxnId
	fmt.Println("esb_txn_id  ", esb_txn_id)

	// SQL query to fetch consumer_id and producer_subscribed_services
	rows, err := db.Db.Raw("SELECT consumer_id, producer_subscribed_services FROM public.producer_to_consumer WHERE producer_id = ? order by created_at asc", producerID).Rows()
	if err != nil {
		// Handle the error if the query fails
		fmt.Println("Error executing SQL query:", err)
		return
	}
	defer rows.Close()

	found := false // Flag to track if the URL is found
	var consumerID string
	var s_name string
	for rows.Next() {

		if found {
			// If a match is already found, break out of the outer loop
			break
		}

		var subscribedServices string

		if err := rows.Scan(&consumerID, &subscribedServices); err != nil {
			// Handle the error if scanning fails
			fmt.Println("Error scanning row:", err)
			continue
		}

		// Unmarshal the JSON data
		var serviceData map[string][]struct {
			ServiceName string `json:"service_name"`
			ServiceURL  string `json:"service_url"`
		}
		if err := json.Unmarshal([]byte(subscribedServices), &serviceData); err != nil {
			// Handle the error if JSON unmarshaling fails
			fmt.Println("Error unmarshaling producer_subscribed_services:", err)
			continue
		}

		// Check if the given URL is present in the services
		for _, services := range serviceData["subsrcibed_services"] {
			if services.ServiceURL == client_url_used {
				found = true
				fmt.Printf("Consumer ID: %s\n", consumerID)
				fmt.Printf("ServiceName : %s\n", services.ServiceName)
				s_name = services.ServiceName
				break

			}
		}
	}

	if !found {
		resp.Code = "0001"
		resp.Status = "Declined"
		if language == "fr" {
			resp.Message = config.L.FR_ReqUrlNotFound
		} else {
			resp.Message = config.L.EN_ReqUrlNotFound
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	fmt.Println("Id of the consumer mapped is  ", consumerID)

	//Code to update esb request details
	res := db.Db.Model(&models.EsbRequestMetadata{}).
		Where("transaction_id = ?", esb_txn_id).
		Updates(map[string]interface{}{"producer_id": producerID, "consumer_id": consumerID, "service": s_name}).Error
	fmt.Println("Update Transactions error", res)

	var status bool
	status_result := db.Db.Model(&models.ProducerToConsumer{}).
		Where("producer_id = ? AND consumer_id = ?", producerID, consumerID).
		Pluck("status", &status)

	if status_result.Error != nil {
		log.Println("Error:", result.Error)
	} else {
		//log.Println("Status of the mapping between Producer and Consumer:", status)
	}

	if status == false {
		resp.Code = "0001"
		resp.Status = "Declined"
		if language == "fr" {
			resp.Message = config.L.FR_P2C_Status
		} else {
			resp.Message = config.L.EN_P2C_Status
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	var consumer models.Consumer
	db.Db.Model(&models.Consumer{}).Where("id = ?", consumerID).Select("consumer_address,status").First(&consumer)

	fmt.Println("Status of the consumer mapped is  ", consumer.Status)

	if consumer.Status == false {

		resp.Code = "0001"
		resp.Status = "Declined"
		if language == "fr" {
			resp.Message = config.L.FR_Consumer_Status
		} else {
			resp.Message = config.L.EN_Consumer_Status
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return

	}

	fmt.Println("Domain/Host Address of the consumer mapped is  ", consumer.ConsumerAddress)

	// Now from here request will be sent to backend consumer service
	consumer_url := consumer.ConsumerAddress + client_url_used
	fmt.Println("consumer_url is :- ", consumer_url)

	consumer_request, err := http.NewRequest("POST", consumer_url, strings.NewReader(string(request)))
	consumer_request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	consumer_resp, err := client.Do(consumer_request)
	if err != nil {
		panic(err)
	}
	defer consumer_resp.Body.Close()

	// fmt.Println("response Status:", consumer_resp.Status)
	// fmt.Println("response Headers:", consumer_resp.Header)

	json.NewDecoder(consumer_resp.Body).Decode(&paybillresponse)
	sresp, err := json.Marshal(&paybillresponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("--------------------------------------------------------")
	fmt.Println("Consumer Backend Service Response :- ", string(sresp))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paybillresponse)
	return

}
