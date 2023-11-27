/*Copyright (C) 2020-2030 Supernet Technologies India Private Limited. All rights reserved.
All material appearing on this repository, folder, sub folder and file(s) (“Content/Code etc”) is protected by copyright laws and is the property of Supernet Technologies India Pvt Ltd. You may not copy, reproduce, distribute, publish, display, perform, modify, create derivative works, transmit, or in any way exploit any such content, nor may you distribute any part of this content over any network, including a local area network, sell, offer it for sale, or use such content to construct any kind of database, you may not alter or remove any copyright or other notice from copies of contents on this file.
Copying or storing any contents except as provided by Supernet Technologies India Pvt Ltd is expressly prohibited without prior written permission of the Supernet Technologies.
For permission to use the content, please contact legal@supernet-india.com.
*/
package entity

type Availability struct {
	Service string `json:"service"`
	Biller  string `json:"biller"`
	Channel string `json:"channel"`
}

type SuperpayPgsRequest struct {
	Billerinfosresp            string       `json:"biller_info"`
	Access_code                string       `json:"access_code"`
	SuperesbProducerAccessCode string       `json:"superesb_producer_access_code"`
	Channel                    string       `json:"channel"`
	Request_id                 string       `json:"request_id"`
	User_Type                  string       `json:"user_type"`
	Service_Name               string       `json:"service_name"`
	Mobile                     string       `json:"mobile" valid:"Required"`
	Token                      string       `json:",omitempty"`
	Operator                   string       `json:"operator"`
	AccountId                  string       `json:"account_id"`
	BillerName                 string       `json:"biller_name"`
	BillerID                   string       `json:"biller_id"`
	Txnnumber                  string       `json:"txnnumber"`
	Notification_URL           string       `json:"notification_url"`
	Return_URL                 string       `json:"return_url"`
	Cancel_URL                 string       `json:"cancel_url"`
	Customerinfo               Customerinfo `json:"Customerinfo"`
	Language                   string       `json:"language`
	ClientURLUsed              string       `json:"client_url_used"`
	EsbTxnId                   string       `json:"esb_txn_id"`
}

type Customerinfo struct {
	Authenticator1 string `json:"authenticator1"`
	Authenticator2 string `json:"authenticator2"`
	Authenticator3 string `json:"authenticator3"`
	Authenticator4 string `json:"authenticator4"`
	Authenticator5 string `json:"authenticator5"`
}

type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
