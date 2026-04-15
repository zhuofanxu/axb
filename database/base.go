package database

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/zhuofanxu/axb/customtype"
	"github.com/zhuofanxu/axb/idgen"
)

// BaseModel 定义所有业务模型共享的基础字段。
type BaseModel struct {
	ID        customtype.SnowflakeID `gorm:"primarykey;autoIncrement:false" json:"id"`
	CreatedBy customtype.SnowflakeID `gorm:"type:bigint;default:0;index;comment:创建人ID" json:"createdBy"`
	CreatedAt customtype.JsonTime    `json:"createdAt"`
	UpdatedAt customtype.JsonTime    `json:"updatedAt"`
}

// BeforeCreate 在创建前补齐主键。
func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == 0 {
		m.ID = customtype.SnowflakeID(idgen.GenSnowflakeId().Int64())
	}
	return nil
}

// SetTimezone 设置进程默认时区，应在应用启动时显式调用，例如：database.SetTimezone("Asia/Shanghai")。
// 作为库不应在 init() 中修改全局时区状态。
func SetTimezone(tz string) error {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return fmt.Errorf("invalid timezone %q: %w", tz, err)
	}
	time.Local = loc
	return nil
}
