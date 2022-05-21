package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	ImportZone(id string, file io.Reader) (SingleZoneResponse, error)
	ExportZone(id string) ([]byte, error)
	ValidateZone(file io.Reader) (ValidateZoneResponse, error)
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

	_, err = c.request("GET", "zones", query, nil, &response)

	return response, err
}

func (c clientImpl) CreateZone(create WriteZoneRequest) (response SingleZoneResponse, err error) {
	encodedBody, err := json.Marshal(create)
	if err != nil {
		return response, err
	}
	_, err = c.request("POST", "zones", nil, bytes.NewBuffer(encodedBody), &response)

	return response, err
}

func (c clientImpl) GetZone(id string) (response SingleZoneResponse, err error) {
	_, err = c.request("GET", fmt.Sprintf("zones/%s", id), nil, nil, &response)
	return response, err
}

func (c clientImpl) UpdateZone(id string, update WriteZoneRequest) (response SingleZoneResponse, err error) {
	encodedBody, err := json.Marshal(update)
	if err != nil {
		return response, err
	}
	_, err = c.request("PATCH", fmt.Sprintf("zones/%s", id), nil, bytes.NewBuffer(encodedBody), &response)

	return response, err
}

func (c clientImpl) DeleteZone(id string) error {
	_, err := c.request("DELETE", fmt.Sprintf("zones/%s", id), nil, nil, nil)
	return err
}

func (c clientImpl) ImportZone(id string, file io.Reader) (response SingleZoneResponse, err error) {
	_, err = c.request("POST", fmt.Sprintf("zones/%s/import", id), nil, file, &response)
	return response, err
}

func (c clientImpl) ExportZone(id string) ([]byte, error) {
	resp, err := c.request("GET", fmt.Sprintf("zones/%s/export", id), nil, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (c clientImpl) ValidateZone(file io.Reader) (response ValidateZoneResponse, err error) {
	_, err = c.request("POST", "/zones/file/validate", nil, file, &response)
	return response, err
}
