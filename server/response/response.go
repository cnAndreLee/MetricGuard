package response

import (
	"github.com/gin-gonic/gin"
)

func Response(ctx *gin.Context, res ResponseStruct) {
	ctx.JSON(res.HttpStatus, gin.H{"status": res.Code, "msg": res.Msg, "data": res.Data})
}

func ResponseRequestError(ctx *gin.Context, str string) {
	ctx.JSON(400, gin.H{"status": FailCode, "msg": str, "data": nil})
}

func ResponseServerError(ctx *gin.Context, str string) {
	ctx.JSON(500, gin.H{"status": ServerErrorCode, "msg": str, "data": nil})

}

type ResponseStruct struct {
	HttpStatus int    //http状态
	Code       int    //状态码
	Data       gin.H  //数据
	Msg        string //信息
}

const (
	SuccessCode      = 2000
	FailCode         = 4000
	CheckFailCode    = 4022
	ServerErrorCode  = 5000
	TokenExpriedCode = 4010 //token过期
	UnauthorizedCode = 4011 //未授权
	ForbiddenCode    = 4030 //无管理员权限
	VipExiredCode    = 4031 //会员过期

	CategoryNameExistedCode          = 4100
	SecondaryCategoryNameExistedCode = 4100
)

// func Success(ctx *gin.Context, code int, data gin.H, msg string) {
//     Response(ctx, http.StatusOK, 200, data, msg)
// }

// func Fail(ctx *gin.Context, code int, data gin.H, msg string) {
//     Response(ctx, http.StatusOK, 400, data, msg)
// }
