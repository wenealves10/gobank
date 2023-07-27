package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/wenealves10/gobank/db/sqlc"
	"github.com/wenealves10/gobank/utils"
	"golang.org/x/crypto/bcrypt"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
    FullName string `json:"full_name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type createUserResponse struct {
    Username string `json:"username"`
    FullName string `json:"full_name"`
    Email string `json:"email"`
    PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        if bcrypt.ErrPasswordTooLong == err {
            ctx.JSON(http.StatusBadRequest, errorResponse(err))
            return
        }

        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    arg := db.CreateUserParams{
        Username: req.Username,
        FullName: req.FullName,
        Email: req.Email,
        HashedPassword: hashedPassword,
    }

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

    rsp := createUserResponse{
        Username: user.Username,
        FullName: user.FullName,
        Email: user.Email,
        PasswordChangedAt: user.PasswordChangedAt,
        CreatedAt: user.CreatedAt,
    }

	ctx.JSON(http.StatusOK, rsp)
}