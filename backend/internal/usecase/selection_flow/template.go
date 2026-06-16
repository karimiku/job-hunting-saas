package selectionflow

// DefaultStageInputs はGUI作成で使う標準選考フロー。
func DefaultStageInputs() []StageInput {
	return []StageInput{
		{StageKind: "application", StageLabel: "エントリー"},
		{StageKind: "document", StageLabel: "書類"},
		{StageKind: "test", StageLabel: "テスト"},
		{StageKind: "interview", StageLabel: "面接"},
		{StageKind: "group", StageLabel: "GD"},
		{StageKind: "offer", StageLabel: "内定"},
	}
}
