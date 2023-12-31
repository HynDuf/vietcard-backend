package response

import "vietcard-backend/internal/domain/entity"

type CreateDeckResponse struct {
	Deck    entity.Deck `json:"deck"`
}

type UpdateDeckResponse struct {
	Success bool        `json:"success"`
	Deck    entity.Deck `json:"deck"`
}

type CopyDeckResponse struct {
	Deck entity.DeckWithCards `json:"deck"`
	DeckReview entity.DeckWithReviewCards `json:"deck_review"`
}

type DeleteDeckResponse struct {
	Success bool `json:"success"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
}
