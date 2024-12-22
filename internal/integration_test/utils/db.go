package utils

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hash/fnv"
)

const MongoObjectIdPrefix = "65acced63a81ecae8d"

var ObjectIdNumber = 0

func NewGenerateObjectIdClosureFunc(fixedCode, initObjectIdNumber int) func() primitive.ObjectID {
	return func() primitive.ObjectID {
		objectID, err := primitive.ObjectIDFromHex(NewGenerateObjectId(fixedCode, initObjectIdNumber))
		if err != nil {
			objectID = primitive.ObjectID{}
		}
		initObjectIdNumber++
		return objectID
	}
}

func NewGenerateObjectId(fixedCode, initObjectIdNumber int) string {
	return fmt.Sprintf("%s%03d%03d", MongoObjectIdPrefix, fixedCode, initObjectIdNumber)
}

type ObjectIdFunc = func() primitive.ObjectID

func NewGenerateObjectIdWithTestNameFunc(testName string, initialNumber int) ObjectIdFunc {
	return func() primitive.ObjectID {
		objectID := GetObjectIdByTestNameAndNumber(testName, initialNumber)
		initialNumber++
		return objectID
	}
}

func GetObjectIdByTestNameAndNumber(testName string, number int) primitive.ObjectID {
	objectIdTemplate := testNameToObjectId(testName)
	prefix := objectIdTemplate[:19]
	suffix := fmt.Sprintf("%05d", number)
	objectID, err := primitive.ObjectIDFromHex(prefix + suffix)
	if err != nil {
		objectID = primitive.ObjectID{}
	}
	return objectID
}

func testNameToObjectId(testName string) string {
	hash := fnv.New64a()
	hash.Write([]byte(testName))
	hashValue := hash.Sum64()

	var idBytes [12]byte
	for i := range idBytes {
		idBytes[i] = byte((hashValue >> uint((i%8)*8)) & 0xFF)
	}

	objectId := primitive.ObjectID(idBytes[:])
	return objectId.Hex()
}
