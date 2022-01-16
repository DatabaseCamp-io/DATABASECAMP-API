package loaders

import (
	"database-camp/internal/models/entities/activity"
	"database-camp/internal/models/entities/exam"
	"database-camp/internal/repositories"
	"sync"
)

type activityExamLoader struct {
	learningRepo repositories.LearningRepository

	activities []exam.Activity
	mutex      sync.Mutex
}

func NewActivityExamLoader(learningRepo repositories.LearningRepository) *activityExamLoader {
	return &activityExamLoader{learningRepo: learningRepo}
}

func (l *activityExamLoader) GetActivities() exam.Activities {
	return l.activities
}

func (l *activityExamLoader) Load(examActivities []exam.ExamActivity) error {
	var wg sync.WaitGroup
	var err error

	wg.Add(len(examActivities))

	for _, activity := range examActivities {
		go func(activityID int, typeID int) {
			defer wg.Done()
			l.loadActivity(activityID, typeID, &err)
		}(activity.ActivityID, activity.ActivityTypeID)
	}

	wg.Wait()
	return err
}

func (l *activityExamLoader) loadActivity(activityID int, typeID int, err *error) {
	var wg sync.WaitGroup
	var _activity *activity.Activity
	var choices activity.Choices

	wg.Add(2)

	go func() {
		defer wg.Done()
		var e error
		_activity, e = l.learningRepo.GetActivity(activityID)
		if e != nil {
			*err = e
		}
	}()

	go func() {
		defer wg.Done()
		var e error
		choices, e = l.learningRepo.GetActivityChoices(activityID, typeID)
		if e != nil {
			*err = e
		}
	}()

	wg.Wait()

	l.mutex.Lock()

	l.activities = append(l.activities, exam.Activity{
		Activity: *_activity,
		Choices:  choices,
	})

	l.mutex.Unlock()
}
