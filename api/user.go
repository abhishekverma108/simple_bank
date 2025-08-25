package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/worker"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/v2"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type createUserResponse struct {
	ID        int32  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	// Span: read body
	spanRead, _ := apm.StartSpan(reqCtx, "read_request_body", "io")
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	spanRead.End()
	if err != nil {
		apm.CaptureError(reqCtx, err).Send()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	// Create decoder with strict validation
	spanDecode, _ := apm.StartSpan(reqCtx, "decode_and_validate_user", "validation")
	decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
	decoder.DisallowUnknownFields()

	var req createUserRequest
	if err := decoder.Decode(&req); err != nil {
		spanDecode.End()
		// Check if it's an unknown field error
		if strings.Contains(err.Error(), "unknown field") {
			apm.CaptureError(reqCtx, err).Send()
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid parameter found",
				"details": err.Error(),
			})
			return
		}
		apm.CaptureError(reqCtx, err).Send()
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Validate the bound data using validator
	if err := validator.New().Struct(&req); err != nil {
		spanDecode.End()
		apm.CaptureError(reqCtx, err).Send()
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	spanDecode.End()

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: req.Password,
		},
		AfterCreate: func(user db.User) error {
			// Span for distributing task
			spanTask, _ := apm.StartSpan(reqCtx, "distribute_send_verify_email", "taskqueue")
			taskPayload := worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}
			err := server.taskDistributor.DistributeTaskSendVerifyEmail(reqCtx, &taskPayload)
			if err != nil {
				apm.CaptureError(reqCtx, err).Send()
			}
			spanTask.End()
			return err
		},
	}

	spanDB, _ := apm.StartSpan(reqCtx, "create_user_tx", "db")
	txResult, err := server.store.CreateUserTx(reqCtx, arg)
	spanDB.End()
	if err != nil {
		apm.CaptureError(reqCtx, err).Send()
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	spanRedis, _ := apm.StartSpan(reqCtx, "redis_set_user_marker", "cache")
	err = server.redisClient.Set(reqCtx, txResult.User.Username, "first of few users", 0).Err()
	spanRedis.End()
	if err != nil {
		apm.CaptureError(reqCtx, err).Send()
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, createUserResponse{
		ID:        txResult.User.ID,
		Username:  txResult.User.Username,
		Email:     txResult.User.Email,
		CreatedAt: txResult.User.CreatedAt.Format(time.RFC3339),
	})

}

//var req createUserRequest
//// Use json.Decoder with DisallowUnknownFields
//decoder := json.NewDecoder(ctx.Request.Body)
//decoder.DisallowUnknownFields()
//if err := decoder.Decode(&req); err != nil {
//	// Check if it's an unknown field error
//	if strings.Contains(err.Error(), "unknown field") {
//		ctx.JSON(http.StatusBadRequest, gin.H{
//			"error":   "Invalid parameter found",
//			"details": err.Error(),
//		})
//		return
//	}
//	ctx.JSON(http.StatusBadRequest, errorResponse(err))
//	return
//}
//if err := ctx.ShouldBind(&req); err != nil {
//	ctx.JSON(http.StatusBadRequest, errorResponse(err))
//	return
//}

//type getUserRequest struct {
//	ID int32 `uri:"id" binding:"required,min=1"`
//}
//type getUserResponse struct {
//	ID        int32  `json:"id"`
//	Username  string `json:"username"`
//	Email     string `json:"email"`
//	CreatedAt string `json:"created_at"`
//}

//
//func (server *Server) getUser(ctx *gin.Context) {
//	var req getUserRequest
//	if err := ctx.ShouldBindUri(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//
//	user, err := server.store.GetUserByID(ctx, req.ID)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			ctx.JSON(http.StatusNotFound, errorResponse(err))
//			return
//		}
//
//		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
//		return
//	}
//	ctx.JSON(http.StatusOK, getUserResponse{
//		ID:        user.ID,
//		Username:  user.Username,
//		Email:     user.Email,
//		CreatedAt: user.CreatedAt.Format(time.RFC3339),
//	})
//}
//type updateUserRequest struct {
//	ID       int32  `json:"id" binding:"required,min=1"`,
//	Username string `json:"username" binding:"required"`
//	Email    string `json:"email" binding:"required,email"`
//	PasswordHash string `json:"password_hash" binding:"required"`
//}
//type updateUserResponse struct {
//	ID        int32  `json:"id"`
//	Username  string `json:"username"`
//	Email     string `json:"email"`
//	CreatedAt string `json:"created_at"`
//}
//func (server *Server) updateUser(ctx *gin.Context) {
//	var req updateUserRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//
//	arg := db.UpdateUserParams{
//		ID:       req.ID,
//		Username: req.Username,
//		Email:    req.Email,
//		PasswordHash: req.PasswordHash,
//	}
//	user, err := server.store.UpdateUser(ctx, arg)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
//		return
//	}
//	ctx.JSON(http.StatusOK, updateUserResponse{
//		ID:        user.ID,
//		Username:  user.Username,
//		Email:     user.Email,
//		CreatedAt: user.CreatedAt.Format(time.RFC3339),
//	})
//}
