package dns

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClientImpl_GetZones(t *testing.T) {
	hitServer := false
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		hitServer = true

		assert.Equal(t, "/zones?name=namey&page=5&per_page=27&search_name=searchy", request.URL.String())

		writer.Write([]byte("{\"meta\":{\"pagination\":{\"page\":1,\"per_page\":100,\"last_page\":1,\"total_entries\":5}},\"zones\":[{\"id\":\"sumtin\",\"name\":\"azone.uk\",\"created\":\"2022-05-20 22:40:44.522 +0000 UTC\",\"modified\":\"2022-05-20 22:40:45.85 +0000 UTC\",\"legacy_dns_host\":\"\",\"legacy_ns\":[\"name-servers.com.\",\"name-servers.com.\"],\"ns\":[\"hydrogen.ns.hetzner.com\",\"oxygen.ns.hetzner.com\",\"helium.ns.hetzner.de\"],\"owner\":\"\",\"paused\":false,\"permission\":\"\",\"project\":\"\",\"records_count\":19,\"registrar\":\"\",\"status\":\"verified\",\"ttl\":86400,\"verified\":\"\",\"is_secondary_dns\":false,\"txt_verification\":{\"name\":\"\",\"token\":\"\"}}]}"))
	}))
	defer server.Close()

	client := NewClient(ClientConfiguration{
		APIToken:   "some_token",
		apiUrl:     server.URL,
		HttpClient: server.Client(),
	})

	zones, err := client.GetZones(GetZonesRequest{
		SearchName: "searchy",
		Name:       "namey",
		PagedRequest: PagedRequest{
			Page:    5,
			PerPage: 27,
		},
	})

	assert.True(t, hitServer)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(zones.Zones))

	assert.Equal(t, GetZonesResponse{
		Metadata: PagedMetadata{
			Pagination: PageMetadata{
				Page:         1,
				PerPage:      100,
				LastPage:     1,
				TotalEntries: 5,
			},
		},
		Zones: []Zone{
			{
				ID:              "sumtin",
				Name:            "azone.uk",
				Created:         HetznerTime(time.Date(2022, 5, 20, 22, 40, 44, 522000000, time.UTC)),
				Modified:        HetznerTime(time.Date(2022, 5, 20, 22, 40, 45, 850000000, time.UTC)),
				LegacyDNSHost:   "",
				LegacyNS:        []string{"name-servers.com.", "name-servers.com."},
				NS:              []string{"hydrogen.ns.hetzner.com", "oxygen.ns.hetzner.com", "helium.ns.hetzner.de"},
				Owner:           "",
				Paused:          false,
				Permission:      "",
				Project:         "",
				RecordsCount:    19,
				Registrar:       "",
				Status:          ZoneVerified,
				TTL:             86400,
				Verified:        nil,
				IsSecondaryDNS:  false,
				TxtVerification: TxtVerification{},
			},
		},
	}, zones)
}
