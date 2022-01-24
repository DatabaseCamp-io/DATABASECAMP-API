package loaders

import (
	"database-camp/internal/models/entities/content"
	"database-camp/internal/models/entities/user"
	"database-camp/internal/repositories"
	"sync"
)

type spiderDataLoader struct {
	userRepo     repositories.UserRepository
	learningRepo repositories.LearningRepository

	spiderDataset user.SpiderDataset
	contentGroups content.ContentGroups
}

func NewSpiderDataLoader(learningRepo repositories.LearningRepository, userRepo repositories.UserRepository) *spiderDataLoader {
	return &spiderDataLoader{learningRepo: learningRepo, userRepo: userRepo}
}

func (l *spiderDataLoader) GetSpiderDataset() user.SpiderDataset {
	return l.spiderDataset
}

func (l *spiderDataLoader) GetContentGroups() content.ContentGroups {
	return l.contentGroups
}

func (l *spiderDataLoader) Load(userID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := Concurrent{Wg: &wg, Err: &err}
	wg.Add(2)
	go l.loadSpiderDatasetAsync(&concurrent, userID)
	go l.loadContentGroupsAsync(&concurrent)
	wg.Wait()
	return err
}

func (l *spiderDataLoader) loadSpiderDatasetAsync(concurrent *Concurrent, id int) {
	defer concurrent.Wg.Done()
	var err error
	l.spiderDataset, err = l.userRepo.GetSpiderDataset(id)
	if err != nil {
		*concurrent.Err = err
	}
}

func (l *spiderDataLoader) loadContentGroupsAsync(concurrent *Concurrent) {
	defer concurrent.Wg.Done()
	var err error
	l.contentGroups, err = l.learningRepo.GetContentGroups()
	if err != nil {
		*concurrent.Err = err
	}
}
