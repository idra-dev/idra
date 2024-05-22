package data

import (
	"github.com/antrad1978/cdc_shared"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"microservices/libraries/custom_errors"
	"microservices/libraries/models"
)

func SaveData(connector cdc_shared.Connector, rows []map[string]interface{}, db *gorm.DB) int {
	i := 0
	tx := db.Begin()
	fieldsUpdate := UpdateColumns(connector, rows)
	for ; i < len(rows); i++ {
		row := rows[i]
		if connector.SaveMode == models.Insert {
			tx := db.Table(connector.Table).Create(row)
			if tx.Error != nil {
				custom_errors.CdcLog(connector,tx.Error)
			}
		} else {
			tx = db.Table(connector.Table).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: connector.IdField}},
				DoUpdates: clause.AssignmentColumns(fieldsUpdate),
			}).Create(row)
			if tx.Error != nil {
				custom_errors.CdcLog(connector, tx.Error)
				tx.Rollback()
			}
		}
	}
	if i>0{
		db.Commit()
	}
	return i
}

func UpdateColumns(connector cdc_shared.Connector, rows []map[string]interface{}) []string {
	var fieldsUpdate []string
	if len(rows) > 0 {
		sampleRow := rows[0]
		for key := range sampleRow {
			if key != connector.IdField {
				fieldsUpdate = append(fieldsUpdate, key)
			}
		}
	}
	return fieldsUpdate
}

