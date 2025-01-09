package integrationtest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jw803/webook/config"
	"github.com/jw803/webook/internal/integration_test/startup"
	"github.com/jw803/webook/internal/integration_test/utils"
	"github.com/jw803/webook/internal/integration_test/utils/test_data"
	"github.com/jw803/webook/internal/pkg/ginx"
	"github.com/jw803/webook/pkg/loggerx"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/mock/gomock"
)

type WebAPIOrderFileUploadCommonSadTestSuite struct {
	suite.Suite
	mongoClient           *mongo.Client
	mdb                   *mongo.Database
	logger                loggerx.Logger
	orderFileUploadLogCol *mongo.Collection
	saleOrderCol          *mongo.Collection
	erpOrderFileCol       *mongo.Collection
	staffCol              *mongo.Collection
	clientCol             *mongo.Collection
	merchantConfigCol     *mongo.Collection
	email                 email.EmailClient
}

func (s *WebAPIOrderFileUploadCommonSadTestSuite) SetupSuite() {
	config.Init()
	s.logger = startup.InitLogger()

	s.mongoClient = startup.InitMongoDBClient()
	s.mdb = startup.InitMongoDB(s.mongoClient)
	s.staffCol = s.mdb.Collection(staffmongo.TableName)
	s.clientCol = s.mdb.Collection(clientmongo.TableName)
	s.merchantConfigCol = s.mdb.Collection(merchantconfigmongo.TableName)
	s.orderFileUploadLogCol = s.mdb.Collection(orderfileuploadlogmongo.TableName)
	s.saleOrderCol = s.mdb.Collection(saleordermongo.TableName)
	s.erpOrderFileCol = s.mdb.Collection(erpsalefilemongo.TableName)
	s.email = local.NewLocalEmailService()
	ginx.InitErrorCounter(prometheus.CounterOpts{
		Name: "test_webapi_order_file_upload_common_sad_path",
		Help: "pod-exportSaleERPFileJob.yaml.tftpl error metrics",
		ConstLabels: map[string]string{
			"app":  "cst-wstyle-integration",
			"lang": "golang",
		},
	})
	ginx.SetLogger(s.logger)
	response.SetLogger(s.logger)
}
func (s *WebAPIOrderFileUploadCommonSadTestSuite) TearDownSuite() {
	s.ClearSeedData()
}

func TestWebAPIUploadFileCommonSadPathHandler(t *testing.T) {
	suite.Run(t, new(WebAPIOrderFileUploadCommonSadTestSuite))
}

func (s *WebAPIOrderFileUploadCommonSadTestSuite) ClearSeedData() {
}

func (s *WebAPIOrderFileUploadCommonSadTestSuite) ClearTestCaseData() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, err := s.orderFileUploadLogCol.DeleteMany(ctx, bson.D{})
	assert.NoError(s.T(), err)
	_, err = s.saleOrderCol.DeleteMany(ctx, bson.D{})
	assert.NoError(s.T(), err)
	_, err = s.staffCol.DeleteMany(ctx, bson.D{})
	assert.NoError(s.T(), err)
	_, err = s.clientCol.DeleteMany(ctx, bson.D{})
	assert.NoError(s.T(), err)
	_, err = s.merchantConfigCol.DeleteMany(ctx, bson.D{})
	assert.NoError(s.T(), err)
}

