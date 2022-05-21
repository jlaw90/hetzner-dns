package dns

type recordType string

const (
	A     recordType = "A"
	AAAA  recordType = "AAAA"
	PTR   recordType = "PTR"
	NS    recordType = "NS"
	MX    recordType = "MX"
	CNAME recordType = "CNAME"
	RP    recordType = "RP"
	TXT   recordType = "TXT"
	SOA   recordType = "SOA"
	HINFO recordType = "HINFO"
	SRV   recordType = "SRV"
	DANE  recordType = "DANE"
	TLSA  recordType = "TLSA"
	DS    recordType = "DS"
	CAA   recordType = "CAA"
)

type Record struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Created  HetznerTime `json:"created"`
	Modified HetznerTime `json:"modified"`
	ZoneID   string      `json:"zone_id"`
	Value    string      `json:"value"`
	TTL      uint64      `json:"ttl"`
	Type     recordType  `json:"type"`
}

type GetRecordsRequest struct {
	ZoneID string
}

type GetRecordsResponse struct {
	Records []Record `json:"records"`
}

type WriteRecordRequest struct {
	Name   string     `json:"name"`
	TTL    *uint64    `json:"ttl,omitempty"`
	Type   recordType `json:"type"`
	Value  string     `json:"value"`
	ZoneID string     `json:"zone_id"`
}

type SingleRecordResponse struct {
	Record Record `json:"record"`
}

type BulkCreateRecordsResponse struct {
	Records        []Record `json:"record"`
	ValidRecords   []Record `json:"valid_records"`
	InvalidRecords []Record `json:"invalid_records"`
}

type BulkUpdateRecord struct {
	ID string
	WriteRecordRequest
}

type BulkUpdateResponse struct {
	Records       []Record `json:"records"`
	FailedRecords []Record `json:"failed_records"`
}

type RecordAPI interface {
	GetRecords(GetRecordsRequest) (GetRecordsResponse, error)
	CreateRecord(WriteRecordRequest) (SingleRecordResponse, error)
	GetRecord(id string) (SingleRecordResponse, error)
	UpdateRecord(id string, update WriteRecordRequest) (SingleRecordResponse, error)
	DeleteRecord(id string) error
	CreateRecords([]WriteRecordRequest) (BulkCreateRecordsResponse, error)
	UpdateRecords([]BulkUpdateRecord) (BulkUpdateResponse, error)
}
