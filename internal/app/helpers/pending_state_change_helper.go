package helpers

import (
	"log"
	"time"

	"github.com/voonik/ss2/internal/app/models"

	"github.com/shopuptech/work"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/goFramework/pkg/database"
	goWorker "github.com/voonik/goFramework/pkg/worker"
)

func ChangePendingState(c *goWorker.VaccountContext, job *work.Job) error {
	log.Println("Start Change Supplier Status from Pending to Verification Failed Job")

	supplierIds := []uint64{}
	noOfDay := aaaModels.GetAppPreferenceServiceInstance().GetValue(c.GetContext(), "supplier_auto_status_change_duration", int64(7)).(int64)
	dateTime := -time.Hour * 24 * time.Duration(noOfDay)
	lastWeek := time.Now().Add(dateTime)
	err := database.DBAPM(c.GetContext()).
		Model(&models.Supplier{}).
		Where("suppliers.status = ?", models.SupplierStatusPending).
		Where("suppliers.updated_at < ?", lastWeek).
		Pluck("id", &supplierIds).
		Error

	if err != nil || len(supplierIds) == 0 {
		log.Printf("ChangePendingState: terminating job. Supplier Count: %d. Error: %v", len(supplierIds), err)
		return err
	}

	log.Printf("ChangePendingState: Number of Supplier in Pending State till date(%s): %d\n", lastWeek.String(), len(supplierIds))

	err = database.DBAPM(c.GetContext()).
		Table("suppliers").
		Where("vaccount_id = ?", c.VaccountID).
		Where("id in (?)", supplierIds).
		Update("status", models.SupplierStatusFailed).
		Update("updated_at", time.Now()).
		Error
	if err != nil {
		log.Println("ChangePendingState: Error while updating supplier status. Error: ", err.Error())
		return err
	}

	log.Println("ChangePendingState: Job has been finished for vaccount = ", c.VaccountID)
	return nil
}
