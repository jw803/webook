package test_data

import (
	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/integration_test/utils"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/repository/dao"
)

func GenerateOrderFIleUploadLogList(testName string, number int) []any {
	orderFileUploadLogList := make([]any, 0)

	for i := 1; i <= number; i++ {
		objectId := utils.GetObjectIdByTestNameAndNumber(testName, i)
		orderFileUploadLog := dao.OrderFileUploadLog{
			Id:         objectId,
			MerchantId: TestMerchantConfigWstyle.MerchantId,
			UploadedAt: utils.InitialTime,
			CreatedAt:  utils.InitialTime,
		}
		orderFileUploadLogList = append(orderFileUploadLogList, orderFileUploadLog)
	}
	return orderFileUploadLogList
}
