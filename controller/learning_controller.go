package controller

import (
	"DatabaseCamp/database"
	"DatabaseCamp/errs"
	"DatabaseCamp/logs"
	"DatabaseCamp/models"
	"DatabaseCamp/repository"
	"DatabaseCamp/services"
	"DatabaseCamp/utils"
	"sync"
	"time"
)

type checkAnswerInfo struct {
	choice   interface{}
	activity *models.ActivityDB
}

type roadmapInfo struct {
	content         *models.ContentDB
	contentActivity []models.ActivityDB
	progression     []models.LearningProgressionDB
}

type activityInfo struct {
	activityHints []models.HintDB
	userHints     []models.UserHintDB
	activity      models.ActivityDB
}

type hintInfo struct {
	activityHints []models.HintDB
	userHints     []models.UserHintDB
	user          models.User
}

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
	GetActivity(userID int, activityID int) (*models.ActivityResponse, error)
	CheckMatchingAnswer(userID int, request models.MatchingChoiceAnswerRequest) (*models.AnswerResponse, error)
	UseHint(userID int, activityID int) (*models.HintDB, error)
	CheckCompletionAnswer(userID int, request models.CompletionAnswerRequest) (interface{}, error)
	GetContentRoadmap(userID int, contentID int) (*models.RoadmapResponse, error)
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

