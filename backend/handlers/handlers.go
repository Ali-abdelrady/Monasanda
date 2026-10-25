package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"mosanda-vacation-attendace/models"
	"mosanda-vacation-attendace/pkg/xhttputil"
	"mosanda-vacation-attendace/service"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type authenticatedUserContextKey struct{}

type Handler struct {
	service *service.Service
}

func NewHandlers(parService *service.Service) *Handler {
	return &Handler{service: parService}
}

func (objHandler *Handler) RequireSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(parResponse http.ResponseWriter, parRequest *http.Request) {
		psCookie, pvError := parRequest.Cookie("session_id")
		if pvError != nil {
			parResponse.WriteHeader(http.StatusUnauthorized)
			xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0300", "Authentication required")
			return
		}

		psSessionID, pvError := uuid.Parse(psCookie.Value)
		if pvError != nil {
			parResponse.WriteHeader(http.StatusUnauthorized)
			xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0301", "Invalid session")
			return
		}

		psUser, pvError := objHandler.service.GetUserFromSession(psSessionID)
		if errors.Is(pvError, sql.ErrNoRows) {
			parResponse.WriteHeader(http.StatusUnauthorized)
			xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0302", "Session is invalid or expired")
			return
		}
		if pvError != nil {
			parResponse.WriteHeader(http.StatusInternalServerError)
			xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0303", "Unable to validate session")
			return
		}

		psContext := context.WithValue(parRequest.Context(), authenticatedUserContextKey{}, psUser)
		next.ServeHTTP(parResponse, parRequest.WithContext(psContext))
	})
}

func authenticatedUserFromRequest(parRequest *http.Request) *models.AuthenticatedUser {
	psUser, _ := parRequest.Context().Value(authenticatedUserContextKey{}).(*models.AuthenticatedUser)
	return psUser
}

func (objHandler *Handler) GetCurrentUser(parResponse http.ResponseWriter, parRequest *http.Request) {
	psUser := authenticatedUserFromRequest(parRequest)
	if psUser == nil {
		parResponse.WriteHeader(http.StatusUnauthorized)
		xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0305", "Authentication required")
		return
	}

	xhttputil.SendResponse(parResponse, parRequest, map[string]any{
		"hasError": false,
		"user":     psUser,
	})
}

func (objHandler *Handler) Login(parResponse http.ResponseWriter, parRequest *http.Request) {
	if parRequest.Method != "POST" {
		xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0100", "Method not allowed")
		return
	}

	var psUser *models.User

	// 1. Decode requet body
	pvError := json.NewDecoder(parRequest.Body).Decode(&psUser)
	if pvError != nil {
		fmt.Println(pvError)
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0101", pvError)
		return
	}

	// 2. Validate data
	pvError = psUser.Validate()
	if pvError != nil {
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0102", pvError)
		return
	}

	// Authenticate psUser
	pvUserID, pvError := objHandler.service.AuthenticateUser(psUser.Email, psUser.Password)
	if pvError != nil {
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0103", pvError)
		return
	}

	// Start New Session
	psSession, pvError := objHandler.service.StartSession(pvUserID)
	if pvError != nil {
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0104", pvError)
		return
	}

	http.SetCookie(parResponse, &http.Cookie{
		Name:     "session_id",
		Value:    psSession.ID.String(),
		Expires:  psSession.ExpiresAt,
		HttpOnly: true,  // Prevents JavaScript access (XSS defense)
		Secure:   false, // Set to true in production when using HTTPS
		Path:     "/",
		SameSite: http.SameSiteLaxMode, // CSRF defense
	})

	pmResponseData := map[string]any{
		"hasError": false,
	}

	xhttputil.SendResponse(parResponse, parRequest, pmResponseData)
}

func (objHandler *Handler) CreateVacationBatch(parResponse http.ResponseWriter, parRequest *http.Request) {
	var psBatch *models.VacationBatch
	if pvError := json.NewDecoder(parRequest.Body).Decode(&psBatch); pvError != nil {
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0200", pvError)
		return
	}
	psUser := authenticatedUserFromRequest(parRequest)
	if psUser == nil {
		parResponse.WriteHeader(http.StatusUnauthorized)
		xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0304", "Authentication required")
		return
	}

	psBatch.CreatedBy = psUser.ID
	psBatch.CreatedAt = time.Now()
	psBatch.UpdatedAt = time.Now()

	if pvError := psBatch.Validate(); pvError != nil {
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0201", pvError)
		return
	}

	psCreatedBatch, pvError := objHandler.service.CreateVacationBatch(psBatch)
	if pvError != nil {
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0202", pvError)
		return
	}
	xhttputil.SendResponse(parResponse, parRequest, map[string]any{
		"hasError":      false,
		"vacationBatch": psCreatedBatch,
	})
}

