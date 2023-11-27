/*Copyright (C) 2020-2030 Supernet Technologies India Private Limited. All rights reserved.
All material appearing on this repository, folder, sub folder and file(s) (“Content/Code etc”) is protected by copyright laws and is the property of Supernet Technologies India Pvt Ltd. You may not copy, reproduce, distribute, publish, display, perform, modify, create derivative works, transmit, or in any way exploit any such content, nor may you distribute any part of this content over any network, including a local area network, sell, offer it for sale, or use such content to construct any kind of database, you may not alter or remove any copyright or other notice from copies of contents on this file.
Copying or storing any contents except as provided by Supernet Technologies India Pvt Ltd is expressly prohibited without prior written permission of the Supernet Technologies.
For permission to use the content, please contact legal@supernet-india.com.
*/
package db

import (
	// "ServiceManagement/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"ServiceManagement/config/config"
)

var Db *gorm.DB

func Init() (err error) {
	dsn := config.CTX.DBConnection
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Migrate the schema

	return err
}
