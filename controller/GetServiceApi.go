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
	// "database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	// "strings"
	"log"
	//"github.com/google/uuid"
	//	"github.com/sirupsen/logrus"
	//	log "github.com/sirupsen/logrus"
	"fmt"
	//	"gorm.io/gorm/clause"
)

type GetServiceRequest struct {
	SuperesbProducerAccessCode string `json:"superesb_producer_access_code"`
	RequestID                  string `json:"request_id"`
	Channel                    string `json:"channel"`
	Language                   string `json:"language"`
	EsbTxnId                   string `json:"esb_txn_id"`
}

type SubscribedService struct {
	ServiceName string `json:"service_name"`
	ServiceURL  string `json:"service_url"`
}

type MappedService struct {
	ConsumerName       string              `json:"consumer_name"`
	ConsumerCode       string              `json:"consumer_code"`
	SubscribedServices []SubscribedService `json:"subscribed_services"`
}

type GetServiceResp struct {
	Code               string          `json:"code"`
	Message            string          `json:"message"`
	MappedServicesList []MappedService `json:"mapped_services_list"`
}

func GetServiceApi(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------------- Request at GetServiceApi --------------")

	var resp GetServiceResp
	// var paybillresponse PayResponse

	var getservicerequest GetServiceRequest
	json.NewDecoder(r.Body).Decode(&getservicerequest)
	request, err := json.Marshal(&getservicerequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Getservice request", string(request))

	var producerID string

	// Perform the GORM query
	result := db.Db.Model(&models.Producer{}).Select("id").Where("access_code = ?", getservicerequest.SuperesbProducerAccessCode).Take(&producerID)

	if result.Error != nil {
		// Handle the error if the query fails
		fmt.Println("Error while fetching producer ID:", result.Error)
		// Handle the error appropriately (e.g., return an error response)
		return
	}

	// The subscriberID variable now contains the ID of the subscriber with the specified access code
	formated_producerID, _ := uuid.Parse(producerID)
	esb_txn_id := getservicerequest.EsbTxnId
	fmt.Println("Producer ID:", formated_producerID)
	fmt.Println("EsbTxnId:", esb_txn_id)

	res := db.Db.Model(&models.EsbRequestMetadata{}).
		Where("transaction_id = ?", esb_txn_id).
		Updates(map[string]interface{}{"producer_id": formated_producerID}).Error
	fmt.Println("Update Transactions error", res)

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
		resp.Message = "No Consumer Mapping exists for given Producer :- " + getservicerequest.SuperesbProducerAccessCode
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	var results []models.ProducerToConsumer

	db.Db.Table("producer_to_consumer").
		Select("consumer_id, producer_subscribed_services").
		Where("producer_id = ?", producerID).
		Find(&results)

	// Define a struct to hold the consumer name and code.
	type ConsumerData struct {
		ConsumerName string
		ConsumerCode string
	}

	// Create a map to store consumer data using their IDs as keys.
	consumerDataMap := make(map[uuid.UUID]ConsumerData)

	for _, result := range results {
		var consumerData ConsumerData
		// Fetch the consumer name and code based on the consumer ID.
		if err := db.Db.Model(&models.Consumer{}).Where("id = ?", result.ConsumerID).Select("consumer_name, consumer_code").Scan(&consumerData).Error; err != nil {
			log.Printf("Error fetching consumer data: %v", err)
			// Handle the error appropriately.
		} else {
			consumerDataMap[result.ConsumerID] = consumerData
		}
	}

	// Now, you have a map containing the consumer data indexed by their IDs.
	// You can access the data as needed.
	mappedServicesList := make([]MappedService, 0)
	for _, result := range results {
		consumerData := consumerDataMap[result.ConsumerID]
		log.Printf("Consumer ID: %s, Consumer Name: %s, Consumer Code: %s,", result.ConsumerID, consumerData.ConsumerName, consumerData.ConsumerCode)

		var producerServices ProducerServices
		if err := json.Unmarshal([]byte(result.ProducerSubscribedServices), &producerServices); err != nil {
			log.Printf("Error unmarshalling JSON: %v", err)
			continue
		}

		subscribedServices := make([]SubscribedService, len(producerServices.SubsrcibedServices))
		for i, service := range producerServices.SubsrcibedServices {
			subscribedServices[i] = SubscribedService{
				ServiceName: service.ServiceName,
				ServiceURL:  service.ServiceURL,
			}
		}

		mappedService := MappedService{
			ConsumerName:       consumerData.ConsumerName,
			ConsumerCode:       consumerData.ConsumerCode,
			SubscribedServices: subscribedServices,
		}
		mappedServicesList = append(mappedServicesList, mappedService)
	}

	mappedDetails := GetServiceResp{
		Code:               "0000",
		Message:            "Producer's mapped Services Fetched",
		MappedServicesList: mappedServicesList,
	}

	w.Header().Set("Content-Type", "application/json")

	// Marshal the mappedDetails into a JSON string
	responseJSON, err := json.Marshal(mappedDetails)
	if err != nil {
		log.Printf("Error marshaling response to JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Log the JSON response
	log.Printf("Response JSON: %s", responseJSON)

	// Write the JSON response to the client
	w.Write(responseJSON)
}
