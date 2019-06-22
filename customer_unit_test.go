package customer_test

import (
	"context"
	pb "customer/pb"
	"testing"

	"google.golang.org/genproto/protobuf/field_mask"
)

func TestCustomerServer_CreateCustomer(t *testing.T) {

	type args struct {
		ctx context.Context
		req *pb.CreateCustomerRequest
	}

	type wants struct {
		Customer *pb.Customer
		wantErr  bool
	}

	tests := []struct {
		name  string
		a     *args
		w     *wants
		setup func(a *args, w *wants)
	}{
		{
			name: "Invalid Customer",
			a: &args{
				ctx: context.Background(),
				req: &pb.CreateCustomerRequest{},
			},
			w: &wants{
				Customer: nil,
				wantErr:  true,
			},
			setup: func(a *args, w *wants) {

				a.req.Customer = getDummyCustomer()

				a.req.Customer.FirstName = "12345678901234567890123456789012345678901234567890"

				a.req.Customer.LastName = "12345678901234567890123456789012345678901234567890"

				a.req.Customer.Email = "Invalid Email"

			},
		},
	}

	srv, err := getServer()

	if err != nil {
		t.Errorf("DB connection error: %s", err)
		return
	}

	for _, tt := range tests {

		tt.setup(tt.a, tt.w)

		_, err = srv.CreateCustomer(tt.a.ctx, tt.a.req)
		if (err != nil) != tt.w.wantErr {
			t.Errorf("Server.CreateUserProfile() error = %v, wantErr = %v", err, tt.w.wantErr)
			return
		}

	}
}

func TestCustomerServer_GetCustomer(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.GetCustomerRequest
	}

	type wants struct {
		Customer *pb.Customer
		wantErr  bool
	}

	tests := []struct {
		name  string
		a     *args
		w     *wants
		setup func(a *args, w *wants)
	}{
		{
			name: "invalid get",
			a: &args{
				ctx: context.Background(),
				req: &pb.GetCustomerRequest{},
			},
			w: &wants{
				Customer: nil,
				wantErr:  true,
			},
			setup: func(a *args, w *wants) {
				a.req.Id = " "
				a.req.ViewMask = &field_mask.FieldMask{Paths: []string{}}
			},
		},
	}

	srv, err := getServer()

	if err != nil {
		t.Errorf("DB connection error: %s", err)
		return
	}

	for _, tt := range tests {
		tt.setup(tt.a, tt.w)

		_, err = srv.GetCustomer(tt.a.ctx, tt.a.req)
		if (err != nil) != tt.w.wantErr {
			t.Errorf("Server.CreateUserProfile() error = %v, wantErr = %v", err, tt.w.wantErr)
			return
		}
	}
}
