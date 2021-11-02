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
	exam := make([]models.ExamResultDB, 0)
	contentExam := make([]models.ContentExamDB, 0)

	concurrent := models.Concurrent{
		Wg:  &wg,
		Err: &err,
	}

	wg.Add(4)
	go c.loadOverviewAsync(&concurrent, &overview)
	go c.loadLearningProgressionAsync(&concurrent, &learningProgression, id)
	go c.loadFailedExamAsync(&concurrent, &exam, id)
	go c.loadContentExamPretestAsync(&concurrent, &contentExam)
	wg.Wait()

	if err != nil {
		return nil, err
	}

	info := models.OverviewInfo{
		Overview:            overview,
		LearningProgression: learningProgression,
		ExamResult:          exam,
		ContentExam:         contentExam,
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

func (c learningController) loadFailedExamAsync(concurrent *models.Concurrent, exam *[]models.ExamResultDB, id int) {
	defer concurrent.Wg.Done()
	result, err := c.userRepo.GetFailedExam(id)
	if err != nil {
		*concurrent.Err = err
	}
	*exam = append(*exam, result...)
}

func (c learningController) loadContentExamPretestAsync(concurrent *models.Concurrent, contentExamDB *[]models.ContentExamDB) {
	defer concurrent.Wg.Done()
	result, err := c.learningRepo.GetContentExam(models.Exam.Pretest)
	if err != nil {
		*concurrent.Err = err
	}
	*contentExamDB = append(*contentExamDB, result...)
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
				g.content[v.ContentID] = c
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

func (c learningController) getRecommendGroupFromExam(info *models.OverviewInfo) map[int]bool {
	recommendedGroup := map[int]bool{}
	examGroupMap := map[int]int{}
	for _, v := range info.ContentExam {
		examGroupMap[v.ActivityID] = v.GroupID
	}
	for _, v := range info.ExamResult {
		recommendedGroup[examGroupMap[v.ActivityID]] = true
	}
	return recommendedGroup
}

func (c learningController) calculateProgress(progress int, total int) int {
	if total == 0 {
		return 100
	} else {
		return (progress / total) * 100
	}

}

func (c learningController) prepareOverviewResponse(info *models.OverviewInfo, data overviewData) models.OverviewResponse {
	var lastedActivityID int
	var lastedGroup *models.LastedGroup
	userActivityCount := map[int]int{}
	contentGroupOverview := make([]models.ContentGroupOverview, 0)

	countRecommend := 0
	recommendGroup := c.getRecommendGroupFromExam(info)

	for i, v := range info.LearningProgression {
		if i == 0 {
			lastedActivityID = v.ActivityID
		}
		userActivityCount[v.ActivityID]++
	}

	for ko, vo := range data.group {
		content := make([]models.ContentOverview, 0)
		_isLasted := false
		_isGroupLasted := false
		countUserActivity := 0
		countActivity := 0
		for kc, vc := range vo.content {

			countActivity += data.activityCount[kc]
			countUserActivity += userActivityCount[kc]
			_isLasted = lastedActivityID == kc
			progress := c.calculateProgress(userActivityCount[kc], data.activityCount[kc])
			content = append(content, models.ContentOverview{
				ContentID:   kc,
				ContentName: vc.name,
				IsLasted:    _isLasted,
				Progress:    progress,
			})
			if _isLasted {
				_isGroupLasted = true
				lastedGroup = &models.LastedGroup{
					GroupID:     ko,
					ContentID:   kc,
					GroupName:   vo.name,
					ContentName: vc.name,
					Progress:    progress,
				}
			}
		}

		isRecommend := recommendGroup[ko]
		if isRecommend {
			countRecommend++
			if countRecommend > 3 {
				isRecommend = false
			}
		}

		contentGroupOverview = append(contentGroupOverview, models.ContentGroupOverview{
			GroupID:     ko,
			IsRecommend: isRecommend,
			IsLasted:    _isGroupLasted,
			GroupName:   vo.name,
			Progress:    c.calculateProgress(countUserActivity, countActivity),
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
