package tests

import (
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"microservices/libraries/data"
	"os"
	"testing"
)


var SqlServerConnectionString = os.Getenv("DB_SQLSERVER")


func TestMssqlMax(t *testing.T){
	manager := data.MssqlManager{}
	connector := cdc_shared.Connector{}
	connector.Table = "tbl_Account_Charges"
	connector.ConnectionString = SqlServerConnectionString
	connector.IdField = "n_Account_Charges_Id"
	id := manager.GetMaxTableId(connector)
	fmt.Println(id)
}

func TestMssqlMaxTimestamp(t *testing.T){
	manager := data.MssqlManager{}
	connector := cdc_shared.Connector{}
	connector.ConnectionString = SqlServerConnectionString
	connector.Table = "tbl_Account_Charges"
	connector.TimestampField = "d_Created_Date"
	ts,_ := manager.GetMaxTimestamp(connector)
	fmt.Println(" "+ts.String())
}

func TestMssqlDBConnector(t *testing.T){
	manager := data.MssqlManager{}
	connector := cdc_shared.Connector{}
	connector.Table = "tbl_Account_Charges"
	connector.IdField = "n_Account_Charges_Id"
	connector.ConnectionString = SqlServerConnectionString
	rows, res := manager.GetRowsById(connector, 0)
	fmt.Println(res, len(rows))
	//manager2 := data.MssqlManager{}
	//dsn2 := "sqlserver://tonio:cacatone123_@localhost?database=trading"
	//manager2.ConnectionString = dsn2
	//manager2.DatabaseType = "Mssql"
	//rows2 := manager2.GetRecords("dbo.tbl_Account_Charges", "n_Account_Charges_Id", 0)
	//fmt.Println(len(rows2))
	//manager2.InsertRows("tbl_Account_Charges",rows)
}
