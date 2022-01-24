package user

import (
	"database-camp/internal/models/entities/content"
)

type SpiderData struct {
	ContentGroupID   int     `gorm:"column:content_group_id" json:"content_group_id"`
	ContentGroupName string  `gorm:"-" json:"content_group_name"`
	Stat             float64 `gorm:"column:stat" json:"stat"`
}

type SpiderDataset []SpiderData

func (dataset *SpiderDataset) FillContentGroups(contentGroups content.ContentGroups) {
	newDataset := make(SpiderDataset, 0)
	oldDataset := *dataset

	mapData := map[int]SpiderData{}
	for _, data := range oldDataset {
		mapData[data.ContentGroupID] = data
	}

	for _, group := range contentGroups {
		newDataset = append(newDataset, SpiderData{
			ContentGroupID:   group.ID,
			ContentGroupName: group.Name,
			Stat:             mapData[group.ID].Stat,
		})
	}

	*dataset = newDataset
}
