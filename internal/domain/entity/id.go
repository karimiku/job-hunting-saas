package entity

import "github.com/google/uuid"

type UserID = uuid.UUID
type CompanyID = uuid.UUID
type EntryID = uuid.UUID
type TaskID = uuid.UUID
type EvidenceID = uuid.UUID
type StageHistoryID = uuid.UUID
type CompanyAliasID = uuid.UUID

func NewID() uuid.UUID {
	return uuid.New()
}
