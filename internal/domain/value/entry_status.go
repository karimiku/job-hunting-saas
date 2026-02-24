package value

import "errors"

var (
	ErrEntryStatusEmpty   = errors.New("entry status must not be empty")
	ErrEntryStatusInvalid = errors.New("entry status is invalid")
)

const (
	entryStatusInProgress = "in_progress"
	entryStatusOffered    = "offered"
	entryStatusAccepted   = "accepted"
	entryStatusRejected   = "rejected"
	entryStatusWithdrawn  = "withdrawn"
)

var validEntryStatuses = map[string]bool{
	entryStatusInProgress: true,
	entryStatusOffered:    true,
	entryStatusAccepted:   true,
	entryStatusRejected:   true,
	entryStatusWithdrawn:  true,
}

type EntryStatus struct {
	value string
}

func NewEntryStatus(raw string) (EntryStatus, error) {
	if raw == "" {
		return EntryStatus{}, ErrEntryStatusEmpty
	}
	if !validEntryStatuses[raw] {
		return EntryStatus{}, ErrEntryStatusInvalid
	}
	return EntryStatus{value: raw}, nil
}

func (s EntryStatus) String() string {
	return s.value
}

func (s EntryStatus) Equals(other EntryStatus) bool {
	return s.value == other.value
}

func (s EntryStatus) IsOpen() bool {
	return s.value == entryStatusInProgress || s.value == entryStatusOffered
}