func (c learningController) loadContentActivityAsync(concurrent *models.Concurrent, contentID int, activity *[]models.ActivityDB) {
	defer concurrent.Wg.Done()
	result, err := c.learningRepo.GetContentActivity(contentID)
	if err != nil {
		*concurrent.Err = err
	}
	*activity = append(*activity, result...)
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

func (c learningController) getChoice(activityID int, typeID int) (interface{}, error) {
	if typeID == 1 {
		return c.learningRepo.GetMatchingChoice(activityID)
	} else if typeID == 2 {
		return c.learningRepo.GetMultipleChoice(activityID)
	} else if typeID == 3 {
		return c.learningRepo.GetCompletionChoice(activityID)
	} else {
		return nil, errs.NewInternalServerError("เกิดข้อผิดพลาด", "Internal Server Error")
	}
}

func (c learningController) prepareChoice(typeID int, choice interface{}) interface{} {
	activityManager := services.NewActivityManager()
	if typeID == 1 {
		return activityManager.PrepareMatchingChoice(choice.([]models.MatchingChoiceDB))
	} else if typeID == 2 {
		return activityManager.PrepareMultipleChoice(choice.([]models.MultipleChoiceDB))
	} else if typeID == 3 {
		return activityManager.PrepareCompletionChoice(choice.([]models.CompletionChoiceDB))
	} else {
		return nil
	}
}

func (c learningController) loadActivityInfo(userID int, activityID int) (*activityInfo, error) {
	var wg sync.WaitGroup
	var err error
	var activity *models.ActivityDB
	concurrent := models.Concurrent{Wg: &wg, Err: &err}
	userHints := make([]models.UserHintDB, 0)
	activityHints := make([]models.HintDB, 0)
	wg.Add(3)
	go c.loadActivityAsync(&concurrent, activityID, &activity)
	go c.loadActivityHints(&concurrent, activityID, &activityHints)
	go c.loadUserHintsAsync(&concurrent, userID, activityID, &userHints)
	wg.Wait()
	if err != nil {
		return nil, err
	}
	info := activityInfo{
		activityHints: activityHints,
		userHints:     userHints,
		activity:      *activity,
	}
	return &info, nil
}

func (c learningController) loadActivityAsync(concurrent *models.Concurrent, activityID int, activity **models.ActivityDB) {
	defer concurrent.Wg.Done()
	var err error
	*activity, err = c.learningRepo.GetActivity(activityID)
	if err != nil {
		*concurrent.Err = err
	}
}

func (c learningController) prepareActivityHint(info activityInfo) *models.ActivityHint {
	var nextHintPoint *int
	usedHint := make([]models.HintDB, 0)

	for _, v := range info.activityHints {
		if c.isUsedHint(info.userHints, v.ID) {
			usedHint = append(usedHint, v)
		} else if nextHintPoint == nil {
			nextHintPoint = &v.PointReduce
		}
	}

	return &models.ActivityHint{
		TotalHint:     len(info.activityHints),
		UsedHints:     usedHint,
		NextHintPoint: nextHintPoint,
	}
}

func (c learningController) GetActivity(userID int, activityID int) (*models.ActivityResponse, error) {
	info, err := c.loadActivityInfo(userID, activityID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewNotFoundError("ไม่พบกิจกรรม", "Activity Not Found")
	}

	choice, err := c.getChoice(info.activity.ID, info.activity.TypeID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewNotFoundError("ไม่พบกิจกรรม", "Activity Not Found")
	}

	preparedChoice := c.prepareChoice(info.activity.TypeID, choice)
	activityHint := c.prepareActivityHint(*info)

	res := models.ActivityResponse{
		Activity: info.activity,
		Choice:   preparedChoice,
		Hint:     *activityHint,
	}

	return &res, nil
}

func (c learningController) finishActivityTrasaction(userID int, activityID int, addPoint int) error {
	tx := database.NewTransaction()

	progression := models.LearningProgressionDB{
		UserID:           userID,
		ActivityID:       activityID,
		CreatedTimestamp: time.Now().Local(),
	}
	_, err := c.userRepo.InsertLearningProgressionTransaction(tx, progression)
	if err != nil && !utils.NewHelper().IsSqlDuplicateError(err) {
		tx.Rollback()
		return err
	}

	if !utils.NewHelper().IsSqlDuplicateError(err) {
		err = c.userRepo.ChangePointTransaction(tx, userID, addPoint, models.Mode.Add)
	}

	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	tx.Close()

	return nil
}

func (c learningController) loadCheckAnswerInfo(activityID int, activityType int) (*checkAnswerInfo, error) {
	var wg sync.WaitGroup
	var err error
	var activity *models.ActivityDB
	var choice interface{}
	concurrent := models.Concurrent{
		Wg:  &wg,
		Err: &err,
	}
	wg.Add(2)
	go c.loadActivityAsync(&concurrent, activityID, &activity)
	go c.getChioceAsync(&concurrent, activityID, activityType, &choice)
	wg.Wait()
	if err != nil {
		return nil, err
	}

	info := checkAnswerInfo{
		choice:   choice,
		activity: activity,
	}

	return &info, nil
}

func (c learningController) getChioceAsync(concurrent *models.Concurrent, activityID int, activityType int, choice *interface{}) {
	defer concurrent.Wg.Done()
	var err error
	*choice, err = c.getChoice(activityID, activityType)
	if err != nil {
		*concurrent.Err = err
	}
}

func (c learningController) CheckMatchingAnswer(userID int, request models.MatchingChoiceAnswerRequest) (*models.AnswerResponse, error) {
	info, err := c.loadCheckAnswerInfo(*request.ActivityID, 1)
	if err != nil || info.activity == nil {
		logs.New().Error(err)
		return nil, errs.NewNotFoundError("ไม่พบกิจกรรม", "Activity Not Found")
	}

	matchingChoice := info.choice.([]models.MatchingChoiceDB)
	isCorrect := true

	if len(matchingChoice) != len(request.Answer) {
		return nil, errs.NewBadRequestError("รูปแบบของคำตอบไม่ถูกต้อง", "Invalid Answer Format")
	}

	for _, correct := range matchingChoice {
		for _, answer := range request.Answer {
			if (correct.PairItem1 == *answer.Item1) && (correct.PairItem2 != *answer.Item2) {
				isCorrect = false
				break
			}
		}
	}

	response := models.AnswerResponse{
		ActivityID: *request.ActivityID,
		IsCorrect:  isCorrect,
	}

	if isCorrect {
		err = c.finishActivityTrasaction(userID, *request.ActivityID, info.activity.Point)
		if err != nil && !utils.NewHelper().IsSqlDuplicateError(err) {
			logs.New().Error(err)
			return nil, errs.NewInternalServerError("เกิดข้อผิดพลาดในการบันทึกกิจกรรม", "Saved Activity Failed")
		}
	}

	return &response, nil
}

func (c learningController) loadHintInfo(userID int, activityID int) (hintInfo, error) {
	var wg sync.WaitGroup
	var err error
	var user models.User
	concurrent := models.Concurrent{Wg: &wg, Err: &err}
	userHints := make([]models.UserHintDB, 0)
	activityHints := make([]models.HintDB, 0)
	wg.Add(3)
	go c.loadActivityHints(&concurrent, activityID, &activityHints)
	go c.loadUserHintsAsync(&concurrent, userID, activityID, &userHints)
	go c.loadUser(&concurrent, userID, &user)
	wg.Wait()
	hintInfo := hintInfo{
		activityHints: activityHints,
		userHints:     userHints,
		user:          user,
	}
	return hintInfo, err
}

func (c learningController) loadUser(concurrent *models.Concurrent, userID int, user *models.User) {
	defer concurrent.Wg.Done()
	var err error
	*user, err = c.userRepo.GetUserByID(userID)
	if err != nil {
		*concurrent.Err = err
	}
}

func (c learningController) loadUserHintsAsync(concurrent *models.Concurrent, userID int, activityID int, hints *[]models.UserHintDB) {
	defer concurrent.Wg.Done()
	userHints, e := c.userRepo.GetUserHint(userID, activityID)
	if e != nil {
		*concurrent.Err = e
	}
	*hints = append(*hints, userHints...)
}

func (c learningController) loadActivityHints(concurrent *models.Concurrent, activityID int, hints *[]models.HintDB) {
	defer concurrent.Wg.Done()
	activityHints, e := c.learningRepo.GetActivityHints(activityID)
	if e != nil {
		*concurrent.Err = e
	}
	*hints = append(*hints, activityHints...)
}

func (c learningController) isUsedHint(userHints []models.UserHintDB, hintID int) bool {
	for _, userHint := range userHints {
		if userHint.HintID == hintID {
			return true
		}
	}
	return false
}

func (c learningController) getNextLevelHint(info hintInfo) *models.HintDB {
	for _, activityHint := range info.activityHints {
		if !c.isUsedHint(info.userHints, activityHint.ID) {
			return &activityHint
		}
	}
	return nil
}

func (c learningController) UseHint(userID int, activityID int) (*models.HintDB, error) {
	hintInfo, err := c.loadHintInfo(userID, activityID)
	if err != nil || len(hintInfo.activityHints) == 0 {
		logs.New().Error(err)
		return nil, errs.NewNotFoundError("ไม่พบคำใบ้ของกิจกรรม", "Activity Hints Not Found")
	}
	nextLevelHint := c.getNextLevelHint(hintInfo)

	if nextLevelHint == nil {
		return nil, errs.NewBadRequestError("ได้ใช้คำใบ้ทั้งหมดของกิจกรรมแล้ว", "Activity Hints Has Been Used")
	}

	if hintInfo.user.Point < nextLevelHint.PointReduce {
		return nil, errs.NewBadRequestError("แต้มไม่เพียงพอในการขอคำใบ้", "Not Enough Points")
	}

	err = c.useHintTransaction(userID, nextLevelHint.PointReduce, nextLevelHint.ID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewInternalServerError("เกิดข้อผิดพลาด", "Internal Server Error")
	}
	return nextLevelHint, nil
}

func (c learningController) useHintTransaction(userID int, reducePoint int, hintID int) error {
	var wg sync.WaitGroup
	var err error
	tx := database.NewTransaction()

	ct := models.ConcurrentTransaction{
		Concurrent: &models.Concurrent{
			Wg:  &wg,
			Err: &err,
		},
		Transaction: tx,
	}

	wg.Add(2)
	go c.updateUserPointAsyncTrasaction(&ct, userID, reducePoint, models.Mode.Reduce)
	go c.insertUserHintAsyncTransaction(&ct, userID, hintID)
	wg.Wait()
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	tx.Close()

	return nil
}

func (c learningController) updateUserPointAsyncTrasaction(ct *models.ConcurrentTransaction, userID int, updatePoint int, mode models.ChangePointMode) {
	defer ct.Concurrent.Wg.Done()
	err := c.userRepo.ChangePointTransaction(ct.Transaction, userID, updatePoint, mode)
	if err != nil {
		*ct.Concurrent.Err = err
	}
}

func (c learningController) insertUserHintAsyncTransaction(ct *models.ConcurrentTransaction, userID int, hintID int) {
	defer ct.Concurrent.Wg.Done()
	hint := models.UserHintDB{
		UserID:           userID,
		HintID:           hintID,
		CreatedTimestamp: time.Now().Local(),
	}
	_, err := c.userRepo.InsertUserHintTransaction(ct.Transaction, hint)
	if err != nil {
		*ct.Concurrent.Err = err
	}
}

func (c learningController) CheckCompletionAnswer(userID int, request models.CompletionAnswerRequest) (interface{}, error) {
	info, err := c.loadCheckAnswerInfo(*request.ActivityID, 1)
	if err != nil || info.activity == nil {
		logs.New().Error(err)
		return nil, errs.NewNotFoundError("ไม่พบกิจกรรม", "Activity Not Found")
	}

	CompletionContent := info.choice.([]models.CompletionChoiceDB)
	isCorrect := true

	if len(CompletionContent) != len(request.Answer) {
		return nil, errs.NewBadRequestError("รูปแบบของคำตอบไม่ถูกต้อง", "Invalid Answer Format")
	}

	for _, correct := range CompletionContent {
		for _, answer := range request.Answer {
			if (correct.ID == *answer.ID) && (correct.Content != *answer.Content) {
				isCorrect = false
				break
			}
		}
	}

	response := models.AnswerResponse{
		ActivityID: *request.ActivityID,
		IsCorrect:  isCorrect,
	}

	if isCorrect {
		err = c.finishActivityTrasaction(userID, *request.ActivityID, info.activity.Point)
		if err != nil && !utils.NewHelper().IsSqlDuplicateError(err) {
			logs.New().Error(err)
			return nil, errs.NewInternalServerError("เกิดข้อผิดพลาดในการบันทึกกิจกรรม", "Saved Activity Failed")
		}
	}

	return response, nil
}

func (c learningController) loadContentAsync(concurrent *models.Concurrent, contentID int, content **models.ContentDB) {
	defer concurrent.Wg.Done()
	result, err := c.learningRepo.GetContent(contentID)
	if err != nil {
		*concurrent.Err = err
	}
	*content = result
}

func (c learningController) loadRoadmapInfo(userID int, contentID int) (roadmapInfo, error) {
	var wg sync.WaitGroup
	var err error
	concurrent := models.Concurrent{Wg: &wg, Err: &err}
	var content *models.ContentDB
	contentActivity := make([]models.ActivityDB, 0)
	learningProgression := make([]models.LearningProgressionDB, 0)
	wg.Add(3)
	go c.loadLearningProgressionAsync(&concurrent, &learningProgression, userID)
	go c.loadContentActivityAsync(&concurrent, contentID, &contentActivity)
	go c.loadContentAsync(&concurrent, contentID, &content)
	wg.Wait()
	info := roadmapInfo{
		contentActivity: contentActivity,
		progression:     learningProgression,
		content:         content,
	}
	return info, err
}

func (c learningController) isLearnedActivity(progression []models.LearningProgressionDB, activityID int) bool {
	for _, v := range progression {
		if v.ActivityID == activityID {
			return true
		}
	}
	return false
}

func (c learningController) prepareRoadmap(info roadmapInfo) *models.RoadmapResponse {
	items := make([]models.RoadmapItem, 0)

	for _, v := range info.contentActivity {
		items = append(items, models.RoadmapItem{
			ActivityID: v.ID,
			IsLearned:  c.isLearnedActivity(info.progression, v.ID),
			Order:      v.Order,
		})
	}

	res := models.RoadmapResponse{
		ContentID:   info.content.ID,
		ContentName: info.content.Name,
		Items:       items,
	}

	return &res
}

func (c learningController) GetContentRoadmap(userID int, contentID int) (*models.RoadmapResponse, error) {
	info, err := c.loadRoadmapInfo(userID, contentID)
	if err != nil {
		logs.New().Error(err)
		return nil, errs.NewNotFoundError("ไม่พบเนื้อหา", "Content Not Found")
	}
	roadmap := c.prepareRoadmap(info)
	return roadmap, nil
}
