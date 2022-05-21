package main

import (
	"bytes"
	"fmt"
	"os"

	dns "github.com/jlaw90/hetzner-dns"
)

func main() {
	key := os.Getenv("HETZNER_API_KEY")
	if key == "" {
		panic("API KEY MUST BE SET")
	}

	client := dns.NewClient(dns.ClientConfiguration{
		APIToken: key,
	})

	zones, err := client.GetZones(dns.GetZonesRequest{})
	fmt.Printf("GET ZONES: %v\n%+v\n", err, zones.Zones)

	zone, err := client.GetZone(zones.Zones[0].ID)

	fmt.Printf("GET ZONE: %v\n%+v\n", err, zone.Zone)

	export, err := client.ExportZone(zone.Zone.ID)

	fmt.Printf("EXPORT ZONE: %v\n%+v\n", err, string(export))

	result, err := client.ImportZone(zone.Zone.ID, bytes.NewBuffer(export))
	fmt.Printf("IMPORT ZONE: %v\n%+v\n", err, result)
}
