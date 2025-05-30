package httpv1

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/horiondreher/go-parking-lot/users/internal/adapters/grpc/proto"
	"github.com/horiondreher/go-parking-lot/users/internal/adapters/http/httputils"
	"github.com/horiondreher/go-parking-lot/users/internal/domain/domainerr"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ReservartionResponse struct {
	ID            uuid.UUID  `json:"id"`
	Type          string     `json:"type"`
	RemainingTime *time.Time `json:"remaining_time,omitempty"`
}

func (adapter *HTTPAdapter) GetReservation(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	usersServiceURL := "http://go-parking-lot-parkings-service:8080/api/v1/parkings"

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, usersServiceURL, nil)
	if err != nil {
		log.Err(err).Msg("error creating request")
		return domainerr.NewInternalError(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Msg("error sending request")
		return domainerr.NewInternalError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Msgf("unexpected status code: %d", resp.StatusCode)
		return domainerr.NewDomainError(resp.StatusCode, domainerr.NotFoundError, "client error", nil)
	}

	var reservation ReservartionResponse
	if err := json.NewDecoder(resp.Body).Decode(&reservation); err != nil {
		log.Err(err).Msg("error decoding response")
		return domainerr.NewInternalError(err)
	}

	domainErr := httputils.Encode(w, r, http.StatusOK, reservation)
	if domainErr != nil {
		return domainErr
	}

	return nil
}

func (adapter *HTTPAdapter) GetReservationgRPC(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	conn, err := grpc.NewClient("go-parking-lot-parkings-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Err(err).Msg("failed to connect")
		return domainerr.NewInternalError(err)
	}
	defer conn.Close()

	client := proto.NewReservationServiceClient(conn)
	response, err := client.GetReservation(r.Context(), &proto.GetReservationRequest{})
	if err != nil {
		log.Err(err).Msg("failed to connect ")
		return domainerr.NewInternalError(err)
	}

	reservationID, _ := uuid.Parse(response.Id)
	reservarionRemaimingTime := response.RemainingTime.AsTime()

	reservation := ReservartionResponse{
		ID:            reservationID,
		Type:          response.Type,
		RemainingTime: &reservarionRemaimingTime,
	}
	domainErr := httputils.Encode(w, r, http.StatusOK, reservation)
	if domainErr != nil {
		return domainErr
	}

	return nil
}
