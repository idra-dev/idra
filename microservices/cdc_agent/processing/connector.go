package processing

import (
	"encoding/json"
	"fmt"
	"github.com/antrad1978/cdc_shared"
	"microservices/libraries/custom_errors"
	"microservices/libraries/data"
	"microservices/libraries/models"
	"os"
	"time"
)


func ProcessStatic() {
	var syncs []cdc_shared.Sync
	path, err := os.Getwd()
	custom_errors.LogAndDie(err)
	fmt.Println(path)
	staticFilePath := os.Getenv(models.StaticFilePath)
	dat, err2 := os.ReadFile(staticFilePath)
	custom_errors.LogAndDie(err2)

	json.Unmarshal(dat, &syncs)
	for {
		for _, sync := range syncs {
			fmt.Println(sync.SyncName)
			data.SyncData(sync, sync.Mode)
		}
		time.Sleep(1 * 3600 * time.Second)
	}
}