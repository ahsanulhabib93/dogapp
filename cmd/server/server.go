package main

import (
	"github.com/qor/validations"
	kampb "github.com/voonik/goConnect/api/go/ss2/key_account_manager"
	paymentpb "github.com/voonik/goConnect/api/go/ss2/payment_account_detail"
	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	addresspb "github.com/voonik/goConnect/api/go/ss2/supplier_address"
	"github.com/voonik/goFramework/pkg/config"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/goFramework/pkg/grpc/server"
	"github.com/voonik/ss2/internal/app/handlers"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/migrations"
	"google.golang.org/grpc/reflection"
)

func init() {
	migrations.Migrate()
	validations.RegisterCallbacks(database.DB())
}

func main() {
	server.Init()
	supplierpb.RegisterSupplierServer(server.GrpcServer, handlers.GetSupplierInstance())
	addresspb.RegisterSupplierAddressServer(server.GrpcServer, handlers.GetSupplierAddressInstance())
	paymentpb.RegisterPaymentAccountDetailServer(server.GrpcServer, handlers.GetPaymentAccountDetailInstance())
	kampb.RegisterKeyAccountManagerServer(server.GrpcServer, handlers.GetKeyAccountManagerInstance())

	if config.GRPCServerConfigReflection() {
		reflection.Register(server.GrpcServer)
	}

	if config.AsynqConfigEnabled() {
		helpers.InitGoJobsWorker()
	}

	server.Start()
	defer server.Finish()

}
