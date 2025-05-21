package startup

import (
	"fmt"
	"github.com/jw803/webook/pkg/mongox"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hash/fnv"
)

const MongoObjectIdPrefix = "65acced63a81ecae8d"

func NewGenerateObjectId(fixedCode, initObjectIdNumber int) string {
	return fmt.Sprintf("%s%03d%03d", MongoObjectIdPrefix, fixedCode, initObjectIdNumber)
}

func NewObjectIdWithTestNameFunc(testName string, initialNumber int) mongox.ObjectIdFunc {
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
