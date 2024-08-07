package scan

import (
	"fmt"
	"log"
	"time"

	"github.com/aceberg/WatchYourLAN/internal/db"
	"github.com/aceberg/WatchYourLAN/internal/models"
	"github.com/aceberg/WatchYourLAN/internal/notify"
)

func hostsCompare() {

	for _, oneHost := range dbHosts {

		host, exists := foundHostsMap[oneHost.Mac]

		if exists && (appConfig.IgnoreIP == "yes" || host.IP == oneHost.IP) {

			oneHost.Date = host.Date

			if oneHost.Now == 0 {
				histAdd(oneHost, 1)
			}

			oneHost.Now = 1

			delete(foundHostsMap, oneHost.Mac)

			db.Update(appConfig.DbPath, oneHost)

		} else if oneHost.Now == 1 {
			oneHost.Now = 0
			oneHost.Date = time.Now().Format("2006-01-02 15:04:05")

			db.Update(appConfig.DbPath, oneHost)

			histAdd(oneHost, 0)
		}
	}

	for _, oneHost := range foundHostsMap {

		msg := fmt.Sprintf("💻 ACHTUNG Unbekannter Host. IP: '%s', MAC: '%s', Hw: '%s'", oneHost.IP, oneHost.Mac, oneHost.Hw)
		log.Println("WARN:", msg)
		notify.Shoutrrr(msg, appConfig.ShoutURL) // Notify through Shoutrrr

		db.Insert(appConfig.DbPath, oneHost)

		histAdd(oneHost, 1)
	}
}

func histAdd(oneHost models.Host, state int) {
	var history models.History

	if oneHost.ID == 0 {
		dbHosts = db.Select(appConfig.DbPath)

		for _, host := range dbHosts {
			if host.IP == oneHost.IP && host.Mac == oneHost.Mac {
				oneHost.ID = host.ID
				break
			}
		}
	}

	history.Host = oneHost.ID
	history.Name = oneHost.Name
	history.IP = oneHost.IP
	history.Mac = oneHost.Mac
	history.Hw = oneHost.Hw
	history.Date = oneHost.Date
	history.Known = oneHost.Known
	history.State = state

	db.InsertHist(appConfig.DbPath, history)
}
