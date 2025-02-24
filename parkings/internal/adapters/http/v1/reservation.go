package v1

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/horiondreher/go-parking-lot/parkings/internal/adapters/http/httputils"
	"github.com/horiondreher/go-parking-lot/parkings/internal/domain/domainerr"
)

type ReservartionResponse struct {
	ID            uuid.UUID  `json:"id"`
	Type          string     `json:"type"`
	RemainingTime *time.Time `json:"remaining_time,omitempty"`
}

func (adapter *HTTPAdapter) GetReservation(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	timeRemaining := time.Now().Add(2 * time.Hour)
	err := httputils.Encode(w, r, http.StatusOK, ReservartionResponse{
		ID:            uuid.New(),
		Type:          "car",
		RemainingTime: &timeRemaining,
	})
	if err != nil {
		return err
	}

	return nil
}
