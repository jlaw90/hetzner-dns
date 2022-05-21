package dns

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"
)

type TxtVerification struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type zoneStatus string

const (
	ZonePending  zoneStatus = "pending"
	ZoneVerified zoneStatus = "verified"
	zoneFailed   zoneStatus = "failed"
)

type Zone struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Created         HetznerTime     `json:"created"`
	Modified        HetznerTime     `json:"modified"`
	LegacyDNSHost   string          `json:"legacy_dns_host"`
	LegacyNS        []string        `json:"legacy_ns"`
	NS              []string        `json:"ns"`
	Owner           string          `json:"owner"`
	Paused          bool            `json:"paused"`
	Permission      string          `json:"permission"`
	Project         string          `json:"project"`
	RecordsCount    uint64          `json:"records_count"`
	Registrar       string          `json:"registrar"`
	Status          zoneStatus      `json:"status"`
	TTL             uint64          `json:"ttl"`
	Verified        *HetznerTime    `json:"verified,omitempty"`
	IsSecondaryDNS  bool            `json:"is_secondary_dns"`
	TxtVerification TxtVerification `json:"txt_verification"`
}

func (z *Zone) UnmarshalJSON(data []byte) error {
	type zoneAlias Zone
	za := zoneAlias{}
	if err := json.Unmarshal(data, &za); err != nil {
		return err
	}
	if za.Verified != nil && time.Time(*(za.Verified)).IsZero() {
		za.Verified = nil
	}

	*z = Zone(za)

	return nil
}

type GetZonesRequest struct {
	PagedRequest

	Name       string
	SearchName string
}

type GetZonesResponse struct {
	Metadata PagedMetadata `json:"meta"`
	Zones    []Zone        `json:"zones"`
}

type WriteZoneRequest struct {
	Name string  `json:"name"`
	TTL  *uint64 `json:"ttl,omitempty"`
}

type SingleZoneResponse struct {
	Zone Zone `json:"zone"`
}

type ValidateZoneResponse struct {
	ParsedRecords int      `json:"parsed_records"`
	ValidRecords  []Record `json:"valid_records"`
}

type ZoneAPI interface {
	GetZones(GetZonesRequest) (GetZonesResponse, error)
	CreateZone(WriteZoneRequest) (SingleZoneResponse, error)
	GetZone(id string) (SingleZoneResponse, error)
	UpdateZone(id string, update WriteZoneRequest) (SingleZoneResponse, error)
	DeleteZone(id string) error
	ImportZone(id string, file io.ReadCloser) (SingleZoneResponse, error)
	ExportZone(id string) (string, error)
	ValidateZone(file io.ReadCloser) (ValidateZoneResponse, error)
}

func (c clientImpl) GetZones(req GetZonesRequest) (response GetZonesResponse, err error) {
	query := url.Values{}
	if req.Name != "" {
		query.Add("name", req.Name)
	}
	if req.SearchName != "" {
		query.Add("search_name", req.SearchName)
	}
	addPagedQueryParams(query, req.PagedRequest)

	err = c.request("GET", "zones").AddQueryParams(query).ReadJSON(&response)
	return response, err
}

func (c clientImpl) CreateZone(create WriteZoneRequest) (response SingleZoneResponse, err error) {
	err = c.request("POST", "zones").JSON(create, &response)
	return response, err
}

func (c clientImpl) GetZone(id string) (response SingleZoneResponse, err error) {
	err = c.request("GET", fmt.Sprintf("zones/%s", id)).ReadJSON(&response)
	return response, err
}

func (c clientImpl) UpdateZone(id string, update WriteZoneRequest) (response SingleZoneResponse, err error) {
	err = c.request("PATCH", fmt.Sprintf("zones/%s", id)).JSON(update, &response)
	return response, err
}

func (c clientImpl) DeleteZone(id string) error {
	_, err := c.request("DELETE", fmt.Sprintf("zones/%s", id)).Send()
	return err
}

func (c clientImpl) ImportZone(id string, file io.ReadCloser) (response SingleZoneResponse, err error) {
	err = c.request("POST", fmt.Sprintf("zones/%s/import", id)).WritePlain(file).ReadJSON(&response)
	return response, err
}

func (c clientImpl) ExportZone(id string) (string, error) {
	return c.request("GET", fmt.Sprintf("zones/%s/export", id)).ReadPlain()
}

func (c clientImpl) ValidateZone(file io.ReadCloser) (response ValidateZoneResponse, err error) {
	err = c.request("POST", "/zones/file/validate").WritePlain(file).ReadJSON(&response)
	return response, err
}
