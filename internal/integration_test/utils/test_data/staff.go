package test_data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"bitbucket.org/starlinglabs/cst-wstyle-integration/internal/repository/dao"
)

var TestStaff1Id, _ = primitive.ObjectIDFromHex("65acced63a81ecae8d000001")

var TestStaff1 = dao.Staff{
	Id:    TestStaff1Id,
	Name:  "staff_wstyle",
	Email: "staff_wstyle@mail.com",
	Permission: map[string]dao.Permission{
		TestMerchantIdWstyle.Hex(): {
			IsEnableExport:           true,
			IsEnableImport:           true,
			IsEnableMailNotification: true,
		},
	},
}

var TestStaff2Id, _ = primitive.ObjectIDFromHex("65acced63a81ecae8d000002")

var TestStaff2 = dao.Staff{
	Id:    TestStaff2Id,
	Name:  "staff_unt",
	Email: "staff_unt@mail.com",
	Permission: map[string]dao.Permission{
		TestMerchantIdUNT.Hex(): {
			IsEnableExport:           true,
			IsEnableImport:           true,
			IsEnableMailNotification: true,
		},
	},
}

var TestStaff3Id, _ = primitive.ObjectIDFromHex("65acced63a81ecae8d000003")

var TestStaff3NoPermission = dao.Staff{
	Id:    TestStaff3Id,
	Name:  "staff_wstyle",
	Email: "staff_wstyle@mail.com",
	Permission: map[string]dao.Permission{
		TestMerchantIdWstyle.Hex(): {
			IsEnableExport:           false,
			IsEnableImport:           false,
			IsEnableMailNotification: false,
		},
	},
}
