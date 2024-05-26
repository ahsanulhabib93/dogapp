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
		addSuppliersInBulk(),
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}

func addSuppliersInBulk() cli.Command {
	return cli.Command{
		Name:  "addSuppliersInBulk",
		Usage: "Add new suppliers in bulk",
		Action: func(ctx *cli.Context) error {
			bgCtx := context.Background()
			threadObject := &misc.ThreadObject{
				VaccountId:    2,
				PortalId:      2,
				CurrentActId:  2,
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
