package request

// request.learning.go
/**
 * 	This file is a part of models, used to collect request of learning
 */

import "DatabaseCamp/errs"

/**
 * 	This class represent multiple choice answer request
 */
type MultipleChoiceAnswerRequest struct {
	ActivityID *int `json:"activity_id"`
	Answer     *int `json:"answer"`
}

/**
 * Validate multiple choice answer request
 *
 * @return the error of validating request
 */
func (r MultipleChoiceAnswerRequest) Validate() error {
	if r.ActivityID == nil {
		return errs.NewBadRequestError("ไม่พบไอดีของกิจกรรมในคำร้องขอ", "Activity ID Not Found")
	} else if r.Answer == nil {
		return errs.NewBadRequestError("ไม่พบคำตอบในคำร้องขอ", "Answer Not Found")
	}
	return nil
}

/**
 * 	This class represent matching choice answer request
 */
type MatchingChoiceAnswerRequest struct {
	ActivityID *int              `json:"activity_id"`
	Answer     []PairItemRequest `json:"answer"`
}

/**
 * Validate pair item request
 *
 * @param 	pairItem 	PairItemRequest
 *
 * @return the error of validating request
 */
func (r MatchingChoiceAnswerRequest) validatePairItem(pairItem PairItemRequest) error {
	if pairItem.Item1 == nil || pairItem.Item2 == nil {
		return errs.NewBadRequestError("ไม่พบคำตอบในคำร้องขอ", "Answer Not Found")
	}
	return nil
}

/**
 * Validate matching choice answer request
 *
 * @return the error of validating request
 */
func (r MatchingChoiceAnswerRequest) Validate() error {
	if r.ActivityID == nil {
		return errs.NewBadRequestError("ไม่พบไอดีของกิจกรรมในคำร้องขอ", "Activity ID Not Found")
	} else if len(r.Answer) == 0 {
		return errs.NewBadRequestError("ไม่พบคำตอบในคำร้องขอ", "Answer Not Found")
	} else {
		for _, v := range r.Answer {
			e := r.validatePairItem(v)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

/**
 * 	This class represent completion choice answer request
 */
type CompletionAnswerRequest struct {
	ActivityID *int                 `json:"activity_id"`
	Answer     []PairContentRequest `json:"answer"`
}

/**
 * Validate pair content request
 *
 * @param pairContent PairContent request
 *
 * @return the error of validating request
 */
func (r CompletionAnswerRequest) validatePairItem(pairContent PairContentRequest) error {
	if pairContent.Content == nil || pairContent.ID == nil {
		return errs.NewBadRequestError("ไม่พบคำตอบในคำร้องขอ", "Answer Not Found")
	}
	return nil
}

/**
 * Validate completion choice answer request
 *
 * @return the error of validating request
 */
func (r CompletionAnswerRequest) Validate() error {
	if r.ActivityID == nil {
		return errs.NewBadRequestError("ไม่พบไอดีของกิจกรรมในคำร้องขอ", "Activity ID Not Found")
	} else if len(r.Answer) == 0 {
		return errs.NewBadRequestError("ไม่พบคำตอบในคำร้องขอ", "Answer Not Found")
	} else {
		for _, v := range r.Answer {
			e := r.validatePairItem(v)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

// Model for matching choice answer request
type PairItemRequest struct {
	Item1 *string `json:"item1"`
	Item2 *string `json:"item2"`
}

/**
 * Validate pair item request request
 *
 * @return the error of validating request
 */
func (m PairItemRequest) Validate() error {
	if m.Item1 == nil || m.Item2 == nil {
		return errs.NewBadRequestError("ไม่พบเนื้อหาของคำตอบในคำร้องขอ", "Content Answer Not Found")
	}
	return nil
}

// Model for completion choice answer request
type PairContentRequest struct {
	ID      *int    `json:"completion_choice_id"`
	Content *string `json:"content"`
}

/**
 * Validate pair content request request
 *
 * @return the error of validating request
 */
func (m PairContentRequest) Validate() error {
	if m.ID == nil {
		return errs.NewBadRequestError("ไม่พบไอดีของช้อยในคำร้องขอ", "Choice ID Not Found")
	} else if m.Content == nil {
		return errs.NewBadRequestError("ไม่พบเนื้อหาของคำตอบในคำร้องขอ", "Content Answer Not Found")
	}
	return nil
}
