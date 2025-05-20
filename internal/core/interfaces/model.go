package interfaces

import "gorm.io/gorm"

type Entity interface {
	SetID(id uint)
	GetID() uint
	TableName() string
}

type BaseModel struct {
	gorm.Model
}

func (b *BaseModel) GetID() uint {
	return b.ID
}

func (b *BaseModel) SetID(id uint) {
	b.ID = id
}
