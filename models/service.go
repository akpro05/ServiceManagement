/*Copyright (C) 2020-2030 Supernet Technologies India Private Limited. All rights reserved.
All material appearing on this repository, folder, sub folder and file(s) (“Content/Code etc”) is protected by copyright laws and is the property of Supernet Technologies India Pvt Ltd. You may not copy, reproduce, distribute, publish, display, perform, modify, create derivative works, transmit, or in any way exploit any such content, nor may you distribute any part of this content over any network, including a local area network, sell, offer it for sale, or use such content to construct any kind of database, you may not alter or remove any copyright or other notice from copies of contents on this file.
Copying or storing any contents except as provided by Supernet Technologies India Pvt Ltd is expressly prohibited without prior written permission of the Supernet Technologies.
For permission to use the content, please contact legal@supernet-india.com.
*/
package models

import (
	"time"

	"github.com/google/uuid"
	//"gorm.io/gorm"
)

type Producer struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProducerName string
	Email        string
	AccessCode   string
	CreatedBy    uuid.UUID
	CreatedAt    time.Time
	UpdatedBy    uuid.UUID
	UpdatedAt    time.Time
	Status       bool
}

type EsbRequestMetadata struct {
	TransactionId      uuid.UUID `gorm:"type:uuid"`
	CreatedAt          time.Time
	Request_ID         string
	Url                string
	Service            string
	OutRequest         string
	OutResponse        string
	ProducerAccessCode string
	ProducerId         uuid.UUID `gorm:"type:uuid"`
	ConsumerId         uuid.UUID `gorm:"type:uuid"`
}

type Consumer struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
	ConsumerName     string
	Email            string
	ConsumerServices string
	ConsumerCode     string
	ConsumerAddress  string
	AccessCode       string
	Status           bool
	CreatedBy        uuid.UUID
	CreatedAt        time.Time
	UpdatedBy        uuid.UUID
	UpdatedAt        time.Time
}

type ProducerToConsumer struct {
	ID                         uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProducerID                 uuid.UUID `gorm:"type:uuid"`
	ConsumerID                 uuid.UUID `gorm:"type:uuid"`
	ProducerSubscribedServices string
	Status                     bool
	CreatedBy                  uuid.UUID
	CreatedAt                  time.Time
	UpdatedBy                  uuid.UUID
	UpdatedAt                  time.Time
}

func (ProducerToConsumer) TableName() string {
	return "producer_to_consumer"
}
