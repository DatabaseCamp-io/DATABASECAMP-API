package models

type OverviewInfo struct {
	Overview            []OverviewDB
	LearningProgression []LearningProgressionDB
	ExamResult          []ExamResultDB
	ContentExam         []ContentExamDB
}

type ActivityHint struct {
	TotalHint     int           `json:"total_hint"`
	UsedHints     []HintDB      `json:"used_hints"`
	NextHintPoint *int          `json:"next_hint_point"`
	HintRoadMap   []HintRoadMap `json:"hint_roadmap"`
}

type RoadmapItem struct {
	ActivityID int  `json:"activity_id"`
	IsLearned  bool `json:"is_learned"`
	Order      int  `json:"order"`
}

type HintRoadMap struct {
	Level       int `json:"level"`
	ReducePoint int `json:"reduce_point"`
}
