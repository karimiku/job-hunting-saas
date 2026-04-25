package entity

import "github.com/google/uuid"

// Named types for entity IDs.
// type alias（=）ではなく named type を使うことで、
// 異なるエンティティの ID を取り違えた場合にコンパイルエラーになる。
// 例: func GetEntry(userID UserID, entryID EntryID) に対して
//
//	GetEntry(entryID, userID) はコンパイルエラーになる。

// UserID はユーザを一意に識別するための ID 型。
type UserID uuid.UUID

// CompanyID は企業を一意に識別するための ID 型。
type CompanyID uuid.UUID

// EntryID は応募 (Entry) を一意に識別するための ID 型。
type EntryID uuid.UUID

// TaskID はタスクを一意に識別するための ID 型。
type TaskID uuid.UUID

// StageHistoryID は選考フェーズ変更履歴を一意に識別するための ID 型。
type StageHistoryID uuid.UUID

// CompanyAliasID は企業の別名を一意に識別するための ID 型。
type CompanyAliasID uuid.UUID

// PasswordCredentialID はパスワード認証情報を一意に識別するための ID 型。
type PasswordCredentialID uuid.UUID

// ExternalIdentityID は外部認証連携情報を一意に識別するための ID 型。
type ExternalIdentityID uuid.UUID

// --- Constructors ---

// NewUserID は新しい UserID を生成する。
func NewUserID() UserID { return UserID(uuid.New()) }

// NewCompanyID は新しい CompanyID を生成する。
func NewCompanyID() CompanyID { return CompanyID(uuid.New()) }

// NewEntryID は新しい EntryID を生成する。
func NewEntryID() EntryID { return EntryID(uuid.New()) }

// NewTaskID は新しい TaskID を生成する。
func NewTaskID() TaskID { return TaskID(uuid.New()) }

// NewStageHistoryID は新しい StageHistoryID を生成する。
func NewStageHistoryID() StageHistoryID { return StageHistoryID(uuid.New()) }

// NewCompanyAliasID は新しい CompanyAliasID を生成する。
func NewCompanyAliasID() CompanyAliasID { return CompanyAliasID(uuid.New()) }

// NewPasswordCredentialID は新しい PasswordCredentialID を生成する。
func NewPasswordCredentialID() PasswordCredentialID { return PasswordCredentialID(uuid.New()) }

// NewExternalIdentityID は新しい ExternalIdentityID を生成する。
func NewExternalIdentityID() ExternalIdentityID { return ExternalIdentityID(uuid.New()) }

// --- String ---

// String は UserID を文字列表現で返す。
func (id UserID) String() string { return uuid.UUID(id).String() }

// String は CompanyID を文字列表現で返す。
func (id CompanyID) String() string { return uuid.UUID(id).String() }

// String は EntryID を文字列表現で返す。
func (id EntryID) String() string { return uuid.UUID(id).String() }

// String は TaskID を文字列表現で返す。
func (id TaskID) String() string { return uuid.UUID(id).String() }

// String は StageHistoryID を文字列表現で返す。
func (id StageHistoryID) String() string { return uuid.UUID(id).String() }

// String は CompanyAliasID を文字列表現で返す。
func (id CompanyAliasID) String() string { return uuid.UUID(id).String() }

// String は PasswordCredentialID を文字列表現で返す。
func (id PasswordCredentialID) String() string { return uuid.UUID(id).String() }

// String は ExternalIdentityID を文字列表現で返す。
func (id ExternalIdentityID) String() string { return uuid.UUID(id).String() }

// --- IsZero: IDが未設定（ゼロ値）かどうかを判定する ---

// IsZero は UserID がゼロ値 (未設定) かを返す。
func (id UserID) IsZero() bool { return uuid.UUID(id) == uuid.Nil }

// IsZero は CompanyID がゼロ値 (未設定) かを返す。
func (id CompanyID) IsZero() bool { return uuid.UUID(id) == uuid.Nil }

// IsZero は EntryID がゼロ値 (未設定) かを返す。
func (id EntryID) IsZero() bool { return uuid.UUID(id) == uuid.Nil }

// IsZero は TaskID がゼロ値 (未設定) かを返す。
func (id TaskID) IsZero() bool { return uuid.UUID(id) == uuid.Nil }

// IsZero は StageHistoryID がゼロ値 (未設定) かを返す。
func (id StageHistoryID) IsZero() bool { return uuid.UUID(id) == uuid.Nil }

// IsZero は CompanyAliasID がゼロ値 (未設定) かを返す。
func (id CompanyAliasID) IsZero() bool { return uuid.UUID(id) == uuid.Nil }

// IsZero は PasswordCredentialID がゼロ値 (未設定) かを返す。
func (id PasswordCredentialID) IsZero() bool { return uuid.UUID(id) == uuid.Nil }

// IsZero は ExternalIdentityID がゼロ値 (未設定) かを返す。
func (id ExternalIdentityID) IsZero() bool { return uuid.UUID(id) == uuid.Nil }
