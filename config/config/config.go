/*Copyright (C) 2020-2030 Supernet Technologies India Private Limited. All rights reserved.
All material appearing on this repository, folder, sub folder and file(s) (“Content/Code etc”) is protected by copyright laws and is the property of Supernet Technologies India Pvt Ltd. You may not copy, reproduce, distribute, publish, display, perform, modify, create derivative works, transmit, or in any way exploit any such content, nor may you distribute any part of this content over any network, including a local area network, sell, offer it for sale, or use such content to construct any kind of database, you may not alter or remove any copyright or other notice from copies of contents on this file.
Copying or storing any contents except as provided by Supernet Technologies India Pvt Ltd is expressly prohibited without prior written permission of the Supernet Technologies.
For permission to use the content, please contact legal@supernet-india.com.
*/
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServerURL    string
	DBConnection string
}

type Language struct {
	EN_NoMappedConsumer string
	EN_ReqUrlNotFound   string
	EN_P2C_Status       string
	EN_Consumer_Status  string

	FR_NoMappedConsumer string
	FR_ReqUrlNotFound   string
	FR_P2C_Status       string
	FR_Consumer_Status  string
}

var CTX Config
var L Language

func Init() {

	// assign LOCAL or SERVER as deployment

	deployment := "LOCAL"

	if deployment == "SERVER" {

		// viper.SetConfigType("yaml")
		// viper.SetConfigName("config")
		// viper.AddConfigPath("/conf")
		// err := viper.ReadInConfig()
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// viper.AutomaticEnv()
		// viper.SetConfigType("yml")
		// err1 := viper.Unmarshal(&CTX)
		// if err1 != nil {
		// 	fmt.Printf("Unable to decode into struct, %v", err1)
		// }
		// err2 := viper.Unmarshal(&L)
		// if err2 != nil {
		// 	fmt.Printf("Unable to decode into struct, %v", err2)
		// }

	} else if deployment == "LOCAL" {

		// // Set the file name of the configurations file
		viper.SetConfigName("config")

		// // //8001 Set the path to look for the configurations file
		viper.AddConfigPath(".")

		// // // Enable VIPER to read Environment Variables
		viper.AutomaticEnv()

		viper.SetConfigType("yml")

		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Error reading config file, %s", err)
		}

		err := viper.Unmarshal(&CTX)
		if err != nil {
			fmt.Printf("Unable to decode into struct, %v", err)
		}

		err2 := viper.Unmarshal(&L)
		if err2 != nil {
			fmt.Printf("Unable to decode into struct, %v", err2)
		}

	}

}
