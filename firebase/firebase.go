package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/saygik/go-payu/payu"
	"google.golang.org/api/option"
)

func NewFirestoreClient() (context.Context, *firestore.Client, error) {
	ctx := context.Background()
	sa := option.WithCredentialsFile("crystal-rental-car-go.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, nil, err
	}
	client, err := app.Firestore(ctx)

	if err != nil {
		return nil, nil, err
	}

	return ctx, client, nil
}

/*
func NewFirestoreClient1() (firestore.Client)  {
	ctx := context.Background()

	sa := option.WithCredentialsFile("crystal-rental-car-go.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	//	Rta1hPo7Cd4kZPAFvoH3
	dsnap, err := client.Collection("cars").Doc("Rta1hPo7Cd4kZPAFvoH3").Get(ctx)
	if err != nil {
		c.JSON(200, gin.H{"Message": "OK firebase"})
	}
	m := dsnap.Data()
	fmt.Printf("Document data: %#v\n", m)

	c.JSON(200, gin.H{"Message": "OK firebase"})

}
*/

func GetOneCar() error {
	ctx, client, err := NewFirestoreClient()
	defer client.Close()
	if err != nil {
		return err
	}
	dsnap, err := client.Collection("orders").Doc("1ORTGLo8w2Zklxearf9q").Get(ctx)
	if err != nil {
		return err
	}
	m := dsnap.Data()
	fmt.Printf("Document data: %#v\n", m)
	d := m["sta"]
	fmt.Printf("Document data: %#v\n", d)

	return nil
}

func UpdateOrderStatus(payuNotifyer payu.PayuNotifyer) (string, string, error) {
	orderId := payuNotifyer.Order.ExtOrderId
	status := payuNotifyer.Order.Status
	ctx, client, err := NewFirestoreClient()
	defer client.Close()
	if err != nil {
		return "", "", err
	}
	snapshot, err := client.Collection("orders").Doc(orderId).Get(ctx)
	if err != nil {
		return "", "", err
	}
	snapshotData := snapshot.Data()
	currentOrderStatus := snapshotData["status"]
	currentCarId := snapshotData["carId"]
	startDate := snapshotData["startDate"]
	endDate := snapshotData["endDate"]
	if currentOrderStatus == nil {
		err := fmt.Errorf("Not status for order with id %s", orderId)
		return "", "", err
	}
	if currentOrderStatus != status && currentOrderStatus != "COMPLETED" {
		if status == "COMPLETED" {
			//						ss:= firestore.ServerTimestamp
			ref := client.Collection("complatedOrders").NewDoc()

			// later...
			_, err := ref.Set(ctx, payuNotifyer)
			if err != nil {
				return "", "", err
			}
			_, err = ref.Set(ctx, map[string]interface{}{
				"createdAt": firestore.ServerTimestamp,
				"carId":     currentCarId,
				"startDate": startDate,
				"endDate":   endDate,
			}, firestore.MergeAll)

			//			_, _, err = client.Collection("complatedOrders").Add(ctx, payuNotifyer)
			if err != nil {
				return "", "", err
			}
		}
		_, err = client.Collection("orders").Doc(orderId).Set(ctx, map[string]interface{}{
			"status": status,
		}, firestore.MergeAll)
		if err != nil {
			return "", "", err
		}
	}
	return fmt.Sprintf("%v", startDate), fmt.Sprintf("%v", endDate), nil
}

//func UpdateOrderStatus3 (orderId string, status string) error {
//	ctx, client, err := NewFirestoreClient()
//	defer client.Close()
//	if err != nil {
//		return err
//	}
//	_, err = client.Collection("orders").Doc(orderId).Get(ctx)
//	if err != nil {
//		return err
//	}
//	_, err = client.Collection("orders").Doc(orderId).Set(ctx, map[string]interface{}{
//		"status": status,
//	}, firestore.MergeAll)
//	if err != nil {
//		// Handle any errors in an appropriate way, such as returning them.
//		return err
//	}
//	return nil
//}
func UpdateOrder(orderId string, status string) error {
	ctx, client, err := NewFirestoreClient()
	defer client.Close()
	if err != nil {
		return err
	}
	_, err = client.Collection("orders").Doc(orderId).Get(ctx)
	if err != nil {
		return err
	}
	_, err = client.Collection("orders").Doc(orderId).Set(ctx, map[string]interface{}{
		"status": status,
	}, firestore.MergeAll)
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		return err
	}
	return nil
}
