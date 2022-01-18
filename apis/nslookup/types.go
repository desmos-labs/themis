package nslookup

// TxtRecord represents a single TXT record for a specific domain
type TxtRecord struct {
	Text string `json:"text"`
}

// NewTxtRecord returns a new TxtRecord instance
func NewTxtRecord(text string) TxtRecord {
	return TxtRecord{
		Text: text,
	}
}

// LookupResponse represents a lookup response
type LookupResponse struct {
	TxtRecords []TxtRecord `json:"txt"`
}

// NewLookupResponse returns a new LookupResponse instance
func NewLookupResponse(txtRecords []TxtRecord) LookupResponse {
	return LookupResponse{
		TxtRecords: txtRecords,
	}
}
