//go:build wireinject
// +build wireinject

package startup

import (
	"time"

	"github.com/google/wire"

	orderfileevent "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/interface/event/orderfile/aws/upload_consumer"
	eventorderfile "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/interface/event/orderfile/local"
	app_installation_token_create "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/interface/event/slwebhook/access_token/aws"
	invoiceevent "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/interface/event/slwebhook/invoice/aws"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/interface/event/slwebhook/member_point_create/aws"
	order_cancel "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/interface/event/slwebhook/order_cancel/aws"
	order_create "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/interface/event/slwebhook/order_create/aws"
	order_update "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/interface/event/slwebhook/order_delivery_update/aws"
	return_order_create "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/interface/event/slwebhook/return_order_create/aws"
	exportinvoiceerpfile "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/interface/job/export_erp_invoice_file"
	exportrefunderpfile "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/interface/job/export_refund_erp_file"
	exportsaleerpfile "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/interface/job/export_sale_erp_file"
	queryoneshiplogisticsorder "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/interface/job/query_oneship_logistics_order"
	erpfilelogcontroller "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/interface/web/erp_file"
	orderfilelogcontroller "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/interface/web/order_file"
	excelgenerator "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/pkg/excel_generator"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/repository"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/repository/dao"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/repository/oneship"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/repository/shopline"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/service/client"
	emailsvc "bitbucket.org/starlinglabs/cst-wstyle-integration/internal/service/email"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/service/erp_file"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/service/invoice"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/service/merchant_config"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/service/order_file_upload_log"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/service/sale_order"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/service/staff"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/pkg/email"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/pkg/storage"
	"bitbucket.org/starlinglabs/cst-wstyle-integration/pkg/storage/local"
)

var thirdProvider = wire.NewSet(InitMongoDBClient, InitMongoDB, InitLogger)
var storageProvider = wire.NewSet(local.NewLocalStorage)
var webServerProvider = wire.NewSet()
var eventConsumerProvider = wire.NewSet(NewConsumers, orderfileevent.NewOrderFileUploadEventConsumer)
var webhookConsumerProvider = wire.NewSet(NewConsumers, invoiceevent.NewInvoiceCreateConsumer, invoiceevent.NewInvoiceUpdateConsumer)
var eventProducerProvider = wire.NewSet(eventorderfile.NewLocalOrderFileProducer)

func InitOrderFileHandler(
	saleOrderDAO dao.SaleOrderDAO, orderFileUploadLogDAO dao.OrderFileUploadLogDAO,
	merchantConfigDAO dao.MerchantConfigDAO, staffDAO dao.StaffDAO, clientDAO dao.ClientDAO,
	staffSLOA shopline.StaffOpenAPI,
	orderSLOA shopline.OrderOpenAPI,
	returnOrderSLOA shopline.ReturnOrderOpenAPI,
	emailSvc email.EmailClient,
	orderDeliveryOneship oneship.OrderDeliveryOneship,
	saleOrderOneship oneship.SaleOrderOneship,
	logisticsOrderOneship oneship.LogisticsOrderOneship,
	nowFunc func() time.Time) *orderfilelogcontroller.OrderFileHandler {
	wire.Build(
		InitLogger,
		local.NewLocalStorage,

		eventProducerProvider,

		repository.NewMerchantConfigRepository,
		repository.NewStaffRepository,
		repository.NewSaleOrderRepository,
		repository.NewOrderFileUploadLogRepository,
		repository.NewClientRepository,

		sale_order.NewSaleOrderService,
		order_file_upload_log.NewOrderFileLogService,
		client.NewClientService,
		merchant_config.NewMerchantConfigService,
		staff.NewStaffService,
		emailsvc.NewEmailService,

		orderfilelogcontroller.NewOrderFileHandler,
	)
	return new(orderfilelogcontroller.OrderFileHandler)
}

