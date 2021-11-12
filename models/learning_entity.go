package models

type OverviewInfo struct {
	Overview            []OverviewDB
	LearningProgression []LearningProgressionDB
	ExamResult          []ExamResultDB
	ContentExam         []ContentExamDB
}

type ContentGroupOverview struct {
	GroupID     int               `json:"group_id"`
	IsRecommend bool              `json:"is_recommend"`
	IsLasted    bool              `json:"is_lasted"`
	GroupName   string            `json:"group_name"`
	Progress    int               `json:"progress"`
	Contents    []ContentOverview `json:"contents"`
}

type ActivityHint struct {
	TotalHint     int      `json:"total_hint"`
	UsedHints     []HintDB `json:"used_hints"`
	NextHintPoint *int     `json:"next_hint_point"`
}

type RoadmapItem struct {
	ActivityID int  `json:"activity_id"`
	IsLearned  bool `json:"is_learned"`
	Order      int  `json:"order"`
}

type LastedGroup struct {
	GroupID     int    `json:"group_id"`
	ContentID   int    `json:"content_id"`
	GroupName   string `json:"group_name"`
	ContentName string `json:"content_name"`
	Progress    int    `json:"progress"`
}

type ContentOverview struct {
	ContentID   int    `json:"content_id"`
	ContentName string `json:"content_name"`
	IsLasted    bool   `json:"is_lasted"`
	Progress    int    `json:"progress"`
}
