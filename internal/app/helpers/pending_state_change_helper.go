package helpers

import (
	"log"
	"time"

	"github.com/voonik/ss2/internal/app/models"

	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/goFramework/pkg/database"
	goWorker "github.com/voonik/goFramework/pkg/worker"
	"github.com/voonik/work"
)

func ChangePendingState(c *goWorker.VaccountContext, job *work.Job) error {
	log.Println("Start Change Supplier Status from Pending to Verification Failed Job")

	supplierIds := []uint64{}
	noOfDay := aaaModels.GetAppPreferenceServiceInstance().GetValue(c.GetContext(), "supplier_auto_status_change_duration", int64(7)).(int64)
	lastWeek := time.Now().Add(-time.Hour * 24 * time.Duration(noOfDay))
	err := database.DBAPM(c.GetContext()).Model(&models.Supplier{}).
		Where("suppliers.status = ?", models.SupplierStatusPending).
		Where("suppliers.updated_at < ?", lastWeek).Pluck("id", &supplierIds).
		Error

	if err != nil {
		log.Println("ChangePendingState: Failed to move Pending supplier to Verification Failed State. Error:", err.Error())
		return err
	}

	log.Printf("ChangePendingState: Number of Supplier in Pending State till date(%s): %d\n", lastWeek.String(), len(supplierIds))
	if len(supplierIds) == 0 {
		return nil
	}

	err = database.DBAPM(c.GetContext()).
		Table("suppliers").
		Where("id in (?)", supplierIds).
		Select("status").
		Update("status", models.SupplierStatusFailed).
		Error
	if err != nil {
		log.Println("ChangePendingState: Error while updating supplier status. Error: ", err.Error())
		return err
	}

	return nil
}
