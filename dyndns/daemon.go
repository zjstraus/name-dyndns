package dyndns

import (
	"github.com/mfycheng/name-dyndns/api"
	"github.com/mfycheng/name-dyndns/log"
	"sync"
	"time"
)

var wg sync.WaitGroup

func updateDNSRecord(a api.API, domain, recordId string, newRecord api.DNSRecord) error {
	log.Logger.Printf("Deleting DNS record for %s.%s.\n", newRecord.Name, domain)
	err := a.DeleteDNSRecord(domain, newRecord.RecordId)
	if err != nil {
		return err
	}

	log.Logger.Printf("Creating DNS record for %s.%s: %s\n", newRecord.Name, domain, newRecord)
	return a.CreateDNSRecord(domain, newRecord)
}

func runConfig(c api.Config, daemon bool) {
	defer wg.Done()

	a := api.NewAPIFromConfig(c)
	for {
		log.Logger.Printf("Running update check for %s.", c.Domain)
		ip, err := GetExternalIP()
		if err != nil {
			log.Logger.Print("Failed to retreive IP: ")
			if daemon {
				log.Logger.Println("Will retry...")
				continue
			} else {
				log.Logger.Println("Giving up.")
				break
			}
		}

		// GetRecords retrieves a list of DNSRecords,
		// 1 per hostname with the associated domain.
		// If the content is not the current IP, then
		// update it.
		records, err := a.GetRecords(c.Domain)
		if err != nil {
			log.Logger.Printf("Failed to retreive records for%s\n", c.Domain)
			if daemon {
				log.Logger.Print("Will retry...")
				continue
			} else {
				log.Logger.Print("Giving up.")
				break
			}
		}

		for _, r := range records {
			if r.Content != ip {
				r.Content = ip
				err = updateDNSRecord(a, c.Domain, r.RecordId, r)
				if err != nil {
					log.Logger.Printf("Failed to update record %s [%s.%s] with IP: %s\n\t%s\n", r.RecordId, r.Name, c.Domain, ip, err)
				} else {
					log.Logger.Printf("Attempting to update record %s [%s.%s] with IP: %s\n", r.RecordId, r.Name, c.Domain, ip)
				}
			}
		}

		if !daemon {
			log.Logger.Println("Non daemon mode, stopping.")
			return
		}

		time.Sleep(time.Duration(c.Interval) * time.Second)
	}
}

// For each domain, check if the host record matches
// the current external IP. If it does not, it updates.
// If daemon is true, then Run will run forever, polling at
// an interval specified in each config.
func Run(configs []api.Config, daemon bool) {
	for _, config := range configs {
		wg.Add(1)
		go runConfig(config, daemon)
	}

	wg.Wait()
}
