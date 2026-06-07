// Package jobemail は選考メール本文から就活管理に必要な候補情報を抽出する。
package jobemail

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ExtractInput は選考メール抽出ユースケースへの入力。
type ExtractInput struct {
	Subject     string
	Text        string
	CompanyName string
	Now         time.Time
}

// ExtractOutput は選考メールから抽出した候補情報。
type ExtractOutput struct {
	CompanyName          string                 `json:"companyName"`
	StageKind            string                 `json:"stageKind"`
	StageLabel           string                 `json:"stageLabel"`
	EventAt              *string                `json:"eventAt"`
	DeadlineAt           *string                `json:"deadlineAt"`
	EntryUpdateCandidate map[string]string      `json:"entryUpdateCandidate"`
	SuggestedTasks       []map[string]string    `json:"suggestedTasks"`
	Notes                []string               `json:"notes"`
	RawSignals           map[string]interface{} `json:"rawSignals"`
}

// Extract は選考メール本文から就活管理用の候補情報を抽出するUseCase。
type Extract struct{}

// NewExtract は選考メール抽出ユースケースを生成する。
func NewExtract() *Extract {
	return &Extract{}
}

// Execute は選考メール本文を解析し、エントリー更新候補とタスク候補を返す。
func (uc *Extract) Execute(input ExtractInput) ExtractOutput {
	now := input.Now
	if now.IsZero() {
		now = time.Now()
	}
	text := strings.TrimSpace(input.Subject + "\n" + input.Text)
	company := strings.TrimSpace(input.CompanyName)
	if company == "" {
		company = guessCompanyName(text)
	}
	stageKind, stageLabel := guessStage(text)
	dates := extractDateTimes(text, now)

	var eventAt *string
	if len(dates) > 0 {
		eventAt = &dates[0]
	}
	var deadlineAt *string
	for _, sentence := range splitSentences(text) {
		if !strings.Contains(sentence, "締切") && !strings.Contains(sentence, "期限") && !strings.Contains(sentence, "まで") {
			continue
		}
		ds := extractDateTimes(sentence, now)
		if len(ds) > 0 {
			deadlineAt = &ds[0]
			break
		}
	}

	tasks := make([]map[string]string, 0, 2)
	if deadlineAt != nil {
		tasks = append(tasks, map[string]string{
			"title":   stageLabel + "の対応",
			"type":    "deadline",
			"dueDate": *deadlineAt,
		})
	}
	if eventAt != nil && (deadlineAt == nil || *eventAt != *deadlineAt) {
		tasks = append(tasks, map[string]string{
			"title":   stageLabel + "の準備",
			"type":    "schedule",
			"dueDate": *eventAt,
		})
	}

	return ExtractOutput{
		CompanyName: company,
		StageKind:   stageKind,
		StageLabel:  stageLabel,
		EventAt:     eventAt,
		DeadlineAt:  deadlineAt,
		EntryUpdateCandidate: map[string]string{
			"companyName": company,
			"stageKind":   stageKind,
			"stageLabel":  stageLabel,
		},
		SuggestedTasks: tasks,
		Notes: []string{
			"この抽出はMCPサーバー内のルールベース処理です。LLM APIは呼んでいません。",
			"Entry更新やTask作成の前に、会社名・日時・締切をユーザーが確認してください。",
		},
		RawSignals: map[string]interface{}{
			"dateTimeCandidates": dates,
			"subject":            input.Subject,
		},
	}
}

func guessCompanyName(text string) string {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(株式会社[^\s　、。\n]{1,30})`),
		regexp.MustCompile(`([^\s　、。\n]{1,30}株式会社)`),
		regexp.MustCompile(`(合同会社[^\s　、。\n]{1,30})`),
	}
	for _, pattern := range patterns {
		if m := pattern.FindStringSubmatch(text); len(m) > 1 {
			return strings.Trim(m[1], "「」[]【】")
		}
	}
	return ""
}

func guessStage(text string) (string, string) {
	checks := []struct {
		words []string
		kind  string
		label string
	}{
		{[]string{"内定", "オファー"}, "offer", "内定・オファー"},
		{[]string{"最終面接"}, "interview", "最終面接"},
		{[]string{"一次面接", "1次面接"}, "interview", "一次面接"},
		{[]string{"二次面接", "2次面接"}, "interview", "二次面接"},
		{[]string{"面接", "面談"}, "interview", "面接"},
		{[]string{"グループディスカッション", "GD"}, "group", "グループ選考"},
		{[]string{"SPI", "適性検査", "Webテスト", "WEBテスト"}, "test", "適性検査"},
		{[]string{"ES", "エントリーシート", "書類"}, "document", "書類・ES"},
	}
	for _, check := range checks {
		for _, word := range check.words {
			if strings.Contains(text, word) {
				return check.kind, check.label
			}
		}
	}
	return "application", "応募"
}

func extractDateTimes(text string, now time.Time) []string {
	var out []string
	seen := map[string]struct{}{}
	add := func(t time.Time) {
		value := t.Format(time.RFC3339)
		if _, ok := seen[value]; ok {
			return
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}

	ymd := regexp.MustCompile(`(20\d{2})[-/](\d{1,2})[-/](\d{1,2})(?:[ T　]*(\d{1,2})[:：時](\d{2})?)?`)
	for _, m := range ymd.FindAllStringSubmatch(text, -1) {
		year, _ := strconv.Atoi(m[1])
		month, _ := strconv.Atoi(m[2])
		day, _ := strconv.Atoi(m[3])
		hour, minute := parseHourMinute(m[4], m[5])
		add(time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.Local))
	}

	jp := regexp.MustCompile(`(\d{1,2})月(\d{1,2})日(?:\([^)]*\))?(?:[ 　]*(\d{1,2})[:：時](\d{2})?)?`)
	for _, m := range jp.FindAllStringSubmatch(text, -1) {
		month, _ := strconv.Atoi(m[1])
		day, _ := strconv.Atoi(m[2])
		hour, minute := parseHourMinute(m[3], m[4])
		year := now.Year()
		t := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.Local)
		if t.Before(now.AddDate(0, -1, 0)) {
			t = t.AddDate(1, 0, 0)
		}
		add(t)
	}
	return out
}

func parseHourMinute(hourRaw, minuteRaw string) (int, int) {
	if hourRaw == "" {
		return 0, 0
	}
	hour, _ := strconv.Atoi(hourRaw)
	minute := 0
	if minuteRaw != "" {
		minute, _ = strconv.Atoi(minuteRaw)
	}
	return hour, minute
}

func splitSentences(text string) []string {
	return regexp.MustCompile(`[。\n]`).Split(text, -1)
}
