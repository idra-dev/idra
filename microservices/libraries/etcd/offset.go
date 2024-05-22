package etcd

import (
	"microservices/libraries"
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
	if value==nil{
		return -1
	}
	res, _ := strconv.ParseInt(string(value), 10, 64)
	return res
}

func (offset IntOffset) GetValue(connectorId string) []byte {
	value := libraries.GetKey("/offsets/" + connectorId)
	if value.Kvs == nil{
		return nil
	}
	return value.Kvs[0].Value
}

func (IntOffset) SetOffsetId(connectorId string, offset int64) int64{
	return libraries.SetIntKey("/offsets/" + connectorId, offset)
}

func SetOffsetToken(connectorId string, token string){
	libraries.SaveKey("/offsets/" + connectorId, []byte(token))
}

func GetOffsetToken(connectorId string) string {
	value := libraries.GetKey("/offsets/" + connectorId)
	if value.Kvs == nil{
		return ""
	}
	return string(value.Kvs[0].Value)
}

func GetInt64FromTime(time time.Time) int64{
	return time.UnixMilli()
}

func GetTimeFromInt64(milliseconds int64) time.Time {
	return time.UnixMilli(milliseconds)
}