package activity

type PeerProblem struct {
	ID      int          `json:"peer_problem_id"`
	Problem string       `json:"problem"`
	Choices []PeerChoice `json:"choices"`
}
type PeerChoice struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

var PeerProblems = []PeerProblem{
	{
		ID:      1,
		Problem: "ในด้าน Entity มีความถูกต้องหรือไม่",
		Choices: []PeerChoice{
			{ID: 1, Value: "ไม่มีความถูกต้อง เนื่องจากยังมี Entity ไม่ถูกต้องตามความต้องการของระบบ"},
			{ID: 2, Value: "ไม่มีความถูกต้อง เนื่องจาก Entity ยังสามารถปรับให้ใช้งานง่ายขึ้นได้"},
			{ID: 3, Value: "มีความถูกต้อง เนื่องจากมี Entity ครบถ้วน"},
		},
	},
	{
		ID:      2,
		Problem: "ในด้าน Attribute มีความถูกต้องหรือไม่",
		Choices: []PeerChoice{
			{ID: 4, Value: "ไม่มีความถูกต้อง เนื่องจากยังมี Attribute ไม่ครบตามความต้องการของระบบ"},
			{ID: 5, Value: "ไม่มีความถูกต้อง เนื่องจาก Attribute ยังสามารถเพิ่มเติมให้ใช้งานง่ายขึ้นได้"},
			{ID: 6, Value: "มีความถูกต้อง เนื่องจากมี Attribute ครบถ้วน"},
		},
	},
	{
		ID:      3,
		Problem: "ในด้าน Relationship มีความถูกต้องหรือไม่",
		Choices: []PeerChoice{
			{ID: 7, Value: "ไม่มีความถูกต้อง เนื่องจากยังมี Relationship ไม่ครบตามความต้องการของระบบ"},
			{ID: 8, Value: "ไม่มีความถูกต้อง เนื่องจาก Relationship ยังสามารถเพิ่มเติมให้ดีขึ้นได้"},
			{ID: 9, Value: "มีความถูกต้อง เนื่องจากมี Relationship ครบถ้วนและถูกต้อง"},
		},
	},
	{
		ID:      4,
		Problem: "ในด้านความซ้ำซ้อนของข้อมูล (Redundancy) มีความถูกต้องหรือไม่",
		Choices: []PeerChoice{
			{ID: 10, Value: "มีความถูกต้อง เนื่องจากข้อมูลไม่สามารถเกิดปัญหาจากความซ้ำซ้อนได้"},
			{ID: 11, Value: "ไม่ถูกต้อง เนื่องจากข้อมูลเกิดความซ้ำซ้อนได้"},
		},
	},
}
