package request

import (
	"database-camp/internal/errs"
)

type CheckAnswerRequest struct {
	ActivityID     *int        `json:"activity_id"`
	ActivityTypeID *int        `json:"activity_type_id"`
	Answer         interface{} `json:"answer"`
}

func (r CheckAnswerRequest) Validate() error {
	if r.ActivityID == nil {
		return errs.ErrActivittyIDNotFound
	} else if r.ActivityTypeID == nil {
		return errs.ErrActivittyTypeIDNotFound
	} else if r.Answer == nil {
		return errs.ErrAnswerNotFound
	} else {
		return nil
	}
}

type PeerReviewRequest struct {
	ERAnswerID *int     `json:"er_answer_id"`
	Reviews    []string `json:"reviews"`
}

func (r PeerReviewRequest) Validate() error {
	if r.ERAnswerID == nil {
		return errs.ErrERAnswerIDNotFound
	} else {
		return nil
	}
}
