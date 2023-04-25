package script

import (
	"context"

	"log"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/rodaine/table"

	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

const (
	fileName  = "internal/app/helpers/scripts/supplier_type_change.xlsx"
	sheetName = "Sheet1"
)

type SupplierData struct {
	ID           uint64
	SupplierType utils.SupplierType
}

var supplierTypeValueMap = map[string]utils.SupplierType{
	"L0":      utils.L0,
	"L1":      utils.L1,
	"L2":      utils.L2,
	"L3":      utils.L3,
	"Hlc":     utils.Hlc,
	"Captive": utils.Captive,
	"Driver":  utils.Driver,
}

var dataFromExcel []SupplierData

func UpdateSupplierType(ctx context.Context) {
	startTime := time.Now()
	log.Println("\n\nSTART script")

	readData()
	updateSuppliers(ctx)

	endTime := time.Now()
	diff := endTime.Sub(startTime)
	log.Println("\nSCRIPT COMPLETED IN ", diff)
}

func updateSuppliers(ctx context.Context) {

	log.Println("Total number of parsed rows: ", len(dataFromExcel))
	for _, supplier := range dataFromExcel {
		log.Println(supplier.ID, supplier.SupplierType)
		err := database.DBAPM(ctx).Debug().Model(&models.Supplier{}).Where("id = ?", supplier.ID).Update("supplier_type", supplier.SupplierType).Error
		if err != nil {
			log.Println("ERROR: ", err.Error())
		}
	}
}
func readData() {
	file, openErr := excelize.OpenFile(fileName)
	if openErr != nil {
		log.Fatal("error in opening file")
	}

	rows := file.GetRows(sheetName)
	log.Println("Total no of rows in excel: ", len(rows))
	log.Println("Data From EXCEL")
	tbl := table.New("ID", "SupplierType", "Parsed SupplierType")
	tbl.WithWriter(log.Writer())
	for index, row := range rows {
		if len(row) != 2 {
			log.Printf("invalid data at row: %d\n", index+1)
			continue
		}
		parsedID, idErr := strconv.ParseUint(row[0], 10, 64)
		parsedSupplierType := supplierTypeValueMap[row[1]]
		tbl.AddRow(row[0], row[1], parsedSupplierType)

		if idErr != nil {
			log.Printf("invalid supplier_id at row: %d\n", index+1)
			continue
		}
		if parsedSupplierType == 0 {
			log.Printf("invalid supplier_type at row: %d\n", index+1)
			continue
		}
		parsedRow := SupplierData{
			ID:           parsedID,
			SupplierType: parsedSupplierType,
		}
		dataFromExcel = append(dataFromExcel, parsedRow)
	}
	tbl.Print()
}
