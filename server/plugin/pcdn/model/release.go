package model

import "github.com/flipped-aurora/gin-vue-admin/server/global"

// PcdnAgentRelease agent 版本发布（OTA 升级源）
type PcdnAgentRelease struct {
	global.GVA_MODEL
	Version     string `json:"version" gorm:"column:version;type:varchar(32);uniqueIndex;comment:版本号"`
	DownloadURL string `json:"downloadUrl" gorm:"column:download_url;type:varchar(500);comment:下载地址"`
	Checksum    string `json:"checksum" gorm:"column:checksum;type:varchar(128);comment:sha256校验"`
	Stable      bool   `json:"stable" gorm:"column:stable;comment:是否稳定版"`
	Force       bool   `json:"force" gorm:"column:force;comment:是否强制升级"`
	Remark      string `json:"remark" gorm:"column:remark;type:varchar(500);comment:备注"`
	DeptID      uint   `json:"deptId" gorm:"column:dept_id;comment:归属部门(数据权限);index"`
	CreatedBy   uint   `json:"createdBy" gorm:"column:created_by;comment:创建人(数据权限)"`
}

// TableName agent 版本表名
func (PcdnAgentRelease) TableName() string {
	return "gva_pcdn_agent_release"
}
