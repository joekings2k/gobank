package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/joekings2k/gobank/db/sqlc"
)

type createAccountRequest struct{
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func(server *Server) createAccount(ctx *gin.Context){
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err!= nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}
	arg := db.CreateAccountParams{
		Owner: req.Owner,
		Currency: req.Currency,
		Balance: 0,
	}

	account,err  := server.store.CreateAccount(ctx,arg)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK,account)
}

func (server *Server) checkHealth(ctx *gin.Context){
	ctx.JSON(http.StatusOK,gin.H{
		"status":"ok",
		"message":"server is healthy",
	})
}

type getAccountRequest struct{
	ID int64 `uri:"id" binding:"required,min=1"`
}

func(server *Server) getAccount(ctx *gin.Context){
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req);err!=nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}
	account,err :=server.store.GetAccount(ctx,req.ID)
	if err !=nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound,errorResponse(err))
		return
		}
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,account)
}

type ListAccountsRequest struct{
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`

}
func(server *Server) ListAccounts(ctx *gin.Context){
	var req ListAccountsRequest
	if err := ctx.ShouldBindQuery(&req);err!=nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit: req.PageSize,
		Offset: (req.PageID - 1) *req.PageSize,
	}
	accounts,err :=server.store.ListAccounts(ctx,arg)
	if err !=nil{
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,accounts)
}


type updateAccountUri struct{
	ID int64 `uri:"id" binding:"required,min=1"`
	
}

type updateAccountBody struct {
	Balance int64 `json:"balance" binding:"required"`
}
func(server *Server) updateAccount(ctx *gin.Context){
	var uri updateAccountUri
	if err := ctx.ShouldBindUri(&uri);err!=nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}
	var body updateAccountBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
      ctx.JSON(http.StatusBadRequest, errorResponse(err))
      return
  }

	arg := db.UpdateAccountParams{
		ID: uri.ID,
		Balance:body.Balance ,
	}
	accounts,err :=server.store.UpdateAccount(ctx,arg)
	if err !=nil{
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,accounts)
}

type deleteAccountUri struct{
	ID int64 `uri:"id" binding:"required,min=1"`
}

func(server *Server) deleteAccount(ctx *gin.Context){
	var uri deleteAccountUri
	if err := ctx.ShouldBindUri(&uri);err!=nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}

	err :=server.store.DeleteAccount(ctx,uri.ID)
	if err !=nil{
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK,gin.H{
		"status":"ok",
		"message":"Account has been deleted",
	})
}