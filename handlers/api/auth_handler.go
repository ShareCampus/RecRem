package api

import (
	"net/http"
	"recrem/config/setting"
	"recrem/forms"
	"recrem/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
}

func (a *AuthHandler) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "hello world")
	return
}

func (a *AuthHandler) Register(ctx *gin.Context) {
	regForm := forms.RegisterForm{}
	result := utils.Result{
		Code: utils.Success,
		Msg:  "注册成功",
		Data: nil,
	}
	if err := ctx.ShouldBindJSON(&regForm); err != nil { // 表单校验失败
		result.Code = utils.RequestError     // 请求数据有误
		result.Msg = utils.GetFormError(err) // 获取表单错误信息
		ctx.JSON(http.StatusOK, result)
		return
	}
	user := regForm.BindToModel() // 绑定表单数据到用户
	u, _ := user.GetByUsername()  // 根据用户名获取用户
	if u.Username != "" {         // 账号已经被注册
		result.Code = utils.RequestError
		result.Msg = "账号已经被注册"
		ctx.JSON(http.StatusOK, result)
		return
	}
	if err := user.Create(); err != nil { // 创建用户 + 异常处理
		result.Code = utils.ServerError
		result.Msg = "服务器端错误"
		ctx.JSON(http.StatusOK, result) // 返回 json
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (a *AuthHandler) Login(ctx *gin.Context) {
	loginForm := forms.LoginForm{}
	result := utils.Result{ // 定义 api 返回信息结构
		Code: utils.Success,
		Msg:  "登录成功",
		Data: nil,
	}
	if err := ctx.ShouldBindJSON(&loginForm); err != nil { // 表单校验失败
		result.Code = utils.RequestError     // 请求数据有误
		result.Msg = utils.GetFormError(err) // 获取表单错误信息
		ctx.JSON(http.StatusOK, result)      // 返回 json
		return
	}
	captchaConfig := &utils.CaptchaConfig{
		Id:          loginForm.CaptchaId,
		VerifyValue: loginForm.CaptchaVal,
	}
	if !utils.CaptchaVerify(captchaConfig) { // 校验验证码
		result.Code = utils.RequestError // 请求数据有误
		result.Msg = "验证码错误"
		ctx.JSON(http.StatusOK, result) // 返回 json
		return
	}
	user := loginForm.BindToModel() // 绑定表单数据到实体类
	u, _ := user.GetByUsername()    // 根据用户名获取用户
	if u.Username == "" {           // 用户不存在
		result.Code = utils.RequestError
		result.Msg = "不存在该用户"
		ctx.JSON(http.StatusOK, result)
		return
	}
	if !utils.VerifyPwd(u.Pwd, user.Pwd) { // 密码错误
		result.Code = utils.RequestError
		result.Msg = "密码错误"
		ctx.JSON(http.StatusOK, result) // 返回 json
		return
	}

	j := utils.NewJWT()                             // 创建一个jwt
	token, err := j.CreateToken(utils.CustomClaims{ // 生成 JWT token
		Username: u.Username,
		UserImg:  u.UserImg,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.
				Duration(setting.Config.Server.TokenExpireTime)).Unix(), // 设置过期时间
			IssuedAt: time.Now().Unix(),
		},
	})
	if err != nil { // 异常处理
		result.Code = utils.ServerError
		result.Msg = "服务器端错误"
		ctx.JSON(http.StatusOK, result) // 返回 json
		return
	}
	result.Data = utils.Token{ // 封装 Token 信息
		Token:    token,
		UserId:   u.ID,
		Username: u.Username,
		UserImg:  u.UserImg,
	}
	ctx.JSON(http.StatusOK, result)
}
