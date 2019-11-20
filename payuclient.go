package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/saygik/go-payu/payu"
)

// OAuth protocol - client_id :	369722
// OAuth protocol - client_secret:	aa3272c123df4a43b3d26e4b4794c0c1
func getAuth(c *gin.Context) {

	fmt.Println("------------------")
	client := payu.NewClient("369722", "aa3272c123df4a43b3d26e4b4794c0c1", payu.APIBaseSandBox)
	accessToken, err := client.GetAccessToken()

	if err != nil {
		c.JSON(400, gin.H{"Message": fmt.Sprintf("Could not get auth :  "), "Error": "Not auth"})
	} else {
		c.JSON(200, gin.H{"Message": fmt.Sprintf("Information  auth-token : %s ", accessToken.Token), "Error": ""})
	}
}

func setNotify(c *gin.Context) {
	//bodyBytes, _ := c.GetRawData()
	//bodyString := string(bodyBytes)
	//
	//fmt.Printf("%q", bodyString)
	fmt.Println("------------------")
	var PayuNotifyer payu.PayuNotifyer
	c.BindJSON(&PayuNotifyer)
	//	fmt.Println(PayuNotifyer)
	fmt.Println("Payment status: %s", PayuNotifyer.Order.Status)

	if PayuNotifyer.Order.Status != "COMPLETED" {
		c.JSON(202, gin.H{"Message": "Ok"})
	} else {
		c.JSON(200, gin.H{"Message": "Ok"})
	}
}
func createOrder(c *gin.Context) {

	fmt.Println("------------------")
	client := payu.NewClient("369722", "aa3272c123df4a43b3d26e4b4794c0c1", payu.APIBaseSandBox)

	//	accessToken, err := client.CreatePayment(p)
	var payOrder payu.Order
	c.BindJSON(&payOrder)
	v, err := client.CreateOrder(payOrder)
	if err != nil {
		c.JSON(400, gin.H{"Message": fmt.Sprintf("Could not create order"), "Error": "No order"})
	} else {
		c.JSON(200, gin.H{"data": v})
	}
}