func InitERPFileHandler(
	clientDAO dao.ClientDAO,
	invoiceDao dao.InvoiceDAO,
	saleOrderDAO dao.SaleOrderDAO,
	erpOrderFileDAO dao.ERPOrderFileDAO,
	erpInvoiceFileDAO dao.ERPInvoiceFileDAO,
	merchantConfigDAO dao.MerchantConfigDAO,
	staffDAO dao.StaffDAO,
	staffSLOA shopline.StaffOpenAPI,
	orderSLOA shopline.OrderOpenAPI,
	returnOrderSLOA shopline.ReturnOrderOpenAPI,
	nowFunc func() time.Time,
	storage storage.Storage,
	emailSvc email.EmailClient,
	orderDeliveryOneship oneship.OrderDeliveryOneship,
	saleOrderOneship oneship.SaleOrderOneship,
	logisticsOrderOneship oneship.LogisticsOrderOneship,
	excelSaleOrderGeneratorFactory excelgenerator.SaleOrderERPExcelGeneratorFactory,
	excelInvoiceGeneratorFactory excelgenerator.InvoiceERPExcelGeneratorFactory,
) *erpfilelogcontroller.ERPFileHandler {
	wire.Build(
		InitLogger,

		repository.NewClientRepository,
		repository.NewMerchantConfigRepository,
		repository.NewStaffRepository,
		repository.NewERPOrderFileRepository,
		repository.NewSaleOrderRepository,
		repository.NewInvoiceRepository,
		repository.NewERPInvoiceFileRepository,

		client.NewClientService,
		sale_order.NewSaleOrderService,
		erp_file.NewERPFileService,
		merchant_config.NewMerchantConfigService,
		invoice.NewInvoiceService,
		staff.NewStaffService,
		emailsvc.NewEmailService,

		erpfilelogcontroller.NewERPFileHandler,
	)
	return new(erpfilelogcontroller.ERPFileHandler)
}

func InitOrderFileUploadConsumer(
	orderFileUploadDAO dao.OrderFileUploadLogDAO,
	saleOrderDAO dao.SaleOrderDAO,
	erpFileUploadLogDAO dao.ERPOrderFileDAO,
	merchantConfigDAO dao.MerchantConfigDAO,
	staffDAO dao.StaffDAO,
	clientDAO dao.ClientDAO,
	staffSLOA shopline.StaffOpenAPI,
	orderSLOA shopline.OrderOpenAPI,
	returnOrderSLOA shopline.ReturnOrderOpenAPI,
	orderDeliveryOneship oneship.OrderDeliveryOneship,
	saleOrderOneship oneship.SaleOrderOneship,
	logisticsOrderOneship oneship.LogisticsOrderOneship,
	emailSvc email.EmailClient,
	nowFunc func() time.Time) *orderfileevent.OrderFileEventConsumer {
	wire.Build(
		InitSQS,
		InitLogger,
		local.NewLocalStorage,

		repository.NewMerchantConfigRepository,
		repository.NewStaffRepository,
		repository.NewClientRepository,
		repository.NewSaleOrderRepository,
		repository.NewOrderFileUploadLogRepository,

		client.NewClientService,
		merchant_config.NewMerchantConfigService,
		emailsvc.NewEmailService,
		sale_order.NewSaleOrderService,

		orderfileevent.NewOrderFileUploadEventConsumer,
	)
	return new(orderfileevent.OrderFileEventConsumer)
}

func InitSLWebhookOrderCreateConsumer(
	saleOrderDAO dao.SaleOrderDAO,
	clientDAO dao.ClientDAO,
	merchantDAO dao.MerchantConfigDAO,
	orderSLOA shopline.OrderOpenAPI,
	returnOrderSLOA shopline.ReturnOrderOpenAPI,
	orderDeliveryOneship oneship.OrderDeliveryOneship,
	saleOrderOneship oneship.SaleOrderOneship,
	logisticsOrderOneship oneship.LogisticsOrderOneship,
	emailSvc email.EmailClient,
	nowFunc func() time.Time) *order_create.SLWebhookOrderCreateConsumer {
	wire.Build(
		InitSQS,
		InitLogger,
		local.NewLocalStorage,

		repository.NewClientRepository,
		repository.NewMerchantConfigRepository,
		repository.NewSaleOrderRepository,

		emailsvc.NewEmailService,
		client.NewClientService,
		merchant_config.NewMerchantConfigService,
		sale_order.NewSaleOrderService,

		order_create.NewSLWebhookOrderCreateConsumer,
	)
	return new(order_create.SLWebhookOrderCreateConsumer)
}

