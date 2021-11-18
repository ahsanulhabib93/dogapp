package main

import (
	"github.com/qor/validations"
	supplierpb "github.com/voonik/goConnect/api/go/supplier_service/supplier"
	"github.com/voonik/goFramework/pkg/config"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/goFramework/pkg/grpc/server"
	"github.com/voonik/supplier_service/internal/app/handlers"
	"github.com/voonik/supplier_service/internal/app/migrations"
	"google.golang.org/grpc/reflection"
)

func init() {
	migrations.Migrate()
	validations.RegisterCallbacks(database.DB())
}

func main() {
	server.Init()
	supplierpb.RegisterSupplierServer(server.GrpcServer, handlers.GetSupplierInstance())

	if config.GRPCServerConfigReflection() {
		reflection.Register(server.GrpcServer)
	}

	server.Start()
	defer server.Finish()

}
