package controller

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/models"
	"DatabaseCamp/repository"
	"DatabaseCamp/services"
)

type learningController struct {
	repo    repository.ILearningRepository
	service services.IAwsService
}

type ILearningController interface {
	GetVideoLecture(id int) (*models.VideoLectureResponse, error)
}

func NewLearningController(repo repository.ILearningRepository, service services.IAwsService) learningController {
	return learningController{repo: repo, service: service}
}

func (c learningController) GetVideoLecture(id int) (*models.VideoLectureResponse, error) {
	content, err := c.repo.GetContent(id)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewNotFoundError("ไม่พบเนื้อหา", "Content Not Found")
	}

	videoLink, err := c.service.GetFileLink(content.VideoPath)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewServiceUnavailableError("Service ไม่พร้อมใช้งาน", "Service Unavailable")
	}

	res := models.VideoLectureResponse{
		ContentID:   content.ID,
		ContentName: content.Name,
		VideoLink:   videoLink,
	}

	return &res, nil
}
