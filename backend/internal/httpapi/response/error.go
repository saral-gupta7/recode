package response

import (
	"github.com/gin-gonic/gin"

	"github.com/saral-gupta7/recode/backend/internal/observability"
)

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func AbortWithError(
	c *gin.Context,
	status int,
	code string,
	message string,
) {
	requestID, _ := observability.RequestIDFromContext(c.Request.Context())

	c.AbortWithStatusJSON(status,
		errorEnvelope{
			Error: errorBody{
				Code:      code,
				Message:   message,
				RequestID: requestID,
			},
		},
	)
}
