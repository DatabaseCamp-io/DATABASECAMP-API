package models

type Content struct {
	ID        int    `gorm:"primaryKey;column:content_id" json:"content_id"`
	GroupID   int    `gorm:"column:content_group_id" json:"content_group_id"`
	Name      string `gorm:"column:name" json:"name"`
	VideoPath string `gorm:"column:view_path" json:"view_path"`
	SlidePath string `gorm:"column:slide_path" json:"slide"`
}

type VideoLectureResponse struct {
	ContentID   int    `json:"content_id"`
	ContentName string `json:"content_name"`
	VideoLink   string `json:"video_link"`
}
