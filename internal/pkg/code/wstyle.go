package code

const (
	ErrBadRequest int = 400000
	// ErrBadRequestFile - 400: Bad request file.
	ErrBadRequestFile int = 400001
	// ErrFileTooLarge - 400: Request file too large.
	ErrFileTooLarge int = 400002
	// ErrFileExtension - 400: Merchant config is not ready.
	ErrFileExtension int = 400003
	// ErrFileHasAOrderNumberBothFullReturnAndPartialReturn - 400: File has an order number both full return and partial return.
	ErrFileHasAOrderNumberBothFullReturnAndPartialReturn int = 400004
	// ErrFileHasAlreadyExportedOrderIds - 400: File has already exported order ids.
	ErrFileHasAlreadyExportedOrderIds int = 400005
	// ErrNotAllowUploadAtThisTime - 400: Not allow upload at this time.
	ErrNotAllowUploadAtThisTime int = 400006
	// ErrThisIsYesterdayRequest - 400: This is yesterday request.
	ErrThisIsYesterdayRequest int = 400007
	// ErrFileHasInvalidOrderIdFormat - 400: File has invalid order id format.
	ErrFileHasInvalidOrderIdFormat int = 400008
	// ErrRefundFileHasUnknownOrderId - 400: Refund file has unknown order id.
	ErrRefundFileHasUnknownOrderId int = 400009
	// ErrSomeoneUploadFileBeforeYou - 400: Someone upload file before you.
	ErrSomeoneUploadFileBeforeYou int = 400010

	// ErrFileNoProductRow 檔案沒有商品資料
	ErrFileNoProductRow int = 400100
	ErrFileDataInvalid  int = 400101

	ErrNotAllowUploadDueToProcessing int = 422000
)

const (
	ErrInternal int = iota + 500000
	// ErrMerchantConfigNotReady - 500: Merchant config is not ready.
	ErrMerchantConfigNotReady
	ErrClientNotReady
	ErrSaleOrderNotFound
	ErrInvoiceNotFound
	ErrProcessReturnInventory
	ErrDuplicateSaleOrder
	ErrDuplicateSaleInvoice
)

const (
	ErrOneShip int = iota + 500100
)

const (
	ErrSLOA int = iota + 500200
	ErrSLOANotFound
)
