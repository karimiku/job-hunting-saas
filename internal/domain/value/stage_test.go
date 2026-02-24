package value

import (
	"testing"
)

func TestNewStage(t *testing.T) {
	tests := []struct {
		name    string
		kind    string
		label   string
		wantErr bool
	}{
		// 正常系
		{"application", "application", "応募", false},
		{"document", "document", "ES提出", false},
		{"test", "test", "Webテスト", false},
		{"coding test", "test", "コーディングテスト", false},
		{"interview first", "interview", "一次面接", false},
		{"interview final", "interview", "最終面接", false},
		{"interview casual", "interview", "カジュアル面談", false},
		{"group discussion", "group", "GD", false},
		{"offer", "offer", "内定", false},
		{"other", "other", "座談会", false},

		// 異常系
		{"empty kind", "", "応募", true},
		{"empty label", "interview", "", true},
		{"both empty", "", "", true},
		{"invalid kind", "unknown", "何か", true},
		{"kind leading space", " interview", "面接", true},
		{"label leading space", "interview", " 一次面接", true},
		{"label trailing space", "interview", "一次面接 ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewStage(tt.kind, tt.label)
			if tt.wantErr && err == nil {
				t.Errorf("NewStage(%q, %q) should return error, but got nil", tt.kind, tt.label)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("NewStage(%q, %q) should succeed, but got error: %v", tt.kind, tt.label, err)
			}
		})
	}
}

func TestStage_Kind(t *testing.T) {
	s, err := NewStage("interview", "一次面接")
	if err != nil {
		t.Fatalf("NewStage failed: %v", err)
	}
	if s.Kind() != "interview" {
		t.Errorf("Kind() = %q, want %q", s.Kind(), "interview")
	}
}

func TestStage_Label(t *testing.T) {
	s, err := NewStage("interview", "一次面接")
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
		kindA string
		lblA  string
		kindB string
		lblB  string
		want  bool
	}{
		{"same kind and label", "interview", "一次面接", "interview", "一次面接", true},
		{"same kind different label", "interview", "一次面接", "interview", "最終面接", false},
		{"different kind same label", "test", "テスト", "interview", "テスト", false},
		{"both different", "document", "ES", "interview", "一次面接", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stageA, err := NewStage(tt.kindA, tt.lblA)
			if err != nil {
				t.Fatalf("NewStage(%q, %q) failed: %v", tt.kindA, tt.lblA, err)
			}
			stageB, err := NewStage(tt.kindB, tt.lblB)
			if err != nil {
				t.Fatalf("NewStage(%q, %q) failed: %v", tt.kindB, tt.lblB, err)
			}
			if got := stageA.Equals(stageB); got != tt.want {
				t.Errorf("Stage(%q/%q).Equals(Stage(%q/%q)) = %v, want %v",
					tt.kindA, tt.lblA, tt.kindB, tt.lblB, got, tt.want)
			}
		})
	}
}
