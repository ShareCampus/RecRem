package api

import (
	"net/http"
	"recrem/config/setting"
	"recrem/forms"
	"recrem/log"
	"recrem/utils"
	"time"

	"recrem/models"

	"github.com/gin-gonic/gin"
	"github.com/go-gomail/gomail"
)

type AuthHandler struct {
}

func (a *AuthHandler) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "hello world")
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
		log.Logger.Sugar().Error("error: ", err.Error())
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
	// captchaConfig := &utils.CaptchaConfig{
	// 	Id:          loginForm.CaptchaId,
	// 	VerifyValue: loginForm.CaptchaVal,
	// }
	// if !utils.CaptchaVerify(captchaConfig) { // 校验验证码
	// 	result.Code = utils.RequestError // 请求数据有误
	// 	result.Msg = "验证码错误"
	// 	ctx.JSON(http.StatusOK, result) // 返回 json
	// 	return
	// }
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

	// j := utils.NewJWT()                             // 创建一个jwt
	// token, err := j.CreateToken(utils.CustomClaims{ // 生成 JWT token
	// 	Username: u.Username,
	// 	UserImg:  u.UserImg,
	// 	StandardClaims: jwt.StandardClaims{
	// 		ExpiresAt: time.Now().Add(time.Second * time.
	// 			Duration(setting.Config.Server.TokenExpireTime)).Unix(), // 设置过期时间
	// 		IssuedAt: time.Now().Unix(),
	// 	},
	// })

	// if err != nil { // 异常处理
	// 	log.Logger.Sugar().Error("error: ", err.Error())
	// 	result.Code = utils.ServerError
	// 	result.Msg = "服务器端错误"
	// 	ctx.JSON(http.StatusOK, result) // 返回 json
	// 	return
	// }

	// result.Data = utils.Token{ // 封装 Token 信息
	// 	Token:    token,
	// 	UserId:   u.ID,
	// 	Username: u.Username,
	// 	UserImg:  u.UserImg,
	// }

	ctx.JSON(http.StatusOK, result)
}

func (a *AuthHandler) CreateCaptcha(ctx *gin.Context) {
	captcha := utils.CaptchaConfig{} // 创建验证码配置结构
	result := utils.Result{          // 返回数据结构
		Code: utils.Success,
		Msg:  "验证码创建成功",
		Data: nil,
	}

	base64, err := utils.GenerateCaptcha(&captcha) // 创建验证码
	if err != nil {                                // 异常处理
		result.Code = utils.ServerError
		result.Msg = "服务器端错误"
		ctx.JSON(http.StatusOK, result)
		return
	}

	result.Data = gin.H{ // 封装 data
		"captcha_id":  captcha.Id,
		"captcha_url": base64,
	}

	ctx.JSON(http.StatusOK, result) // 返回 json 数据
}

func (a *AuthHandler) ForgetPwd(ctx *gin.Context) {
	forgetPwdForm := forms.ForgetPwdForm{}
	if err := ctx.ShouldBindJSON(&forgetPwdForm); err != nil {
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.RequestError,
			Msg:  utils.GetFormError(err),
			Data: nil,
		})
		return
	}

	user, _ := models.User{Email: forgetPwdForm.Email}.GetByEmail()
	if user.Username == "" {
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.RequestError,
			Msg:  "不存在该邮箱帐号",
			Data: nil,
		})
		return
	}

	code := ""
	_ = setting.Cache.Get(forgetPwdForm.Email, &code)
	if code == "" {
		verifyCode, err := utils.CreateRandomCode(6)
		if err != nil {
			log.Logger.Sugar().Error("创建验证码失败：", err.Error())
			ctx.JSON(http.StatusOK, utils.Result{
				Code: utils.ServerError,
				Msg:  "创建验证码失败",
				Data: nil,
			})
			return
		}
		code = verifyCode
		_ = setting.Cache.Set(forgetPwdForm.Email, code, time.Minute*15)
	}

	msg := gomail.NewMessage()
	// 设置收件人
	msg.SetHeader("To", forgetPwdForm.Email)
	// 设置发件人
	msg.SetAddressHeader("From", setting.Config.SMTP.Account, setting.Config.SMTP.Account)
	// 主题
	msg.SetHeader("Subject", "忘记密码验证")
	// log.Logger.Sugar().Info("verifyCode: ", code)
	// 正文
	msg.SetBody("text/html", utils.GetForgetPwdEmailHTML(user.Username, code))
	// 设置 SMTP 参数
	d := gomail.NewDialer(setting.Config.SMTP.Address, setting.Config.SMTP.Port,
		setting.Config.SMTP.Account, setting.Config.SMTP.Password)

	// 发送
	err := d.DialAndSend(msg)
	if err != nil {
		log.Logger.Sugar().Error("验证码发送失败：", err.Error())
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.ServerError,
			Msg:  "验证码发送失败，请检查 smtp 配置",
			Data: nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.Result{
		Code: utils.Success,
		Msg:  "验证码发送成功，请前往邮箱查看",
		Data: nil,
	})
}

func (a *AuthHandler) ResetPwd(ctx *gin.Context) {
	resetPwdForm := forms.ResetPwdForm{}
	if err := ctx.ShouldBindJSON(&resetPwdForm); err != nil {
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.RequestError,
			Msg:  utils.GetFormError(err),
			Data: nil,
		})
		return
	}

	verifyCode := ""
	_ = setting.Cache.Get(resetPwdForm.Email, &verifyCode)
	if verifyCode != resetPwdForm.VerifyCode {
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.RequestError,
			Msg:  "验证码无效或错误",
			Data: nil,
		})
		return
	}

	user := resetPwdForm.BindToModel()
	err := user.UpdatePwd()
	if err != nil {
		log.Logger.Sugar().Error("error: ", err.Error())
		ctx.JSON(http.StatusOK, utils.Result{
			Code: utils.ServerError,
			Msg:  "服务器端错误",
			Data: nil,
		})
		return
	}

	// 删除缓存中的验证码
	_ = setting.Cache.Delete(resetPwdForm.Email)

	ctx.JSON(http.StatusOK, utils.Result{
		Code: utils.Success,
		Msg:  "重置密码成功",
		Data: nil,
	})
}
