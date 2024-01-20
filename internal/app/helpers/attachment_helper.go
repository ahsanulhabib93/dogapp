package helpers

import (
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

func GetModelByAttachableType(attachableType utils.AttachableType) interface{} {
	var object interface{}

	switch attachableType {
	case utils.AttachableTypeSupplier:
		object = &models.Supplier{}
	case utils.AttachableTypePartnerServiceMapping:
		object = &models.PartnerServiceMapping{}
	}

	return object
}
