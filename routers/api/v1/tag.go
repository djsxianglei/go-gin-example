package v1

import (
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"githubcom/djsxianglei/go-gin-example/models"
	"githubcom/djsxianglei/go-gin-example/pkg/e"
	"githubcom/djsxianglei/go-gin-example/pkg/setting"
	"githubcom/djsxianglei/go-gin-example/pkg/util"
	"net/http"
)

//获取多个文章标签
func GetTags(c *gin.Context) {
	name := c.Query("name")
	maps := make(map[string]interface{})

	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	var state int = -1
	if arg := c.Query("status"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["status"] = state
	}
	code := e.SUCCESS
	data["data"] = models.GetTags(util.GetPage(c), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"data":   data,
		"msg":    e.GetMsg(code),
	})
}

//新增文章标签
func AddTag(c *gin.Context) {
	name := c.PostForm("name")
	state := com.StrTo(c.DefaultPostForm("state", "0")).MustInt()
	createdBy := c.PostForm("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createdBy, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	code := e.INVALID_PARAMS
	msg := ""
	if ! valid.HasErrors() {
		if ! models.ExistTagByName(name) {
			code = e.SUCCESS
			models.AddTag(name, state, createdBy)
		} else {
			code = e.ERROR_EXIST_TAG
		}
	}
	if code == e.INVALID_PARAMS {
		msg = valid.Errors[0].Message
	} else {
		msg = e.GetMsg(code)
	}
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"msg":    msg,
		"data":   make(map[string]string),
	})

}

//修改文章标签
func EditTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	fmt.Println("name=",name)
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	code := e.INVALID_PARAMS
	msg := ""
	if ! valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}

			models.EditTag(id, data)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}
	if code == e.INVALID_PARAMS {
		msg = valid.Errors[0].Message
	} else {
		msg = e.GetMsg(code)
	}
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"msg":    msg,
		"data":   make(map[string]string),
	})
}

//删除文章的标签
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	msg := ""
	if ! valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			models.DeleteTag(id)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}
	if code == e.INVALID_PARAMS {
		msg = valid.Errors[0].Message
	} else {
		msg = e.GetMsg(code)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"msg":    msg,
		"data":   make(map[string]string),
	})
}
