package handler

import (
	"context"
	"github.com/Cosydays/codeql_demo/model"
	"github.com/Cosydays/codeql_demo/rpc_sdk"
	"strconv"
)

func ChangeEmail(ctx context.Context, req model.ChangeEmailRequest) {
	email := req.GetEmail()
	userId := req.GetUserId()

	queryUserReq := rpc_sdk.NewQueryUserRequest()
	queryUserReq.SetUserId(strconv.FormatInt(userId, 10))

	//field email from rpc
	userInfo := rpc_sdk.RpcQueryUser(ctx, queryUserReq)

	//field email flow to rpc
	rpc_sdk.RpcChangeEmail(ctx, userInfo.Email, email)
}
