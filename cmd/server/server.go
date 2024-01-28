package main

import (
	"github.com/qor/validations"
	attachmentpb "github.com/voonik/goConnect/api/go/ss2/attachment"
	kampb "github.com/voonik/goConnect/api/go/ss2/key_account_manager"
	psmpb "github.com/voonik/goConnect/api/go/ss2/partner_service_mapping"
	paymentpb "github.com/voonik/goConnect/api/go/ss2/payment_account_detail"
	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	sampb "github.com/voonik/goConnect/api/go/ss2/seller_account_manager"
	sbdpb "github.com/voonik/goConnect/api/go/ss2/seller_bank_detail"
	spdpb "github.com/voonik/goConnect/api/go/ss2/seller_pricing_detail"
	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	addresspb "github.com/voonik/goConnect/api/go/ss2/supplier_address"
	vapb "github.com/voonik/goConnect/api/go/ss2/vendor_address"
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
	psmpb.RegisterPartnerServiceMappingServer(server.GrpcServer, handlers.GetPartnerServiceMappingInstance())
	spb.RegisterSellerServer(server.GrpcServer, handlers.GetSellerInstance())
	sbdpb.RegisterSellerBankDetailServer(server.GrpcServer, handlers.GetSellerBankDetailInstance())
	spdpb.RegisterSellerPricingDetailServer(server.GrpcServer, handlers.GetSellerPricingDetailInstance())
	vapb.RegisterVendorAddressServer(server.GrpcServer, handlers.GetVendorAddressInstance())
	attachmentpb.RegisterAttachmentServer(server.GrpcServer, handlers.GetAttachmentInstance())
	sampb.RegisterSellerAccountManagerServer(server.GrpcServer, handlers.GetSellerAccountManagerInstance())

	if config.GRPCServerConfigReflection() {
		reflection.Register(server.GrpcServer)
	}

	helpers.InitGoJobsWorker()

	server.Start()
	defer server.Finish()

}