func (s *WebAPIOrderFileUploadCommonSadTestSuite) TestWebAPIOrderFileUpload() {
	t := s.T()
	testCases := []struct {
		name            string
		before          func(t *testing.T)
		mock            func(t *testing.T, ctrl *gomock.Controller)
		after           func(t *testing.T)
		reqFormFilePath string
		reqFileName     string
		reqBody         string
		wantCode        int
		wantResBody     string
	}{
		{
			name: "sent successfully",
			mock: func(t *testing.T, ctrl *gomock.Controller) {
				return
			},
			before: func(t *testing.T) {
			},
			after:           func(t *testing.T) {},
			reqFormFilePath: "files/sales/others/normal.xlsx",
			reqFileName:     "normal.xlsx",
			reqBody: `
				"phone": ""
			`,
			wantCode: 200,
			wantResBody: `
				"code": 0,
				"data": null,
				"msg": ""
			`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r != nil {
					t.Log("Panic arose! Recovered in ", r)
					t.Fail()
				}
				s.ClearTestCaseData()
			}()
			ctrl := gomock.NewController(t)
			tc.before(t)

			now, objectIdGen, staffSLOAMock, orderSLOAMock, returnOrderSLOAMock, emailMock,
				orderDeliveryOneShipMock, saleOrderOneShipMock, logisticsOrderOneshipMock := tc.mock(t, ctrl)

			server := s.setServer(now, objectIdGen, staffSLOAMock, orderSLOAMock, returnOrderSLOAMock, emailMock,
				orderDeliveryOneShipMock, saleOrderOneShipMock, logisticsOrderOneshipMock)

			file, err := os.Open(tc.reqFormFilePath)
			if err != nil {
				panic(err.Error())
				return
			}
			defer file.Close()
			fileBuffer := bytes.Buffer{}
			_, err = fileBuffer.ReadFrom(file)
			if err != nil {
				fmt.Println("Error reading file:", err)
				return
			}

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			part, err := writer.CreateFormFile("file", tc.reqFileName)
			assert.NoError(t, err)
			_, err = part.Write(fileBuffer.Bytes())
			assert.NoError(t, err)
			jsonPart, err := writer.CreateFormField("body")
			assert.NoError(t, err)
			formBodyJsonString, err := json.Marshal(tc.reqBody(t))
			assert.NoError(t, err)
			jsonPart.Write(formBodyJsonString)

			writer.Close()

			req, err := http.NewRequest(http.MethodPost,
				"/order-file/upload", body)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)
			code := recorder.Code

			var result utils.Result[string]
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.Equal(t, tc.wantCode, code)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResBody(t), result)
			if code != http.StatusOK {
				if result.Msg != "" {
					t.Log("API Error & Response::", recorder.Body.Bytes())
				}
				return
			}
			tc.after(t)
		})
	}
}

func (s *WebAPIOrderFileUploadCommonSadTestSuite) setServer(
	nowFunc utils.NowFunc, objectIdFunc utils.ObjectIdFunc,
	staffSLOAMock shopline.StaffOpenAPI, orderSLOAMock shopline.OrderOpenAPI,
	returnOrderSLOAMock shopline.ReturnOrderOpenAPI,
	emailSvcMock email.EmailClient, orderDeliveryOSMock oneship.OrderDeliveryOneship,
	saleOrderOSMock oneship.SaleOrderOneship, logisticsOrderOSMock oneship.LogisticsOrderOneship) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("shopline", ginx.ShoplineClaims{
			Data: ginx.Staff{
				MerchantId: test_data.TestMerchantConfigWstyle.MerchantId,
				StaffId:    test_data.TestStaff1.Id.Hex(),
				ClientId:   test_data.TestClientECShopWstyle.Id.Hex(),
			},
		})
		ctx.Next()
	})
	hdl := startup.InitOrderFileHandler(
		saleordermongo.NewSaleOrderDAO(s.mongoClient, s.mdb, s.logger, nowFunc, objectIdFunc),
		orderfileuploadlogmongo.NewOrderFileUploadLogDAO(s.mdb, s.logger, nowFunc, objectIdFunc),
		merchantconfigmongo.NewMerchantConfigDAO(s.mdb, s.logger, nowFunc, objectIdFunc),
		staffmongo.NewStaffDAO(s.mdb, s.logger, nowFunc, objectIdFunc),
		clientmongo.NewClientDAO(s.mdb, s.logger, nowFunc, objectIdFunc),
		staffSLOAMock,
		orderSLOAMock,
		returnOrderSLOAMock,
		emailSvcMock,
		orderDeliveryOSMock,
		saleOrderOSMock,
		logisticsOrderOSMock,
		nowFunc)
	hdl.RegisterRoutes(server)
	return server
}
