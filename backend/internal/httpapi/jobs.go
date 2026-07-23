package httpapi

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	applicationjobs "github.com/saral-gupta7/recode/backend/internal/application/jobs"
	"github.com/saral-gupta7/recode/backend/internal/httpapi/response"
	"github.com/saral-gupta7/recode/backend/internal/job"
	"github.com/saral-gupta7/recode/backend/internal/storage"
	"github.com/saral-gupta7/recode/backend/internal/task"
)

const multipartOverheadAllowance int64 = 1 << 20

type jobsHandler struct {
	service  *applicationjobs.Service
	maxBytes int64
}

type jobResponse struct {
	ID          string  `json:"id"`
	Operation   string  `json:"operation"`
	Status      string  `json:"status"`
	Progress    int     `json:"progress"`
	Attempt     int     `json:"attempt"`
	ExpiresAt   *string `json:"expires_at,omitempty"`
	ResultReady bool    `json:"result_ready"`
	FailureCode string  `json:"failure_code,omitempty"`
	OwnerToken  string  `json:"owner_token,omitempty"`
}

func (h jobsHandler) create(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(
		c.Writer,
		c.Request.Body,
		h.maxBytes+multipartOverheadAllowance,
	)

	operation, err := job.ParseOperation(c.PostForm("operation"))
	if err != nil {
		writeJobError(c, task.ErrInvalidInput)
		return
	}

	options, err := task.DecodeOptions([]byte(c.PostForm("options")))
	if err != nil {
		writeJobError(c, err)
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		writeJobError(c, task.ErrInvalidInput)
		return
	}
	if c.Request.MultipartForm != nil {
		defer c.Request.MultipartForm.RemoveAll()
	}

	file, err := fileHeader.Open()
	if err != nil {
		writeJobError(c, err)
		return
	}
	defer file.Close()

	created, err := h.service.Create(c.Request.Context(), applicationjobs.CreateInput{
		Operation:        operation,
		Options:          options,
		OriginalFilename: fileHeader.Filename,
		MIMEType:         fileHeader.Header.Get("Content-Type"),
		Content:          file,
	})
	if err != nil {
		writeJobError(c, err)
		return
	}

	body := toJobResponse(created.Record)
	body.OwnerToken = created.OwnerToken
	c.JSON(http.StatusAccepted, body)
}

func (h jobsHandler) get(c *gin.Context) {
	record, err := h.service.Get(c.Request.Context(), c.Param("jobID"), bearerToken(c))
	if err != nil {
		writeJobError(c, err)
		return
	}
	c.JSON(http.StatusOK, toJobResponse(record))
}

func (h jobsHandler) cancel(c *gin.Context) {
	record, err := h.service.Cancel(c.Request.Context(), c.Param("jobID"), bearerToken(c))
	if err != nil {
		writeJobError(c, err)
		return
	}
	c.JSON(http.StatusOK, toJobResponse(record))
}

func (h jobsHandler) result(c *gin.Context) {
	record, reader, err := h.service.OpenResult(
		c.Request.Context(),
		c.Param("jobID"),
		bearerToken(c),
	)
	if err != nil {
		writeJobError(c, err)
		return
	}
	defer reader.Close()

	extension := task.OutputExtension(record.Job.Operation(), record.Options)
	c.Header("Content-Type", outputContentType(extension))
	c.Header(
		"Content-Disposition",
		fmt.Sprintf(`attachment; filename="recode-%s.%s"`, record.Job.ID(), extension),
	)
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, reader)
}

func (h jobsHandler) delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("jobID"), bearerToken(c)); err != nil {
		writeJobError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func toJobResponse(record task.Record) jobResponse {
	snapshot := record.Job.Snapshot()
	responseBody := jobResponse{
		ID:          snapshot.ID,
		Operation:   string(snapshot.Operation),
		Status:      string(snapshot.Status),
		Progress:    snapshot.Progress,
		Attempt:     snapshot.Attempt,
		ResultReady: snapshot.Status == job.StatusCompleted && snapshot.ResultKey != "",
		FailureCode: snapshot.FailureCode,
	}
	if !snapshot.ExpiresAt.IsZero() {
		value := snapshot.ExpiresAt.Format("2006-01-02T15:04:05.999999999Z07:00")
		responseBody.ExpiresAt = &value
	}
	return responseBody
}

func bearerToken(c *gin.Context) string {
	const prefix = "Bearer "
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}

func writeJobError(c *gin.Context, err error) {
	var maxBytesError *http.MaxBytesError
	if errors.As(err, &maxBytesError) {
		response.AbortWithError(c, http.StatusRequestEntityTooLarge, "file_too_large", "The uploaded file is too large.")
		return
	}

	switch {
	case errors.Is(err, task.ErrNotFound):
		response.AbortWithError(c, http.StatusNotFound, "job_not_found", "The job was not found.")
	case errors.Is(err, applicationjobs.ErrUnauthorized):
		response.AbortWithError(c, http.StatusUnauthorized, "unauthorized", "The ownership token is invalid.")
	case errors.Is(err, storage.ErrTooLarge):
		response.AbortWithError(c, http.StatusRequestEntityTooLarge, "file_too_large", "The uploaded file is too large.")
	case errors.Is(err, task.ErrInvalidInput), errors.Is(err, task.ErrInvalidOption):
		response.AbortWithError(c, http.StatusBadRequest, "invalid_request", "The job request is invalid.")
	case errors.Is(err, applicationjobs.ErrConflict):
		response.AbortWithError(c, http.StatusConflict, "job_conflict", "The job cannot perform that action.")
	case errors.Is(err, applicationjobs.ErrNotReady):
		response.AbortWithError(c, http.StatusConflict, "result_not_ready", "The job result is not ready.")
	default:
		response.AbortWithError(c, http.StatusInternalServerError, "internal_error", "The request could not be completed.")
	}
}

func outputContentType(extension string) string {
	switch extension {
	case "jpg", "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "webp":
		return "image/webp"
	case "mp3":
		return "audio/mpeg"
	case "wav":
		return "audio/wav"
	case "m4a":
		return "audio/mp4"
	case "webm":
		return "video/webm"
	case "mov":
		return "video/quicktime"
	default:
		return "video/mp4"
	}
}
