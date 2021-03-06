package controller

import (
	"github.com/ephz3nt/ginEssential/dto"
	"github.com/ephz3nt/ginEssential/response"
	"log"
	"net/http"

	common "github.com/ephz3nt/ginEssential/common"
	model "github.com/ephz3nt/ginEssential/model"
	util "github.com/ephz3nt/ginEssential/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *gin.Context) {
	DB := common.GetDB()
	// 获取参数
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	// 数据验证
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "message": "手机号必须为11位"})
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")

		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "message": "密码不能少于6位"})
		return
	}
	// 如果昵称没有传，给一个10位随机字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	log.Println(name, telephone, password)
	// 判断手机号是否存在
	if isTelephoneExist(DB, telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已存在")

		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "message": "用户已存在"})
		return
	}
	// 创建用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "hash出错")

		//ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "hash出错"})
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hashedPassword),
	}
	DB.Create(&newUser)
	// 返回结果
	//ctx.JSON(200, gin.H{
	//	"code":    200,
	//	"message": "注册成功",
	//})
	response.Success(ctx, nil, "注册成功")
}

func Login(ctx *gin.Context) {
	var DB = common.GetDB()
	// 获取参数
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	// 数据验证
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")

		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "message": "手机号必须为11位"})
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")

		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "message": "密码不能少于6位"})
		return
	}
	// 判断手机号是否存在
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已存在")

		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "message": "用户不存在"})
		return
	}
	// 判断密码是否正确
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 400, nil, "密码错误")

		//ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "密码错误"})
		return
	}
	// 返回token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "token generate failed")

		//ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "token generate failed"})
		log.Printf("token generate error: %v", err)
		return
	}
	//ctx.JSON(200, gin.H{
	//	"code":    200,
	//	"data":    gin.H{"token": token},
	//	"message": "登录成功",
	//})
	response.Success(ctx, gin.H{"token": token}, "登录成功")

}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}})
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
