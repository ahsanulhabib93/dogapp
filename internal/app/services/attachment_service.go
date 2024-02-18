package services

import (
	"context"
	"fmt"

	attachmentpb "github.com/voonik/goConnect/api/go/ss2/attachment"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

type AttachmentService struct{}

type AttachmentServiceInterface interface {
	AddAttachment(ctx context.Context, params *attachmentpb.AddAttachmentParams) (*attachmentpb.BasicApiResponse, error)
	RemoveAttachment(ctx context.Context, params *attachmentpb.RemoveAttachmentParams) (*attachmentpb.BasicApiResponse, error)
	GetAttachmentFileTypes(ctx context.Context, params *attachmentpb.GetAttachmentFileTypesParams) (*attachmentpb.AttachmentFileTypesResponse, error)
}

func GetAttachmentServiceInstance() AttachmentServiceInterface {
	return new(AttachmentService)
}
func (service *AttachmentService) AddAttachment(
	ctx context.Context,
	params *attachmentpb.AddAttachmentParams,
) (*attachmentpb.BasicApiResponse, error) {
	resp := &attachmentpb.BasicApiResponse{Success: false}
	if params.AttachableId == utils.Zero || params.FileUrl == utils.EmptyString || params.ReferenceNumber == utils.EmptyString {
		resp.Message = "Required params missing"
		return resp, nil
	}
	attachableType := utils.AttachableType(params.GetAttachableType())
	fileTypeStr := params.GetFileType()

	fileType, ok := utils.FileTypeMapping[fileTypeStr]
	if !ok {
		resp.Message = "Invalid file type"
		return resp, nil
	}

	validFileTypes, ok := utils.AttachableFileTypeMapping[attachableType]
	if !ok {
		resp.Message = "Invalid attachable type"
		return resp, nil
	}

	if !utils.Includes(validFileTypes, fileType) {
		resp.Message = "Incompatible attachable type and file type"
		return resp, nil
	}
	attachableModel := helpers.GetModelByAttachableType(attachableType)

	existingAttachment := &models.Attachment{}
	database.DBAPM(ctx).Model(&models.Attachment{}).Where(&models.Attachment{AttachableType: attachableType, AttachableID: params.AttachableId, FileType: fileType}).First(existingAttachment)
	if existingAttachment.ID != utils.Zero {
		resp.Message = "Attachable already uploaded for this filetype"
		return resp, nil
	}

	err := database.DBAPM(ctx).Model(attachableModel).Where("id = ?", params.GetAttachableId()).First(attachableModel).Error
	if err != nil {
		resp.Message = "Attachable not found"
		return resp, nil
	}

	attachment := models.Attachment{
		AttachableID:    params.GetAttachableId(),
		FileURL:         params.GetFileUrl(),
		AttachableType:  attachableType,
		ReferenceNumber: params.GetReferenceNumber(),
		FileType:        fileType,
	}

	err = database.DBAPM(ctx).Save(&attachment).Error
	if err != nil {
		resp.Message = fmt.Sprintf("Error while creating attachment: %s", err.Error())
		return resp, nil
	}

	resp.Success = true
	resp.Message = "Attachment added successfully"
	return resp, nil
}

func (service *AttachmentService) RemoveAttachment(ctx context.Context, params *attachmentpb.RemoveAttachmentParams) (*attachmentpb.BasicApiResponse, error) {
	resp := &attachmentpb.BasicApiResponse{Success: false}
	if params.GetAttachmentId() == utils.Zero || params.GetAttachableId() == utils.Zero || params.GetAttachableType() == utils.Zero {
		resp.Message = "Required params missing"
		return resp, nil
	}
	var attachment models.Attachment
	query := database.DBAPM(ctx).Model(&models.Attachment{}).Where(&models.Attachment{
		AttachableID:   params.GetAttachableId(),
		AttachableType: utils.AttachableType(params.GetAttachableType()),
	})
	err := query.First(&attachment, params.AttachmentId).Error
	if err != nil {
		resp.Message = "Attachment not found"
		return resp, nil
	}

	if err := database.DBAPM(ctx).Delete(&attachment).Error; err != nil {
		resp.Message = "Error deleting attachment"
		return resp, nil
	}

	resp.Success = true
	resp.Message = "Attachment removed successfully"
	return resp, nil
}

func (service *AttachmentService) GetAttachmentFileTypes(ctx context.Context, params *attachmentpb.GetAttachmentFileTypesParams) (*attachmentpb.AttachmentFileTypesResponse, error) {
	resp := &attachmentpb.AttachmentFileTypesResponse{
		Data: &attachmentpb.AttachmentFileTypes{
			FileTypes: []string{},
		},
	}
	activeFileTypes := aaaModels.GetAppPreferenceServiceInstance().GetValue(ctx, "active_file_types", []string{}).([]string)
	fileTypes, _ := utils.AttachableFileTypeMapping[utils.AttachableType(params.GetAttachableType())]

	for _, fileType := range fileTypes {
		if utils.Includes(activeFileTypes, fileType.String()) {
			resp.Data.FileTypes = append(resp.Data.FileTypes, fileType.String())
		}
	}

	return resp, nil
}
