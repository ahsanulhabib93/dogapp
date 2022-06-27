package helpers

import (
	"fmt"
	"log"
	"time"

	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"

	"github.com/voonik/goFramework/pkg/database"
	goWorker "github.com/voonik/goFramework/pkg/worker"
	"github.com/voonik/work"
)

func ChangePendingState(c *goWorker.VaccountContext, job *work.Job) error {
	log.Println("Start Change Supplier Status from Pending to Verification Failed Job")

	supplierIds := []uint64{}
	lastWeek := time.Now().Add(-time.Hour * 24 * 7)
	err := database.DBAPM(c.GetContext()).Model(&models.Supplier{}).
		Where("suppliers.status = ?", models.SupplierStatusPending).
		Where("suppliers.created_at < ?", lastWeek).Pluck("id", &supplierIds).Error

	if err != nil {
		log.Println("ChangePendingState: Failed to move Pending supplier to Verification Failed State. Error:", err.Error())
		return err
	}

	log.Printf("ChangePendingState: Number of Supplier in Pending State till date(%s): %d\n", lastWeek.String(), len(supplierIds))
	if len(supplierIds) == 0 {
		return nil
	}

	err = database.DBAPM(c.GetContext()).
		Exec(fmt.Sprintf(
			"update suppliers set status = '%s' , updated_at = '%v' where id in (%v)",
			models.SupplierStatusFailed, time.Now(), utils.IntToString(supplierIds))).Error
	if err != nil {
		log.Println("ChangePendingState: Error while updating supplier status. Error: ", err.Error())
		return err
	}

	return nil
}
