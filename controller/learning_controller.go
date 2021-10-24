package controller

import (
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/models"
	"DatabaseCamp/repository"
	"DatabaseCamp/services"
	"sync"
)

type content struct {
	id       int
	name     string
	activity []*int
}

type group struct {
	id      int
	name    string
	content map[int]*content
}

type overviewData struct {
	group              map[int]*group
	activityContentMap map[int]int
	activityCount      map[int]int
}

type learningController struct {
	learningRepo repository.ILearningRepository
	userRepo     repository.IUserRepository
	service      services.IAwsService
}

type ILearningController interface {
	GetVideoLecture(id int) (*models.VideoLectureResponse, error)
	GetOverview(id int) (*models.OverviewResponse, error)
}

func NewLearningController(
	learningRepo repository.ILearningRepository,
	userRepo repository.IUserRepository,
	service services.IAwsService,
) learningController {
	return learningController{
		learningRepo: learningRepo,
		userRepo:     userRepo,
		service:      service,
	}
}

func (c learningController) GetVideoLecture(id int) (*models.VideoLectureResponse, error) {
	content, err := c.learningRepo.GetContent(id)
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

func (c learningController) loadOverviewInfo(id int) (*models.OverviewInfo, error) {
	var wg sync.WaitGroup
	var err error
	overview := make([]models.OverviewDB, 0)
	learningProgression := make([]models.LearningProgressionDB, 0)
	exam := make([]models.ExamDB, 0)

	concurrent := models.Concurrent{
		Wg:  &wg,
		Err: &err,
	}

	wg.Add(3)
	go c.loadOverviewAsync(&concurrent, &overview)
	go c.loadLearningProgressionAsync(&concurrent, &learningProgression, id)
	go c.loadExamAsync(&concurrent, &exam, id)
	wg.Wait()

	if err != nil {
		return nil, err
	}

	info := models.OverviewInfo{
		Overview:            overview,
		LearningProgression: learningProgression,
		Exam:                exam,
	}

	return &info, err
}

func (c learningController) loadOverviewAsync(concurrent *models.Concurrent, overview *[]models.OverviewDB) {
	defer concurrent.Wg.Done()
	result, err := c.learningRepo.GetOverview()
	if err != nil {
		*concurrent.Err = err
	}
	*overview = append(*overview, result...)
}

func (c learningController) loadLearningProgressionAsync(concurrent *models.Concurrent, learningProgression *[]models.LearningProgressionDB, id int) {
	defer concurrent.Wg.Done()
	result, err := c.userRepo.GetLearningProgression(id)
	if err != nil {
		*concurrent.Err = err
	}
	*learningProgression = append(*learningProgression, result...)
}

func (c learningController) loadExamAsync(concurrent *models.Concurrent, exam *[]models.ExamDB, id int) {
	defer concurrent.Wg.Done()
	result, err := c.userRepo.GetExam(id)
	if err != nil {
		*concurrent.Err = err
	}
	*exam = append(*exam, result...)
}

func (c learningController) prepareOverview(info *models.OverviewInfo) overviewData {
	root := map[int]*group{}
	activityContentMap := map[int]int{}
	activityCount := map[int]int{}
	for _, v := range info.Overview {
		_v := v
		g := root[v.GroupID]
		if g == nil {
			c := &content{
				id:       v.ContentID,
				name:     v.ContentName,
				activity: []*int{_v.ActivityID},
			}

			g = &group{
				id:      v.GroupID,
				name:    v.GroupName,
				content: map[int]*content{v.ContentID: c},
			}
			root[v.GroupID] = g
		} else {
			c := root[v.GroupID].content[v.ContentID]
			if c == nil {
				var _c *content
				c = &content{
					id:       v.ContentID,
					name:     v.ContentName,
					activity: []*int{_v.ActivityID},
				}
				_c = &content{
					id:       v.ContentID,
					name:     v.ContentName,
					activity: []*int{_v.ActivityID},
				}
				c = _c
			} else {
				c.activity = append(c.activity, _v.ActivityID)
			}
		}
		if v.ActivityID != nil {
			activityContentMap[*v.ActivityID] = v.ContentID
			activityCount[*v.ActivityID] += 1
		}
	}

	overviewData := overviewData{
		group:              root,
		activityContentMap: activityContentMap,
		activityCount:      activityCount,
	}

	return overviewData
}

func (c learningController) prepareOverviewResponse(info *models.OverviewInfo, data overviewData) models.OverviewResponse {
	var lastedActivityID int
	var lastedGroup models.LastedGroup
	userActivityCount := map[int]int{}
	contentGroupOverview := make([]models.ContentGroupOverview, 0)

	for i, v := range info.LearningProgression {
		if i == 0 {
			lastedActivityID = v.ActivityID
		}
		userActivityCount[v.ActivityID]++
	}

	for ko, vo := range data.group {
		content := make([]models.ContentOverview, 0)
		_isLasted := false
		countUserActivity := 0
		countActivity := 0
		for kc, vc := range vo.content {
			countActivity += data.activityCount[kc]
			countUserActivity += userActivityCount[kc]
			_isLasted = lastedActivityID == kc
			progress := (userActivityCount[kc] / data.activityCount[kc]) * 100
			content = append(content, models.ContentOverview{
				ContentID:   kc,
				ContentName: vc.name,
				IsLasted:    _isLasted,
				Progress:    progress,
			})
			if _isLasted {
				lastedGroup = models.LastedGroup{
					GroupID:     ko,
					ContentID:   kc,
					GroupName:   vo.name,
					ContentName: vc.name,
					Progress:    progress,
				}
			}
		}

		contentGroupOverview = append(contentGroupOverview, models.ContentGroupOverview{
			GroupID:     ko,
			IsRecommend: false,
			IsLasted:    _isLasted,
			GroupName:   vo.name,
			Progress:    (countUserActivity / countActivity) * 100,
			Contents:    content,
		})
	}
	return models.OverviewResponse{
		LastedGroup:          lastedGroup,
		ContentGroupOverview: contentGroupOverview,
	}
}

func (c learningController) GetOverview(id int) (*models.OverviewResponse, error) {
	info, err := c.loadOverviewInfo(id)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewInternalServerError("เกิดข้อผิดพลาด", "Internal Server Error")
	}
	data := c.prepareOverview(info)
	res := c.prepareOverviewResponse(info, data)
	return &res, nil
}
