// Package dyndns provides a tool for running a
// dynamic dns updating service.
package dyndns

import (
	"fmt"
	"sync"
	"time"

	"github.com/zjstraus/name-dyndns/api"
	"github.com/zjstraus/name-dyndns/log"
)

var wg sync.WaitGroup

func contains(c api.Config, val string) bool {
	for _, v := range c.Hostnames {
		// We have a special case where an empty hostname
		// is equivalent to the domain (i.e. val == domain).
		if val == c.Domain && v == "" {
			return true
		} else if fmt.Sprintf("%s.%s.", v, c.Domain) == val {
			return true
		}
	}
	return false
}

func updateDNSRecord(a api.API, newRecord api.DNSRecord) error {
	log.Logger.Printf("Deleting DNS record for %s: %s\n", newRecord.Host, newRecord.DomainName)
	err := a.DeleteDNSRecord(newRecord.DomainName, newRecord.RecordID)
	if err != nil {
		return err
	}

	log.Logger.Printf("Creating DNS record for %s: %s\n", newRecord.Host, newRecord.DomainName)

	return a.CreateDNSRecord(newRecord)
}

func runConfig(c api.Config) {
	defer wg.Done()

	a := api.NewAPIFromConfig(c)
	for {
		ip, err := GetExternalIP()
		if err != nil {
			log.Logger.Print("Failed to retreive IP: ")
			log.Logger.Print(err)
			log.Logger.Printf("Will retry in %d seconds...\n", c.Interval)
			time.Sleep(time.Duration(c.Interval) * time.Second)
		}

		// GetRecords retrieves a list of DNSRecords,
		// 1 per hostname with the associated domain.
		// If the content is not the current IP, then
		// update it.
		records, err := a.GetDNSRecords(c.Domain)
		if err != nil {
			log.Logger.Printf("Failed to retreive records for %s:\n\t%s\n", c.Domain, err)
			log.Logger.Printf("Will retry in %d seconds...\n", c.Interval)
			time.Sleep(time.Duration(c.Interval) * time.Second)
			continue
		}

		for _, r := range records {
			log.Logger.Printf("Checking against %s", r.FQDN)
			if !contains(c, r.FQDN) {
				continue
			}

			// Only A records should be mapped to an IP.
			// TODO: Support AAAA records.
			if r.Type != "A" {
				continue
			}

			log.Logger.Printf("Running update check for %s.", r.Host)
			if r.Answer != ip {
				r.Answer = ip
				err = updateDNSRecord(a, r)
				if err != nil {
					log.Logger.Printf("Failed to update record %d [%s] with IP: %s\n\t%s\n", r.RecordID, r.Host, ip, err)
				} else {
					log.Logger.Printf("Updated record %d [%s] with IP: %s\n", r.RecordID, r.Host, ip)
				}
			}
		}

		log.Logger.Println("Update complete.")
		log.Logger.Printf("Will update again in %d seconds.\n", c.Interval)

		time.Sleep(time.Duration(c.Interval) * time.Second)
	}
}

// Run will process each configuration in configs.
// If daemon is true, then Run will run forever,
// processing each configuration at its specified
// interval.
//
// Each configuration represents a domain with
// multiple hostnames.
func Run(config api.Config) {
	runConfig(config)

	wg.Wait()
}
