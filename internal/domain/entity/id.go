package entity

import "github.com/google/uuid"

// Named types for entity IDs.
// type alias（=）ではなく named type を使うことで、
// 異なるエンティティの ID を取り違えた場合にコンパイルエラーになる。
// 例: func GetEntry(userID UserID, entryID EntryID) に対して
//
//	GetEntry(entryID, userID) はコンパイルエラーになる。
type UserID uuid.UUID
type CompanyID uuid.UUID
type EntryID uuid.UUID
type TaskID uuid.UUID
type StageHistoryID uuid.UUID
type CompanyAliasID uuid.UUID

// --- Constructors ---

func NewUserID() UserID               { return UserID(uuid.New()) }
func NewCompanyID() CompanyID          { return CompanyID(uuid.New()) }
func NewEntryID() EntryID              { return EntryID(uuid.New()) }
func NewTaskID() TaskID                { return TaskID(uuid.New()) }
func NewStageHistoryID() StageHistoryID { return StageHistoryID(uuid.New()) }
func NewCompanyAliasID() CompanyAliasID { return CompanyAliasID(uuid.New()) }

// --- String ---

func (id UserID) String() string         { return uuid.UUID(id).String() }
func (id CompanyID) String() string      { return uuid.UUID(id).String() }
func (id EntryID) String() string        { return uuid.UUID(id).String() }
func (id TaskID) String() string         { return uuid.UUID(id).String() }
func (id StageHistoryID) String() string { return uuid.UUID(id).String() }
func (id CompanyAliasID) String() string { return uuid.UUID(id).String() }

// --- IsZero: IDが未設定（ゼロ値）かどうかを判定する ---

func (id UserID) IsZero() bool         { return uuid.UUID(id) == uuid.Nil }
func (id CompanyID) IsZero() bool      { return uuid.UUID(id) == uuid.Nil }
func (id EntryID) IsZero() bool        { return uuid.UUID(id) == uuid.Nil }
func (id TaskID) IsZero() bool         { return uuid.UUID(id) == uuid.Nil }
func (id StageHistoryID) IsZero() bool { return uuid.UUID(id) == uuid.Nil }
func (id CompanyAliasID) IsZero() bool { return uuid.UUID(id) == uuid.Nil }
