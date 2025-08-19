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
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	// Create decoder with strict validation
	decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
	decoder.DisallowUnknownFields()

	var req createUserRequest
	if err := decoder.Decode(&req); err != nil {
		// Check if it's an unknown field error
		if strings.Contains(err.Error(), "unknown field") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid parameter found",
				"details": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Validate the bound data using validator
	if err := validator.New().Struct(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: req.Password,
		},
		AfterCreate: func(user db.User) error {
			taskPayload := worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}
			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, &taskPayload)
			//if err != nil {
			//	ctx.JSON(http.StatusInternalServerError, gin.H{
			//		"error":   "Failed to send verification email",
			//		"details": err.Error(),
			//	})
			//	return
			//}
		},
	}

	txResult, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
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
