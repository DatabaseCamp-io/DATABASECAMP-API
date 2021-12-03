package request

import "DatabaseCamp/errs"

type MultipleChoiceAnswerRequest struct {
	ActivityID *int `json:"activity_id"`
	Answer     *int `json:"answer"`
}

// Check multiple choice answer from activity id
func (r MultipleChoiceAnswerRequest) Validate() error {
	if r.ActivityID == nil {
		return errs.NewBadRequestError("ไม่พบไอดีของกิจกรรมในคำร้องขอ", "Activity ID Not Found")
	} else if r.Answer == nil {
		return errs.NewBadRequestError("ไม่พบคำตอบในคำร้องขอ", "Answer Not Found")
	}
	return nil
}

type MatchingChoiceAnswerRequest struct {
	ActivityID *int              `json:"activity_id"`
	Answer     []PairItemRequest `json:"answer"`
}

// Check matching choice answer existence
func (r MatchingChoiceAnswerRequest) validatePairItem(pairItem PairItemRequest) error {
	if pairItem.Item1 == nil || pairItem.Item2 == nil {
		return errs.NewBadRequestError("ไม่พบคำตอบในคำร้องขอ", "Answer Not Found")
	}
	return nil
}

// Check matching choice answer from activity id
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

type CompletionAnswerRequest struct {
	ActivityID *int                 `json:"activity_id"`
	Answer     []PairContentRequest `json:"answer"`
}

// Check completion choice answer existence
func (r CompletionAnswerRequest) validatePairItem(pairContent PairContentRequest) error {
	if pairContent.Content == nil || pairContent.ID == nil {
		return errs.NewBadRequestError("ไม่พบคำตอบในคำร้องขอ", "Answer Not Found")
	}
	return nil
}

// Check completion choice answer from activity id
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

type PairItemRequest struct {
	Item1 *string `json:"item1"`
	Item2 *string `json:"item2"`
}

// Check pair item existence
func (m PairItemRequest) Validate() error {
	if m.Item1 == nil || m.Item2 == nil {
		return errs.NewBadRequestError("ไม่พบเนื้อหาของคำตอบในคำร้องขอ", "Content Answer Not Found")
	}
	return nil
}

type PairContentRequest struct {
	ID      *int    `json:"completion_choice_id"`
	Content *string `json:"content"`
}

// Check pair content existence
func (m PairContentRequest) Validate() error {
	if m.ID == nil {
		return errs.NewBadRequestError("ไม่พบไอดีของช้อยในคำร้องขอ", "Choice ID Not Found")
	} else if m.Content == nil {
		return errs.NewBadRequestError("ไม่พบเนื้อหาของคำตอบในคำร้องขอ", "Content Answer Not Found")
	}
	return nil
}
