package main

import (
  "github.com/qor/validations"
  "github.com/voonik/supplier_service/internal/app/migrations"
  "github.com/voonik/goFramework/pkg/config"
  "github.com/voonik/goFramework/pkg/database"
  "github.com/voonik/goFramework/pkg/grpc/server"
  "google.golang.org/grpc/reflection"
)

func init() {
  migrations.Migrate()
  validations.RegisterCallbacks(database.DB())
}

func main() {
  server.Init()

  if config.GRPCServerConfigReflection() {
    reflection.Register(server.GrpcServer)
  }

  server.Start()
  defer server.Finish()

}
