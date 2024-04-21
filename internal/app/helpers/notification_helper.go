package helpers

// func SendStatusChangeEmailNotification(ctx context.Context, supplier models.Supplier, status, reason string) error {
// 	return nil
// }

// type VigeonAPIHelper struct{}

// type VigeonAPIHelperInterface interface {
// 	SendEmailAPI(ctx context.Context, emailParam *notify.EmailParam) *notify.EmailResp
// }

// var vigeonApiHelper VigeonAPIHelperInterface

// // InjectMockVigeonAPIHelperInstance ...
// func InjectMockVigeonAPIHelperInstance(mockObj VigeonAPIHelperInterface) {
// 	vigeonApiHelper = mockObj
// }

// // getVigeonAPIHelperInstance ...
// func getVigeonAPIHelperInstance() VigeonAPIHelperInterface {
// 	if vigeonApiHelper == nil {
// 		vigeonApiHelper = new(VigeonAPIHelper)
// 	}
// 	return vigeonApiHelper
// }

// func (apiHelper *VigeonAPIHelper) SendEmailAPI(ctx context.Context, emailParam *notify.EmailParam) *notify.EmailResp {
// 	resp, err := vigeon.Notify().EmailNotification(ctx, emailParam)
// 	if err != nil {
// 		log.Println("SentEmailAPI: Failed to sent email. Error: ", err.Error())
// 	}

// 	return resp
// }
