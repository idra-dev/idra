package tests

import (
	"fmt"
	"github.com/antrad1978/cdc_shared"
	data2 "microservices/libraries/data"
	"testing"
)

func TestConnectorS3(t *testing.T){
	manager := data2.MssqlManager{}
	connector := cdc_shared.Connector{}
	connector.ConnectionString = "sqlserver://extdev:@linkfx-sit.database.windows.net?database=TradingSystemProdBackup"
	connector.Table = "tbl_Account_Charges"
	connector.IdField = "n_Account_Charges_Id"
	rows, res := manager.GetRowsById(connector, 0)
	fmt.Println(res, len(rows))
	manager2 := data2.S3JsonConnector{}
	manager2.ConnectionString = "eu-west-1"
	manager2.SaveData( "eu-west-1", "samplengt1", rows, "tbl_Account_Charges")
}