package loaders

import (
	"database-camp/internal/models/entities/content"
	"database-camp/internal/models/entities/user"
	"database-camp/internal/repositories"
	"sync"
)

type recommendLoader struct {
	userRepo     repositories.UserRepository
	learningRepo repositories.LearningRepository

	preTestResults user.PreTestResults
	contentGroups  content.ContentGroups
}

func NewRecommendLoader(learningRepo repositories.LearningRepository, userRepo repositories.UserRepository) *recommendLoader {
	return &recommendLoader{learningRepo: learningRepo, userRepo: userRepo}
}

func (l *recommendLoader) GetPreTestResults() user.PreTestResults {
	return l.preTestResults
}

func (l *recommendLoader) GetContentGroups() content.ContentGroups {
	return l.contentGroups
}

func (l *recommendLoader) Load(userID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := Concurrent{Wg: &wg, Err: &err}
	wg.Add(2)
	go l.loadPreTestResultsAsync(&concurrent, userID)
	go l.loadContentGroupsAsync(&concurrent)
	wg.Wait()
	return err
}

func (l *recommendLoader) loadPreTestResultsAsync(concurrent *Concurrent, id int) {
	defer concurrent.Wg.Done()
	var err error
	l.preTestResults, err = l.userRepo.GetPreTestResults(id)
	if err != nil {
		*concurrent.Err = err
	}
}

func (l *recommendLoader) loadContentGroupsAsync(concurrent *Concurrent) {
	defer concurrent.Wg.Done()
	var err error
	l.contentGroups, err = l.learningRepo.GetContentGroups()
	if err != nil {
		*concurrent.Err = err
	}
}
