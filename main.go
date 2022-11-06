package main

import (
	"context"
	"github.com/Cosydays/codeql_demo/biz/handler"
	"github.com/Cosydays/codeql_demo/client"
	"github.com/Cosydays/codeql_demo/model"
)

func main() {
	ctx := context.Background()
	client.InitRedisClient(ctx)

	infoMap := map[string]string{}
	infoMap["old_email"] = "test5@email.com"
	deleteEmailReq := model.DeleteEmailRequest{
		UserId: 123,
		Email:  "test1@email.com",
		Info:   infoMap,
	}
	handler.DeleteEmail(ctx, deleteEmailReq)

	updateEmailReq := model.UpdateEmailRequest{
		UserId: 456,
		Email:  "test2@email.com",
	}
	handler.UpdateEmail(ctx, updateEmailReq)

	createEmailReq := model.CreateEmailRequest{
		UserId: 789,
		Email:  "test3@email.com",
		Status: "aa",
	}
	handler.CreateEmail(ctx, createEmailReq)

	changeEmailReq := model.ChangeEmailRequest{
		UserId: 123,
		Email:  "test4@email.com",
	}
	handler.ChangeEmail(ctx, changeEmailReq)
}
