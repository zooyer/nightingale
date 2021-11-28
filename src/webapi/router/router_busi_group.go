package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
	"github.com/toolkits/pkg/str"

	"github.com/didi/nightingale/v5/src/models"
)

type busiGroupForm struct {
	Name    string                   `json:"name" binding:"required"`
	Members []models.BusiGroupMember `json:"members"`
}

func busiGroupAdd(c *gin.Context) {
	var f busiGroupForm
	ginx.BindJSON(c, &f)

	if len(f.Members) == 0 {
		ginx.Bomb(http.StatusBadRequest, "members empty")
	}

	username := c.MustGet("username").(string)
	ginx.NewRender(c).Message(models.BusiGroupAdd(f.Name, f.Members, username))
}

func busiGroupPut(c *gin.Context) {
	var f busiGroupForm
	ginx.BindJSON(c, &f)

	username := c.MustGet("username").(string)
	targetbg := c.MustGet("busi_group").(*models.BusiGroup)
	ginx.NewRender(c).Message(targetbg.Update(f.Name, username))
}

func busiGroupMemberAdd(c *gin.Context) {
	var members []models.BusiGroupMember
	ginx.BindJSON(c, &members)

	username := c.MustGet("username").(string)
	targetbg := c.MustGet("busi_group").(*models.BusiGroup)

	ginx.NewRender(c).Message(targetbg.AddMembers(members, username))
}

func busiGroupMemberDel(c *gin.Context) {
	var members []models.BusiGroupMember
	ginx.BindJSON(c, &members)

	username := c.MustGet("username").(string)
	targetbg := c.MustGet("busi_group").(*models.BusiGroup)

	ginx.NewRender(c).Message(targetbg.DelMembers(members, username))
}

func busiGroupDel(c *gin.Context) {
	username := c.MustGet("username").(string)
	targetbg := c.MustGet("busi_group").(*models.BusiGroup)

	err := targetbg.Del()
	if err != nil {
		logger.Infof("busi_group_delete fail: operator=%s, group_name=%s error=%v", username, targetbg.Name, err)
	} else {
		logger.Infof("busi_group_delete succ: operator=%s, group_name=%s", username, targetbg.Name)
	}

	ginx.NewRender(c).Message(err)
}

// 我是超管、或者我是业务组成员
func busiGroupGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", defaultLimit)
	query := ginx.QueryStr(c, "query", "")

	me := c.MustGet("user").(*models.User)
	lst, err := me.BusiGroups(limit, query)

	ginx.NewRender(c).Data(lst, err)
}

// 这个接口只有在活跃告警页面才调用，获取各个BG的活跃告警数量
func busiGroupAlertingsGets(c *gin.Context) {
	ids := ginx.QueryStr(c, "ids", "")
	ret, err := models.AlertNumbers(str.IdsInt64(ids))
	ginx.NewRender(c).Data(ret, err)
}

func busiGroupGet(c *gin.Context) {
	bg := BusiGroup(ginx.UrlParamInt64(c, "id"))
	ginx.Dangerous(bg.FillUserGroups())
	ginx.NewRender(c).Data(bg, nil)
}