package models

import "time"

// AuditRecord is the base object for audit records
type AuditRecord struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	IdentitySub string    `json:"identitySub"`
	ActionType  string    `json:"actionType"`
	RequestURL  string    `json:"requestUrl"`
	Payload     string    `json:"payload"`
	Action      string    `json:"action"`
	TimeStamp   time.Time `json:"timeStamp"`
}
