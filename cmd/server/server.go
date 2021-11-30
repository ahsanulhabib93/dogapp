package main

import (
	"github.com/qor/validations"
	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	paymentpb "github.com/voonik/goConnect/api/go/ss2/payment_account_detail"
	kampb "github.com/voonik/goConnect/api/go/ss2/key_account_manager"
	"github.com/voonik/goFramework/pkg/config"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/goFramework/pkg/grpc/server"
	"github.com/voonik/ss2/internal/app/handlers"
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
	paymentpb.RegisterPaymentAccountDetailServer(server.GrpcServer, handlers.GetPaymentAccountDetailInstance())
	kampb.RegisterKeyAccountManagerServer(server.GrpcServer, handlers.GetKeyAccountManagerInstance())

	if config.GRPCServerConfigReflection() {
		reflection.Register(server.GrpcServer)
	}

	server.Start()
	defer server.Finish()

}
