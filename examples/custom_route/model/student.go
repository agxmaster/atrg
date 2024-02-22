package model

import "github.com/agxmaster/atm/model"

type Student struct {
	Id     int64  `json:"id" column:"id" gorm:"primaryKey"`
	Name   string `json:"name,required" column:"name"`
	Gender int64  `json:"sex" column:"sex"`
	Age    int64  `json:"age" column:"age"`
	Class  string `json:"class" column:"class"`
	model.OperatorTime
}

func (s Student) TableName() string {
	return "student"
}