func InitCronjobExportERPSaleFile(
	clientDAO dao.ClientDAO,
	saleOrderDAO dao.SaleOrderDAO,
	erpOrderFileDAO dao.ERPOrderFileDAO,
	erpInvoiceFileDAO dao.ERPInvoiceFileDAO,
	merchantConfigDAO dao.MerchantConfigDAO,
	invoiceDAO dao.InvoiceDAO,
	staffDAO dao.StaffDAO,
	staffSLOA shopline.StaffOpenAPI,
	orderSLOA shopline.OrderOpenAPI,
	returnOrderSLOA shopline.ReturnOrderOpenAPI,
	storage storage.Storage,
	emailSvc email.EmailClient,
	timeout time.Duration,
	nowFunc func() time.Time,
	orderDeliveryOneship oneship.OrderDeliveryOneship,
	saleOrderOneship oneship.SaleOrderOneship,
	logisticsOrderOneship oneship.LogisticsOrderOneship,
	excelInvoiceGeneratorFactory excelgenerator.InvoiceERPExcelGeneratorFactory,
	excelSaleOrderGeneratorFactory excelgenerator.SaleOrderERPExcelGeneratorFactory) *exportsaleerpfile.ExportSaleERPFileJob {
	wire.Build(
		InitLogger,

		repository.NewClientRepository,
		repository.NewMerchantConfigRepository,
		repository.NewStaffRepository,
		repository.NewSaleOrderRepository,
		repository.NewERPOrderFileRepository,
		repository.NewInvoiceRepository,
		repository.NewERPInvoiceFileRepository,

		client.NewClientService,
		sale_order.NewSaleOrderService,
		erp_file.NewERPFileService,
		merchant_config.NewMerchantConfigService,
		invoice.NewInvoiceService,

		emailsvc.NewEmailService,

		exportsaleerpfile.NewExportSaleERPFileJob,
	)
	return new(exportsaleerpfile.ExportSaleERPFileJob)
}

func InitCronjobExportERPRefundFile(
	clientDAO dao.ClientDAO,
	saleOrderDAO dao.SaleOrderDAO,
	erpOrderFileDAO dao.ERPOrderFileDAO,
	erpInvoiceFileDAO dao.ERPInvoiceFileDAO,
	merchantConfigDAO dao.MerchantConfigDAO,
	staffDAO dao.StaffDAO,
	invoiceDAO dao.InvoiceDAO,
	staffSLOA shopline.StaffOpenAPI,
	orderSLOA shopline.OrderOpenAPI,
	returnOrderSLOA shopline.ReturnOrderOpenAPI,
	storage storage.Storage,
	emailSvc email.EmailClient,
	timeout time.Duration,
	nowFunc func() time.Time,
	orderDeliveryOneship oneship.OrderDeliveryOneship,
	logisticsOrderOneship oneship.LogisticsOrderOneship,
	saleOrderOneship oneship.SaleOrderOneship,
	excelInvoiceGeneratorFactory excelgenerator.InvoiceERPExcelGeneratorFactory,
	excelSaleOrderGeneratorFactory excelgenerator.SaleOrderERPExcelGeneratorFactory) *exportrefunderpfile.ExportRefundERPFileJob {
	wire.Build(
		InitLogger,

		repository.NewClientRepository,
		repository.NewMerchantConfigRepository,
		repository.NewStaffRepository,
		repository.NewSaleOrderRepository,
		repository.NewERPOrderFileRepository,
		repository.NewInvoiceRepository,
		repository.NewERPInvoiceFileRepository,

		client.NewClientService,
		sale_order.NewSaleOrderService,
		erp_file.NewERPFileService,
		merchant_config.NewMerchantConfigService,
		invoice.NewInvoiceService,
		emailsvc.NewEmailService,

		exportrefunderpfile.NewExportRefundERPFileJob,
	)
	return new(exportrefunderpfile.ExportRefundERPFileJob)
}

func InitSLWebhookOrderCancelConsumer(
	saleOrderDAO dao.SaleOrderDAO,
	clientDAO dao.ClientDAO,
	merchantDAO dao.MerchantConfigDAO,
	orderSLOA shopline.OrderOpenAPI,
	returnOrderSLOA shopline.ReturnOrderOpenAPI,
	orderDeliveryOneship oneship.OrderDeliveryOneship,
	saleOrderOneship oneship.SaleOrderOneship,
	logisticsOrderOneship oneship.LogisticsOrderOneship,
	emailSvc email.EmailClient,
	nowFunc func() time.Time) *order_cancel.SLWebhookOrderCancelConsumer {
	wire.Build(
		InitSQS,
		InitLogger,
		local.NewLocalStorage,

		repository.NewClientRepository,
		repository.NewMerchantConfigRepository,
		repository.NewSaleOrderRepository,

		emailsvc.NewEmailService,
		client.NewClientService,
		merchant_config.NewMerchantConfigService,
		sale_order.NewSaleOrderService,

		order_cancel.NewSLWebhookOrderCancelConsumer,
	)
	return new(order_cancel.SLWebhookOrderCancelConsumer)
}