func (objHandler *Handler) GetVacationBatches(parResponse http.ResponseWriter, parRequest *http.Request) {
	pvBatches, pvError := objHandler.service.GetVacationBatches()
	if pvError != nil {
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0203", pvError)
		return
	}
	xhttputil.SendResponse(parResponse, parRequest, map[string]any{
		"hasError":        false,
		"vacationBatches": pvBatches,
	})
}

func (objHandler *Handler) GetVacationBatch(parResponse http.ResponseWriter, parRequest *http.Request) {
	pvBatchID, pvError := strconv.Atoi(parRequest.PathValue("id"))
	if pvError != nil || pvBatchID <= 0 {
		xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0204", "Invalid vacation batch ID")
		return
	}

	psBatch, pvError := objHandler.service.GetVacationBatch(pvBatchID)
	if pvError != nil {
		if errors.Is(pvError, sql.ErrNoRows) {
			xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0205", "Vacation batch not found")
			return
		}
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0206", pvError)
		return
	}
	xhttputil.SendResponse(parResponse, parRequest, map[string]any{
		"hasError":      false,
		"vacationBatch": psBatch,
	})
}

func (objHandler *Handler) UpdateVacationBatch(parResponse http.ResponseWriter, parRequest *http.Request) {
	pvBatchID, pvError := strconv.Atoi(parRequest.PathValue("id"))
	if pvError != nil || pvBatchID <= 0 {
		xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0207", "Invalid vacation batch ID")
		return
	}

	var psBatch *models.VacationBatch
	if pvError = json.NewDecoder(parRequest.Body).Decode(&psBatch); pvError != nil {
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0208", pvError)
		return
	}
	if pvError = psBatch.Validate(); pvError != nil {
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0209", pvError)
		return
	}

	psBatch.CreatedAt = time.Now()
	psBatch.UpdatedAt = time.Now()

	psUpdatedBatch, pvError := objHandler.service.UpdateVacationBatch(pvBatchID, psBatch)
	if pvError != nil {
		if errors.Is(pvError, sql.ErrNoRows) {
			xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0210", "Vacation batch not found")
			return
		}
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0211", pvError)
		return
	}
	xhttputil.SendResponse(parResponse, parRequest, map[string]any{
		"hasError":      false,
		"vacationBatch": psUpdatedBatch,
	})
}

func (objHandler *Handler) DeleteVacationBatch(parResponse http.ResponseWriter, parRequest *http.Request) {
	pvBatchID, pvError := strconv.Atoi(parRequest.PathValue("id"))
	if pvError != nil || pvBatchID <= 0 {
		xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0212", "Invalid vacation batch ID")
		return
	}
	if pvError = objHandler.service.DeleteVacationBatch(pvBatchID); pvError != nil {
		if errors.Is(pvError, sql.ErrNoRows) {
			xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0213", "Vacation batch not found")
			return
		}
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0214", pvError)
		return
	}
	xhttputil.SendResponse(parResponse, parRequest, map[string]any{"hasError": false})
}

func (objHandler *Handler) CreateAttendance(parResponse http.ResponseWriter, parRequest *http.Request) {
	var psAttendance *models.Attendance
	if pvError := json.NewDecoder(parRequest.Body).Decode(&psAttendance); pvError != nil {
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0400", pvError)
		return
	}
	psUser := authenticatedUserFromRequest(parRequest)
	if psUser == nil {
		parResponse.WriteHeader(http.StatusUnauthorized)
		xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0401", "Authentication required")
		return
	}
	if psAttendance.Status == "" {
		psAttendance.Status = models.AttendanceStatusAbsent
	}
	psAttendance.CreatedBy = psUser.ID
	setAttendanceCheckIn(psAttendance, psUser.ID)
	if pvError := psAttendance.Validate(); pvError != nil {
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0402", pvError)
		return
	}

	psCreatedAttendance, pvError := objHandler.service.CreateAttendance(psAttendance)
	if pvError != nil {
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0403", pvError)
		return
	}
	xhttputil.SendResponse(parResponse, parRequest, map[string]any{
		"hasError":   false,
		"attendance": psCreatedAttendance,
	})
}

