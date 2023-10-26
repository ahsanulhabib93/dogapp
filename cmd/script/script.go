package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/voonik/goFramework/pkg/misc"
	script "github.com/voonik/ss2/internal/app/helpers/scripts"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		updateSupplierType(),
		addSuppliersInBulk(),
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}

func updateSupplierType() cli.Command {
	return cli.Command{
		Name:  "updateCargoSupplierType",
		Usage: "Update supplier type of cargo suppliers",
		Action: func(ctx *cli.Context) error {
			bgCtx := context.Background()
			threadObject := &misc.ThreadObject{
				VaccountId:    12,
				PortalId:      12,
				CurrentActId:  12,
				XForwardedFor: "5079327",
				UserData: &misc.UserData{
					UserId: 18,
					Name:   "AFT2User",
					Email:  "aft2user@gmail.com",
					Phone:  "8801855533367",
				},
			}
			context := misc.SetInContextThreadObject(bgCtx, threadObject)
			script.UpdateSupplierType(context)
			return nil
		},
	}
}
func addSuppliersInBulk() cli.Command {
	return cli.Command{
		Name:  "addSuppliersInBulk",
		Usage: "Add new suppliers in bulk",
		Action: func(ctx *cli.Context) error {
			bgCtx := context.Background()
			threadObject := &misc.ThreadObject{
				VaccountId:    12,
				PortalId:      12,
				CurrentActId:  12,
				XForwardedFor: "5079327",
				UserData: &misc.UserData{
					UserId: 18,
					Name:   "AFT2User",
					Email:  "aft2user@gmail.com",
					Phone:  "8801855533367",
				},
			}
			context := misc.SetInContextThreadObject(bgCtx, threadObject)
			script.AddSuppliersFromExcel(context)
			return nil
		},
	}
}