func InitSLWebhookReturnOrderCreateConsumer(
	saleOrderDAO dao.SaleOrderDAO,
	clientDAO dao.ClientDAO,
	merchantDAO dao.MerchantConfigDAO,
	orderSLOA shopline.OrderOpenAPI,
	returnOrderSLOA shopline.ReturnOrderOpenAPI,
	orderDeliveryOneship oneship.OrderDeliveryOneship,
	saleOrderOneship oneship.SaleOrderOneship,
	logisticsOrderOneship oneship.LogisticsOrderOneship,
	emailSvc email.EmailClient,
	nowFunc func() time.Time) *return_order_create.SLWebhookReturnOrderCreateConsumer {
	wire.Build(
		InitSQS,
		InitLogger,
		local.NewLocalStorage,

		repository.NewClientRepository,
		repository.NewMerchantConfigRepository,
		repository.NewSaleOrderRepository,

		emailsvc.NewEmailService,
		client.NewClientService,
		merchant_config.NewMerchantConfigService,
		sale_order.NewSaleOrderService,

		return_order_create.NewSLWebhookReturnOrderCreateConsumer,
	)
	return new(return_order_create.SLWebhookReturnOrderCreateConsumer)
}

func InitCronjobExportERPInvoiceFile(
	clientDAO dao.ClientDAO,
	invoiceDAO dao.InvoiceDAO,
	erpOrderFileDAO dao.ERPOrderFileDAO,
	erpInvoiceFileDAO dao.ERPInvoiceFileDAO,
	merchantConfigDAO dao.MerchantConfigDAO,
	saleOrder dao.SaleOrderDAO,
	staffDAO dao.StaffDAO,
	staffSLOA shopline.StaffOpenAPI,
	orderSLOA shopline.OrderOpenAPI,
	returnOrderSLOA shopline.ReturnOrderOpenAPI,
	storage storage.Storage,
	emailSvc email.EmailClient,
	timeout time.Duration,
	nowFunc func() time.Time,
	orderDeliveryOneship oneship.OrderDeliveryOneship,
	logisticsOrderOneship oneship.LogisticsOrderOneship,
	saleOrderOneship oneship.SaleOrderOneship,
	excelSaleOrderGeneratorFactory excelgenerator.SaleOrderERPExcelGeneratorFactory,
	excelInvoiceGeneratorFactory excelgenerator.InvoiceERPExcelGeneratorFactory) *exportinvoiceerpfile.ExportERPInvoiceFileJob {
	wire.Build(
		InitLogger,

		repository.NewClientRepository,
		repository.NewMerchantConfigRepository,
		repository.NewStaffRepository,
		repository.NewSaleOrderRepository,
		repository.NewInvoiceRepository,
		repository.NewERPOrderFileRepository,
		repository.NewERPInvoiceFileRepository,

		client.NewClientService,
		sale_order.NewSaleOrderService,
		erp_file.NewERPFileService,
		merchant_config.NewMerchantConfigService,
		invoice.NewInvoiceService,
		emailsvc.NewEmailService,

		exportinvoiceerpfile.NewExportERPInvoiceFileJob,
	)
	return new(exportinvoiceerpfile.ExportERPInvoiceFileJob)
}

func InitWebhookInvoiceCreateConsumer(
	merchantConfigDAO dao.MerchantConfigDAO,
	saleOrderDAO dao.SaleOrderDAO,
	invoiceDAO dao.InvoiceDAO,
	orderSLOA shopline.OrderOpenAPI,
	returnOrderSLOA shopline.ReturnOrderOpenAPI,
	nowFunc func() time.Time,
) *invoiceevent.InvoiceCreateConsumer {
	wire.Build(
		InitLogger,
		InitSQS,
		repository.NewMerchantConfigRepository,
		repository.NewInvoiceRepository,
		repository.NewSaleOrderRepository,

		invoice.NewInvoiceService,
		merchant_config.NewMerchantConfigService,
		invoiceevent.NewInvoiceCreateConsumer,
	)
	return new(invoiceevent.InvoiceCreateConsumer)
}

func InitWebhookInvoiceUpdateConsumer(
	merchantConfigDAO dao.MerchantConfigDAO,
	saleOrderDAO dao.SaleOrderDAO,
	invoiceDAO dao.InvoiceDAO,
	orderSLOA shopline.OrderOpenAPI,
	returnOrderSLOA shopline.ReturnOrderOpenAPI,
	nowFunc func() time.Time,
) *invoiceevent.InvoiceUpdateConsumer {
	wire.Build(
		InitLogger,
		InitSQS,
		repository.NewMerchantConfigRepository,
		repository.NewInvoiceRepository,
		repository.NewSaleOrderRepository,

		invoice.NewInvoiceService,
		merchant_config.NewMerchantConfigService,
		invoiceevent.NewInvoiceUpdateConsumer,
	)
	return new(invoiceevent.InvoiceUpdateConsumer)
}

