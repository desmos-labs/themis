package nslookup

import "net"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

// ReadTxtRecords reads the TXT records for the given domain
func (h *Handler) ReadTxtRecords(domain string) ([]TxtRecord, error) {
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		return nil, err
	}

	var records = make([]TxtRecord, len(txtRecords))
	for i, record := range txtRecords {
		records[i] = NewTxtRecord(record)
	}

	return records, nil
}
