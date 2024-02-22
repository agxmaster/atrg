package process

import (
	"github.com/agxmaster/atm"
	"github.com/agxmaster/atrg/core"
	"reflect"
)

type AtmWithModel struct {
	BaseStore
}

func (p *AtmWithModel) Info(id int) (interface{}, error) {
	model := reflect.New(p.ModelConf.Model).Interface()
	data, err := atm.First(p.Ctx, core.Db, model, int64(id))

	if err != nil {
		return data, nil
	}
	res, err := p.ModelConf.FormatModelLine(p.Ctx, data)
	return res, err
}

func (p *AtmWithModel) List(params map[string]interface{}, page int) (interface{}, error) {
	claus, err := p.BaseStore.ProcessList(params, page)
	if err != nil {
		return nil, err
	}
	data, err := atm.QueryPage(p.Ctx, core.Db, p.ModelConf.Model, claus, false)

	if err != nil {
		return nil, err
	}
	data.Data, err = p.ModelConf.FormatModelList(p.Ctx, data.Data)
	return data, err
}

func (p *AtmWithModel) Delete(id int) error {
	return atm.DeleteBatch(p.Ctx, core.Db, reflect.New(p.ModelConf.Model).Interface(), []int64{int64(id)})
}

func (p *AtmWithModel) Create() error {
	data := reflect.New(p.ModelConf.Model).Elem().Interface()
	err := p.C.Bind(&data)

	if err != nil {
		return err
	}

	if p.ModelConf.CreateParamsHandler != nil {
		data, err = p.ModelConf.CreateParamsHandler(p.Ctx, data)
		if err != nil {
			return err
		}
	}
	return atm.Create(p.Ctx, core.Db, reflect.New(p.ModelConf.Model).Interface(), data)
}

func (p *AtmWithModel) BatchCreate() error {
	data := reflect.New(reflect.SliceOf(p.ModelConf.Model)).Elem().Interface()
	err := p.C.Bind(&data)

	if err != nil {
		return err
	}

	if p.ModelConf.CreateBatchParamsHandler != nil {
		data, err = p.ModelConf.CreateBatchParamsHandler(p.Ctx, data)
		if err != nil {
			return err
		}
	}
	return atm.BatchCreate(p.Ctx, core.Db, reflect.New(p.ModelConf.Model).Interface(), data)
}

func (p *AtmWithModel) Update(id int) error {
	data := reflect.New(p.ModelConf.Model).Elem().Interface()

	err := p.C.Bind(&data)
	if err != nil {
		return err
	}

	if p.ModelConf.UpdateParamsHandler != nil {
		data, err = p.ModelConf.UpdateParamsHandler(p.Ctx, data)
		if err != nil {
			return err
		}
	}

	return atm.Update(p.Ctx, core.Db, reflect.New(p.ModelConf.Model).Interface(), int64(id), data)
}
