package handler

import (
	"context"
	"fmt"
	"github.com/Cosydays/codeql_demo/model"
	"github.com/Cosydays/codeql_demo/rpc_sdk"
)

func DeleteEmail(ctx context.Context, req model.DeleteEmailRequest) {
	if req.GetUserId() > 0 {
		fmt.Println(req.GetUserId())
	}
	//field email flow to rpc

	//email from req field
	email := req.GetEmail()
	deleteEmailReq := rpc_sdk.NewDeleteEmailRequest()
	deleteEmailReq.RpcEmail = email

	//old_email from req Info map
	oldEmail := req.Info["old_email"]
	deleteEmailReq.OldEmail = oldEmail
	rpc_sdk.RpcDeleteEmailInfo(ctx, deleteEmailReq)
}
