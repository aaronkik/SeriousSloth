package dynamodbstream

import (
	"github.com/aws/aws-lambda-go/events"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func ToAttributeValue(av events.DynamoDBAttributeValue) ddbtypes.AttributeValue {
	switch av.DataType() {
	case events.DataTypeString:
		return &ddbtypes.AttributeValueMemberS{Value: av.String()}
	case events.DataTypeNumber:
		return &ddbtypes.AttributeValueMemberN{Value: av.Number()}
	case events.DataTypeBoolean:
		return &ddbtypes.AttributeValueMemberBOOL{Value: av.Boolean()}
	case events.DataTypeBinary:
		return &ddbtypes.AttributeValueMemberB{Value: av.Binary()}
	case events.DataTypeNull:
		return &ddbtypes.AttributeValueMemberNULL{Value: true}
	case events.DataTypeList:
		list := av.List()
		out := make([]ddbtypes.AttributeValue, len(list))
		for i, value := range list {
			out[i] = ToAttributeValue(value)
		}
		return &ddbtypes.AttributeValueMemberL{Value: out}
	case events.DataTypeMap:
		return &ddbtypes.AttributeValueMemberM{Value: ToAttributeValueMap(av.Map())}
	case events.DataTypeStringSet:
		return &ddbtypes.AttributeValueMemberSS{Value: av.StringSet()}
	case events.DataTypeNumberSet:
		return &ddbtypes.AttributeValueMemberNS{Value: av.NumberSet()}
	case events.DataTypeBinarySet:
		return &ddbtypes.AttributeValueMemberBS{Value: av.BinarySet()}
	}
	return nil
}

func ToAttributeValueMap(m map[string]events.DynamoDBAttributeValue) map[string]ddbtypes.AttributeValue {
	out := make(map[string]ddbtypes.AttributeValue, len(m))
	for k, v := range m {
		out[k] = ToAttributeValue(v)
	}
	return out
}
