package loaders

import (
	"database-camp/internal/models/entities/badge"
	"database-camp/internal/models/entities/content"
	"database-camp/internal/models/entities/user"
	"database-camp/internal/repositories"
	"sync"
)

type profileLoader struct {
	userRepo     repositories.UserRepository
	learningRepo repositories.LearningRepository

	spiderDataset user.SpiderDataset
	contentGroups content.ContentGroups
	profile       *user.Profile
	badges        []badge.Badge
	userBadges    []badge.UserBadge
}

func NewProfileLoader(learningRepo repositories.LearningRepository, userRepo repositories.UserRepository) *profileLoader {
	return &profileLoader{learningRepo: learningRepo, userRepo: userRepo}
}

func (l *profileLoader) GetSpiderDataset() user.SpiderDataset {
	return l.spiderDataset
}

func (l *profileLoader) GetContentGroups() content.ContentGroups {
	return l.contentGroups
}

func (l *profileLoader) GetProfile() *user.Profile {
	return l.profile
}

func (l *profileLoader) GetUserBadges() []badge.UserBadge {
	return l.userBadges
}

func (l *profileLoader) GetBadges() []badge.Badge {
	return l.badges
}

func (l *profileLoader) Load(userID int) error {
	var wg sync.WaitGroup
	var err error
	concurrent := Concurrent{Wg: &wg, Err: &err}
	wg.Add(5)
	go l.loadSpiderDatasetAsync(&concurrent, userID)
	go l.loadContentGroupsAsync(&concurrent)
	go l.loadProfileAsync(&concurrent, userID)
	go l.loadBadgesAsync(&concurrent)
	go l.loadUserBadgesAsync(&concurrent, userID)
	wg.Wait()
	return err
}

func (l *profileLoader) loadSpiderDatasetAsync(concurrent *Concurrent, id int) {
	defer concurrent.Wg.Done()
	var err error
	l.spiderDataset, err = l.userRepo.GetSpiderDataset(id)
	if err != nil {
		*concurrent.Err = err
	}
}

func (l *profileLoader) loadContentGroupsAsync(concurrent *Concurrent) {
	defer concurrent.Wg.Done()
	var err error
	l.contentGroups, err = l.learningRepo.GetContentGroups()
	if err != nil {
		*concurrent.Err = err
	}
}

func (l *profileLoader) loadProfileAsync(concurrent *Concurrent, id int) {
	defer concurrent.Wg.Done()
	var err error
	l.profile, err = l.userRepo.GetProfile(id)
	if err != nil {
		*concurrent.Err = err
	}
}

func (l *profileLoader) loadBadgesAsync(concurrent *Concurrent) {
	defer concurrent.Wg.Done()
	var err error
	l.badges, err = l.userRepo.GetAllBadge()
	if err != nil {
		*concurrent.Err = err
	}
}

func (l *profileLoader) loadUserBadgesAsync(concurrent *Concurrent, id int) {
	defer concurrent.Wg.Done()
	var err error
	l.userBadges, err = l.userRepo.GetUserBadge(id)
	if err != nil {
		*concurrent.Err = err
	}
}
