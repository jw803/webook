package test_data

const (
	TestFileChannelShopline = "shopline"
	TestFileChannelOthers   = "others"
	TestFileChannelPayeasy  = "payeasy"
	TestFileChannelMomo     = "momo"
	TestFileChannelAmazon   = "amazon"
)

const (
	TestFileTypeSale      = "sale"
	TestFileTypeRefund    = "refund"
	TestFileTypeSalePoint = "sale_point"
)

const (
	TestFileStatusUploading = "uploading"
	TestFileStatusSuccess   = "success"
	TestFileStatusFailed    = "failed"
)

const (
	TestFileExportStatusInit     = "init"
	TestFileExportStatusWait     = "wait"
	TestFileExportStatusFinished = "finished"
	TestFileExportStatusSkip     = "skip"
	TestFileExportStatusFailed   = "failed"
)

const (
	TestERPFileDownloadStatusNotYet     = "not_yet_downloaded"
	TestERPFileDownloadStatusDownloaded = "downloaded"
)

const (
	TestChannelMarkEC      = "A"
	TestChannelMarkShopee  = "B"
	TestChannelMarkPayeasy = "C"
	TestChannelMarkMomo    = "D"
	TestChannelMarkAmazon  = "E"
	TestChannelMarkOthers  = "G"
)

const (
	TestISOCurrencyNTD = "NTD"
	TestISOCurrencyUSD = "USD"
)

const (
	TestWarehouseIdOthers   = "warehouse_id_001"
	TestWarehouseIdShopline = "B201"
)

var (
	TestERPSaleFileTitleColumns = []string{"客戶編號", "客戶全稱", "產品編號", "品名規格", "倉庫編號", "數量", "單價", "金額",
		"含稅金額", "聯絡電話", "行動電話", "送貨地址", "幣別", "帳款歸屬", "客戶訂單", "聯絡人員", "單據編號", "訂單日期", "備註"}
)

const (
	TestSLERPProductSKUDelivery       = "Z04005"
	TestSLERPProductSKUCustomDiscount = "Z04006"
	TestSLERPProductSKUPoint          = "Z04001"
	TestSLERPProductSKUCredit         = "Z04003"
)

const (
	TestSLERPProductNameDelivery       = "運費"
	TestSLERPProductNameCustomDiscount = "自訂折扣"
	TestSLERPProductNamePoint          = "購物金"
	TestSLERPProductNameCredit         = "折價"
)
