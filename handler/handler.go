package handler

import (
	"github.com/agxmaster/atrg/core"
	"github.com/agxmaster/atrg/process"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Atr struct {
}

func (a *Atr) Info(c *gin.Context) {

	id, err := strconv.Atoi(c.Param(core.IdKey))
	if err != nil {
		core.Error(c, err)
		return
	}

	res, err := process.ProcessFactory(c.Request.Context(), c, c.Param(core.ModelKey)).Info(id)
	if err != nil {
		core.Error(c, err)
		return
	}

	core.Success(c, res)
}

func (a *Atr) List(c *gin.Context) {

	var params = make(map[string]interface{})

	for k, v := range c.Request.URL.Query() {
		params[k] = v
	}

	page, err := strconv.Atoi(c.Param(core.PageKey))

	if err != nil {
		core.Error(c, err)
		return
	}
	res, err := process.ProcessFactory(c.Request.Context(), c, c.Param(core.ModelKey)).List(params, page)

	if err != nil {
		core.Error(c, err)
		return
	}
	core.Success(c, res)

}

func (a *Atr) Delete(c *gin.Context) {

	id, err := strconv.Atoi(c.Param(core.IdKey))
	if err != nil {
		core.Error(c, err)
		return
	}

	err = process.ProcessFactory(c.Request.Context(), c, c.Param(core.ModelKey)).Delete(id)
	if err != nil {
		core.Error(c, err)
		return
	}
	core.Success(c, nil)
}

func (a *Atr) Create(c *gin.Context) {
	err := process.ProcessFactory(c.Request.Context(), c, c.Param(core.ModelKey)).Create()

	if err != nil {
		core.Error(c, err)
		return
	}
	core.Success(c, nil)

}

func (a *Atr) BatchCreate(c *gin.Context) {
	err := process.ProcessFactory(c.Request.Context(), c, c.Param(core.ModelKey)).BatchCreate()
	if err != nil {
		core.Error(c, err)
		return
	}
	core.Success(c, nil)

}

func (a *Atr) Update(c *gin.Context) {

	id, err := strconv.Atoi(c.Param(core.IdKey))
	if err != nil {
		core.Error(c, err)
		return
	}

	err = process.ProcessFactory(c.Request.Context(), c, c.Param(core.ModelKey)).Update(id)
	if err != nil {
		core.Error(c, err)
		return
	}
	core.Success(c, nil)

}
