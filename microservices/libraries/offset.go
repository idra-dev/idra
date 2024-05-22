package libraries

import (
	"strconv"
	"time"
)

type connectorIdType string

type IntOffset struct{
	ConnectorId connectorIdType
	Offset      int64
}


func (offset IntOffset) GetOffsetId(connectorId string) int64{
	value := offset.GetValue(connectorId)
	res, _ := strconv.ParseInt(string(value), 10, 64)
	return res
}

func (offset IntOffset) GetValue(connectorId string) []byte {
	value := GetKey("/offsets/" + connectorId)
	return value.Kvs[0].Value
}

func (IntOffset) SetOffsetId(connectorId string, offset int64) int64{
	return SetIntKey("/offsets/" + connectorId, offset)
}

func GetInt64FromTime(time time.Time) int64{
	return time.UnixMilli()
}

func GetTimeFromInt64(milliseconds int64) time.Time {
	return time.UnixMilli(milliseconds)
}