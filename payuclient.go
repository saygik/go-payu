package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/saygik/go-payu/firebase"
	"github.com/saygik/go-payu/payu"
)

func getAuth(c *gin.Context) {

	fmt.Println("------------------")
	client := payu.NewClient(AppConfig.PayU.ClientID, AppConfig.PayU.Secret, AppConfig.PayU.PayUBase)
	accessToken, err := client.GetAccessToken()

	if err != nil {
		c.JSON(400, gin.H{"Message": fmt.Sprintf("Could not get auth :  "), "Error": "Not auth"})
	} else {
		c.JSON(200, gin.H{"Message": fmt.Sprintf("Information  auth-token : %s ", accessToken.Token), "Error": ""})
	}
}

func setNotify(c *gin.Context) {
	fmt.Println("------------------")
	var PayuNotifyer payu.PayuNotifyer
	c.BindJSON(&PayuNotifyer)
	fmt.Println("Payment status: %s", PayuNotifyer.Order.Status)
	_ = firebase.UpdateOrderStatus(PayuNotifyer)
	if PayuNotifyer.Order.Status != "COMPLETED" {
		c.JSON(202, gin.H{"Message": "Ok"})
	} else {
		c.JSON(200, gin.H{"Message": "Ok"})
	}
}
func createOrder(c *gin.Context) {

	fmt.Println("------------------")
	payuClient := payu.NewClient(AppConfig.PayU.ClientID, AppConfig.PayU.Secret, AppConfig.PayU.PayUBase)

	//	accessToken, err := client.CreatePayment(p)
	var payOrder payu.Order
	c.BindJSON(&payOrder)
	payOrder.MerchantPosId = AppConfig.PayU.MerchantPosId
	v, err := payuClient.CreateOrder(payOrder)
	if err != nil {
		c.JSON(400, gin.H{"Message": fmt.Sprintf("Could not create order"), "Error": "No order"})
	} else {
		err = firebase.UpdateOrder(v.ExtOrderId, "SEND_TO_PAYU")
		if err != nil {
			c.JSON(400, gin.H{"Message": fmt.Sprintf("Could not update order in firebase")})
		} else {
			c.JSON(200, gin.H{"data": v})
		}
	}
}
