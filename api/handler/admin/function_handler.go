package admin

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FunctionHandler struct {
	handler.BaseHandler
	db *gorm.DB
}

func NewFunctionHandler(app *core.AppServer, db *gorm.DB) *FunctionHandler {
	h := FunctionHandler{db: db}
	h.App = app
	return &h
}

func (h *FunctionHandler) Save(c *gin.Context) {
	var data vo.Function
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	var f = model.Function{
		Id:          data.Id,
		Name:        data.Name,
		Label:       data.Label,
		Description: data.Description,
		Parameters:  utils.JsonEncode(data.Parameters),
		Required:    utils.JsonEncode(data.Required),
		Action:      data.Action,
		Token:       data.Token,
		Enabled:     data.Enabled,
	}

	res := h.db.Save(&f)
	if res.Error != nil {
		resp.ERROR(c, "error with save data:"+res.Error.Error())
		return
	}
	data.Id = f.Id
	resp.SUCCESS(c, data)
}

func (h *FunctionHandler) Set(c *gin.Context) {
	var data struct {
		Id    uint        `json:"id"`
		Filed string      `json:"filed"`
		Value interface{} `json:"value"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	res := h.db.Model(&model.Function{}).Where("id = ?", data.Id).Update(data.Filed, data.Value)
	if res.Error != nil {
		resp.ERROR(c, "更新数据库失败！")
		return
	}
	resp.SUCCESS(c)
}

func (h *FunctionHandler) List(c *gin.Context) {
	var items []model.Function
	res := h.db.Find(&items)
	if res.Error != nil {
		resp.ERROR(c, "No data found")
		return
	}

	functions := make([]vo.Function, 0)
	for _, v := range items {
		var f vo.Function
		err := utils.CopyObject(v, &f)
		if err != nil {
			continue
		}
		functions = append(functions, f)
	}
	resp.SUCCESS(c, functions)
}

func (h *FunctionHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)

	if id > 0 {
		res := h.db.Delete(&model.Function{Id: uint(id)})
		if res.Error != nil {
			resp.ERROR(c, "更新数据库失败！")
			return
		}
	}
	resp.SUCCESS(c)
}