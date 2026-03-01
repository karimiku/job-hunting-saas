package value

import (
	"errors"
	"testing"
)

func TestNewStageKind(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		wantErr error
	}{
		// 正常系
		{"application", "application", nil},
		{"document", "document", nil},
		{"test", "test", nil},
		{"interview", "interview", nil},
		{"group", "group", nil},
		{"offer", "offer", nil},
		{"other", "other", nil},

		// 異常系
		{"empty", "", ErrStageKindEmpty},
		{"invalid", "unknown", ErrStageKindInvalid},
		{"leading space", " interview", ErrStageKindInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewStageKind(tt.raw)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewStageKind(%q) error = %v, wantErr %v", tt.raw, err, tt.wantErr)
			}
		})
	}
}

func TestStageKind_String(t *testing.T) {
	k, err := NewStageKind("interview")
	if err != nil {
		t.Fatalf("NewStageKind failed: %v", err)
	}
	if k.String() != "interview" {
		t.Errorf("String() = %q, want %q", k.String(), "interview")
	}
}

func TestStageKind_Equals(t *testing.T) {
	a := StageKindInterview()
	b := StageKindInterview()
	c := StageKindOffer()

	if !a.Equals(b) {
		t.Error("same kind should be equal")
	}
	if a.Equals(c) {
		t.Error("different kind should not be equal")
	}
}

func TestStageKind_ConstantConstructors(t *testing.T) {
	tests := []struct {
		name string
		fn   func() StageKind
		want string
	}{
		{"application", StageKindApplication, "application"},
		{"document", StageKindDocument, "document"},
		{"test", StageKindTest, "test"},
		{"interview", StageKindInterview, "interview"},
		{"group", StageKindGroup, "group"},
		{"offer", StageKindOffer, "offer"},
		{"other", StageKindOther, "other"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fn()
			if k.String() != tt.want {
				t.Errorf("String() = %q, want %q", k.String(), tt.want)
			}
		})
	}
}

func TestNewStage(t *testing.T) {
	tests := []struct {
		name    string
		kind    StageKind
		label   string
		wantErr error
	}{
		// 正常系
		{"application", StageKindApplication(), "応募", nil},
		{"document", StageKindDocument(), "ES提出", nil},
		{"test", StageKindTest(), "Webテスト", nil},
		{"coding test", StageKindTest(), "コーディングテスト", nil},
		{"interview first", StageKindInterview(), "一次面接", nil},
		{"interview final", StageKindInterview(), "最終面接", nil},
		{"interview casual", StageKindInterview(), "カジュアル面談", nil},
		{"group discussion", StageKindGroup(), "GD", nil},
		{"offer", StageKindOffer(), "内定", nil},
		{"other", StageKindOther(), "座談会", nil},

		// 異常系
		{"empty label", StageKindInterview(), "", ErrStageLabelEmpty},
		{"label leading space", StageKindInterview(), " 一次面接", ErrStageLabelInvalid},
		{"label trailing space", StageKindInterview(), "一次面接 ", ErrStageLabelInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewStage(tt.kind, tt.label)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewStage(%q, %q) error = %v, wantErr %v", tt.kind.String(), tt.label, err, tt.wantErr)
			}
		})
	}
}

func TestStage_Kind(t *testing.T) {
	s, err := NewStage(StageKindInterview(), "一次面接")
	if err != nil {
		t.Fatalf("NewStage failed: %v", err)
	}
	if s.Kind().String() != "interview" {
		t.Errorf("Kind() = %q, want %q", s.Kind().String(), "interview")
	}
}

func TestStage_Label(t *testing.T) {
	s, err := NewStage(StageKindInterview(), "一次面接")
	if err != nil {
		t.Fatalf("NewStage failed: %v", err)
	}
	if s.Label() != "一次面接" {
		t.Errorf("Label() = %q, want %q", s.Label(), "一次面接")
	}
}

func TestStage_Equals(t *testing.T) {
	tests := []struct {
		name  string
		kindA StageKind
		lblA  string
		kindB StageKind
		lblB  string
		want  bool
	}{
		{"same kind and label", StageKindInterview(), "一次面接", StageKindInterview(), "一次面接", true},
		{"same kind different label", StageKindInterview(), "一次面接", StageKindInterview(), "最終面接", false},
		{"different kind same label", StageKindTest(), "テスト", StageKindInterview(), "テスト", false},
		{"both different", StageKindDocument(), "ES", StageKindInterview(), "一次面接", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stageA, err := NewStage(tt.kindA, tt.lblA)
			if err != nil {
				t.Fatalf("NewStage(%q, %q) failed: %v", tt.kindA.String(), tt.lblA, err)
			}
			stageB, err := NewStage(tt.kindB, tt.lblB)
			if err != nil {
				t.Fatalf("NewStage(%q, %q) failed: %v", tt.kindB.String(), tt.lblB, err)
			}
			if got := stageA.Equals(stageB); got != tt.want {
				t.Errorf("Stage(%q/%q).Equals(Stage(%q/%q)) = %v, want %v",
					tt.kindA.String(), tt.lblA, tt.kindB.String(), tt.lblB, got, tt.want)
			}
		})
	}
}