func InitSLWebhookOrderDeliveryUpdateConsumer(
	saleOrderDAO dao.SaleOrderDAO,
	clientDAO dao.ClientDAO,
	merchantDAO dao.MerchantConfigDAO,
	orderSLOA shopline.OrderOpenAPI,
	returnOrderSLOA shopline.ReturnOrderOpenAPI,
	orderDeliveryOneship oneship.OrderDeliveryOneship,
	logisticsOneship oneship.LogisticsOrderOneship,
	saleOrderOneship oneship.SaleOrderOneship,
	emailSvc email.EmailClient,
	nowFunc func() time.Time) *order_update.SLWebhookOrderDeliveryUpdateConsumer {
	wire.Build(
		InitSQS,
		InitLogger,
		local.NewLocalStorage,

		repository.NewClientRepository,
		repository.NewMerchantConfigRepository,
		repository.NewSaleOrderRepository,

		emailsvc.NewEmailService,
		client.NewClientService,
		merchant_config.NewMerchantConfigService,
		sale_order.NewSaleOrderService,

		order_update.NewSLWebhookOrderDeliveryUpdateConsumer,
	)
	return new(order_update.SLWebhookOrderDeliveryUpdateConsumer)
}

func InitSLWebhookMemberPointCreateConsumer(
	saleOrderDAO dao.SaleOrderDAO,
	clientDAO dao.ClientDAO,
	merchantDAO dao.MerchantConfigDAO,
	orderSLOA shopline.OrderOpenAPI,
	returnOrderSLOA shopline.ReturnOrderOpenAPI,
	orderDeliveryOneship oneship.OrderDeliveryOneship,
	logisticsOneship oneship.LogisticsOrderOneship,
	saleOrderOneship oneship.SaleOrderOneship,
	emailSvc email.EmailClient,
	nowFunc func() time.Time) *aws.SLWebhookMemberPointCreateConsumer {
	wire.Build(
		InitSQS,
		InitLogger,
		local.NewLocalStorage,
		repository.NewClientRepository,
		repository.NewMerchantConfigRepository,
		repository.NewSaleOrderRepository,
		emailsvc.NewEmailService,
		client.NewClientService,
		merchant_config.NewMerchantConfigService,
		sale_order.NewSaleOrderService,
		aws.NewSLWebhookMemberPointCreateConsumer,
	)
	return new(aws.SLWebhookMemberPointCreateConsumer)
}

func InitSLWebhookAppInstallationTokenCreateConsumer(
	merchantDAO dao.MerchantConfigDAO,
	nowFunc func() time.Time) *app_installation_token_create.AppInstallationTokenCreateConsumer {
	wire.Build(
		InitSQS,
		InitLogger,
		repository.NewMerchantConfigRepository,
		app_installation_token_create.NewAppInstallationTokenCreateConsumer,
	)
	return new(app_installation_token_create.AppInstallationTokenCreateConsumer)
}

func InitCronjobQueryOneshipLogisticsOrder(
	clientDAO dao.ClientDAO,
	saleOrderDAO dao.SaleOrderDAO,
	erpOrderFileDAO dao.ERPOrderFileDAO,
	erpInvoiceFileDAO dao.ERPInvoiceFileDAO,
	merchantConfigDAO dao.MerchantConfigDAO,
	staffDAO dao.StaffDAO,
	invoiceDAO dao.InvoiceDAO,
	staffSLOA shopline.StaffOpenAPI,
	orderSLOA shopline.OrderOpenAPI,
	returnOrderSLOA shopline.ReturnOrderOpenAPI,
	storage storage.Storage,
	timeout time.Duration,
	nowFunc func() time.Time,
	orderDeliveryOneship oneship.OrderDeliveryOneship,
	logisticsOrderOneship oneship.LogisticsOrderOneship,
	saleOrderOneship oneship.SaleOrderOneship,
	emailSvc email.EmailClient,
	excelInvoiceGeneratorFactory excelgenerator.InvoiceERPExcelGeneratorFactory,
	excelSaleOrderGeneratorFactory excelgenerator.SaleOrderERPExcelGeneratorFactory) *queryoneshiplogisticsorder.QueryOneshipLogisticsOrderJob {
	wire.Build(
		InitLogger,

		repository.NewClientRepository,
		repository.NewMerchantConfigRepository,
		repository.NewSaleOrderRepository,

		emailsvc.NewEmailService,
		client.NewClientService,
		sale_order.NewSaleOrderService,
		merchant_config.NewMerchantConfigService,

		queryoneshiplogisticsorder.NewQueryOneshipLogisticsOrderJob,
	)
	return new(queryoneshiplogisticsorder.QueryOneshipLogisticsOrderJob)
}
