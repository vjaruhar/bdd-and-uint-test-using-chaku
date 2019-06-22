package customer_test

import (
	"context"
	customer "customer"
	pb "customer/pb"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/golang/protobuf/ptypes"
	. "github.com/smartystreets/goconvey/convey"
	"go.appointy.com/chaku/driver"
	"go.appointy.com/idutil"
	"go.appointy.com/waqt/protos/types"
	"google.golang.org/genproto/protobuf/field_mask"
)

func TestCustomer(t *testing.T) {
	ctx := context.Background()
	srv, err := getServer()
	if err != nil {
		t.Errorf("DB connection error: %s", err)
		return
	}

	Convey("Create Cutomer", t, func() {
		cust := getDummyCustomer()

		createCust, err := srv.CreateCustomer(ctx, &pb.CreateCustomerRequest{Parent: "jhbkjb", Customer: cust})

		So(err, ShouldEqual, nil)

		CUSTId := idutil.GetId(createCust.Id)
		cust.Id = createCust.Id

		Convey("Get Customer", func() {
			_, err := srv.GetCustomer(ctx, &pb.GetCustomerRequest{
				Id:       CUSTId,
				ViewMask: &field_mask.FieldMask{Paths: []string{}}})

			So(err, ShouldEqual, nil)
			//So(getCustResp, ShouldEqual, cust)
		})

		Convey("Update Customer", func() {
			updateCust := cust
			updateCust.FirstName = "New FN"
			updateCust.LastName = "New LN"

			Convey("Success: Update Customer", func() {
				_, err := srv.UpdateCustomer(ctx, &pb.UpdateCustomerRequest{Customer: updateCust, UpdateMask: &field_mask.FieldMask{Paths: []string{"first_name", "last_name", "address"}}})

				So(err, ShouldEqual, nil)

			})

			Convey("Check Update Customer", func() {
				_, err := srv.GetCustomer(ctx, &pb.GetCustomerRequest{Id: updateCust.Id})
				So(err, ShouldEqual, nil)
				//So(resp, ShouldEqual, updateCust)

			})

		})

		Convey("Delete Customer", func() {
			Convey("Success: Delete Customer", func() {

				_, err := srv.DeleteCustomer(ctx, &pb.DeleteCustomerRequest{Id: CUSTId})

				So(err, ShouldEqual, nil)
			})

			Convey("Check Delete Customer", func() {
				_, err := srv.GetCustomer(ctx, &pb.GetCustomerRequest{Id: CUSTId})
				So(err, ShouldNotEqual, nil)

			})
		})
	})
}

func getDummyCustomer() *pb.Customer {
	return &pb.Customer{
		Id:                "svsdfssdfsfs",
		FirstName:         "FN",
		LastName:          "LN",
		Email:             "email@email.com",
		BirthDate:         ptypes.TimestampNow(),
		ProfileImage:      &types.GalleryItem{},
		Telephones:        []string{},
		Address:           &types.Address{},
		Note:              "erssfeferfev",
		Tag:               []string{},
		Timezone:          "kjnilnim",
		PreferredLanguage: "Spanish",
	}
}

func getServer() (pb.CustomersServer, error) {
	store := createStore()

	if err := store.CreateCustomerPGStore(context.Background()); err != nil {
		return nil, err
	}

	return customer.NewCustomersServer(store), nil

}

func createStore() pb.CustomerStore {
	config := getConfig()

	db, err := sql.Open("postgres", config)
	if err != nil {
		panic(fmt.Errorf("connection not open | %v", err.Error()))
	}
	if err = db.Ping(); err != nil {
		panic("ping  " + err.Error())
	}
	return pb.NewPostgresCustomerStore(db, driver.GetUserId)
}

func getConfig() string {
	config := "host=10.0.0.77 port=5432 user=postgres password=manhattan dbname=appointypostgres sslmode=disable"
	// create pgconfig.json and put your credentials in it
	pg, err := os.Open("pgconfig.json")
	if err == nil {
		pgc, err := ioutil.ReadAll(pg)
		if err != nil {
			log.Println(err)
		}

		mp := make(map[string]string, 0)
		err = json.Unmarshal(pgc, &mp)
		if err != nil {
			log.Fatalln(err)
		}
		config = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			mp["host"], mp["port"], mp["user"], mp["password"], mp["dbname"])
	}
	return config
}
