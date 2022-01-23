package content

type Content struct {
	ID        int    `gorm:"primaryKey;column:content_id" json:"content_id"`
	GroupID   int    `gorm:"column:content_group_id" json:"content_group_id"`
	Name      string `gorm:"column:name" json:"name"`
	VideoPath string `gorm:"column:video_path" json:"video_path"`
	SlidePath string `gorm:"column:slide_path" json:"slide"`
}

type ContentGroup struct {
	ContentGroupID int `gorm:"column:content_group_id"`
}

type ContentGroups []ContentGroup