func (objHandler *Handler) GetAttendances(parResponse http.ResponseWriter, parRequest *http.Request) {
	pvAttendances, pvError := objHandler.service.GetAttendances()
	if pvError != nil {
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0404", pvError)
		return
	}
	xhttputil.SendResponse(parResponse, parRequest, map[string]any{
		"hasError":    false,
		"attendances": pvAttendances,
	})
}

func (objHandler *Handler) GetAttendance(parResponse http.ResponseWriter, parRequest *http.Request) {
	pvAttendanceID, pvError := strconv.Atoi(parRequest.PathValue("id"))
	if pvError != nil || pvAttendanceID <= 0 {
		xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0405", "Invalid attendance ID")
		return
	}
	psAttendance, pvError := objHandler.service.GetAttendance(pvAttendanceID)
	if pvError != nil {
		if errors.Is(pvError, sql.ErrNoRows) {
			xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0406", "Attendance record not found")
			return
		}
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0407", pvError)
		return
	}
	xhttputil.SendResponse(parResponse, parRequest, map[string]any{
		"hasError":   false,
		"attendance": psAttendance,
	})
}

func (objHandler *Handler) UpdateAttendance(parResponse http.ResponseWriter, parRequest *http.Request) {
	pvAttendanceID, pvError := strconv.Atoi(parRequest.PathValue("id"))
	if pvError != nil || pvAttendanceID <= 0 {
		xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0408", "Invalid attendance ID")
		return
	}
	var psAttendance *models.Attendance
	if pvError = json.NewDecoder(parRequest.Body).Decode(&psAttendance); pvError != nil {
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0409", pvError)
		return
	}
	psUser := authenticatedUserFromRequest(parRequest)
	if psUser == nil {
		parResponse.WriteHeader(http.StatusUnauthorized)
		xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0410", "Authentication required")
		return
	}
	if psAttendance.Status == "" {
		psAttendance.Status = models.AttendanceStatusAbsent
	}
	setAttendanceCheckIn(psAttendance, psUser.ID)
	if pvError = psAttendance.Validate(); pvError != nil {
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0411", pvError)
		return
	}

	psUpdatedAttendance, pvError := objHandler.service.UpdateAttendance(pvAttendanceID, psAttendance)
	if pvError != nil {
		if errors.Is(pvError, sql.ErrNoRows) {
			xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0412", "Attendance record not found")
			return
		}
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0413", pvError)
		return
	}
	xhttputil.SendResponse(parResponse, parRequest, map[string]any{
		"hasError":   false,
		"attendance": psUpdatedAttendance,
	})
}

func (objHandler *Handler) DeleteAttendance(parResponse http.ResponseWriter, parRequest *http.Request) {
	pvAttendanceID, pvError := strconv.Atoi(parRequest.PathValue("id"))
	if pvError != nil || pvAttendanceID <= 0 {
		xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0414", "Invalid attendance ID")
		return
	}
	if pvError = objHandler.service.DeleteAttendance(pvAttendanceID); pvError != nil {
		if errors.Is(pvError, sql.ErrNoRows) {
			xhttputil.HandleResponseErrorTextMsg(parResponse, parRequest, "HND0415", "Attendance record not found")
			return
		}
		xhttputil.HandleResponseErrorMsg(parResponse, parRequest, "HND0416", pvError)
		return
	}
	xhttputil.SendResponse(parResponse, parRequest, map[string]any{"hasError": false})
}

func setAttendanceCheckIn(parAttendance *models.Attendance, parUserID int) {
	if parAttendance.Status == models.AttendanceStatusPresent {
		pvCheckedInAt := time.Now()
		pvCheckedInBy := parUserID
		parAttendance.CheckedInAt = &pvCheckedInAt
		parAttendance.CheckedInByID = &pvCheckedInBy
		return
	}
	parAttendance.CheckedInAt = nil
	parAttendance.CheckedInByID = nil
}
