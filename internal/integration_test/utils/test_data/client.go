package test_data

import (
	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/repository/dao"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var TestClientIdECShopWstyle, _ = primitive.ObjectIDFromHex("65acced63a81ecae8d000001")
var TestClientECShopWstyle = dao.Client{
	Id:          TestClientIdECShopWstyle,
	Name:        "Shopline_Wstyle",
	SerialId:    "6100",
	Source:      TestMerchantIdWstyle.Hex(),
	SourceId:    TestMerchantIdWstyle.Hex(),
	WarehouseId: "B201",
}

var TestClientIdECShopUNT, _ = primitive.ObjectIDFromHex("65acced63a81ecae8d000002")
var TestClientECShopUNT = dao.Client{
	Id:          TestClientIdECShopUNT,
	Name:        "Shopline_UNT",
	SerialId:    "6000",
	Source:      TestMerchantIdUNT.Hex(),
	SourceId:    TestMerchantIdUNT.Hex(),
	WarehouseId: "B201",
}

var TestClientIdOthers, _ = primitive.ObjectIDFromHex("65acced63a81ecae8d000003")
var TestClientOthers = dao.Client{
	Id:          TestClientIdOthers,
	Name:        "Others",
	SerialId:    "others_id_001",
	Source:      "",
	WarehouseId: "",
}

var TestClientIdOthers2, _ = primitive.ObjectIDFromHex("65acced63a81ecae8d000004")
var TestClientOthers2 = dao.Client{
	Id:          TestClientIdOthers2,
	Name:        "Others",
	SerialId:    "client_001",
	Source:      "",
	WarehouseId: "",
}
