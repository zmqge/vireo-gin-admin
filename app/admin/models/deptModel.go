package models

type Dept struct {
	ID       uint   `gorm:"primaryKey;comment:主键" json:"id"`
	ParentID uint   `gorm:"column:parent_id;comment:父部门ID" json:"parentId"`
	Name     string `gorm:"size:50;comment:部门名称" json:"name"`
	Code     string `gorm:"size:50;comment:部门编码" json:"code"`
	Sort     int    `gorm:"comment:排序" json:"sort"`
	Status   int    `gorm:"comment:状态" json:"status"`
}

type DeptV0 struct {
	ID       uint      `json:"id"`
	ParentID uint      `json:"parentId"`
	Name     string    `json:"name"`
	Code     string    `json:"code"`
	Sort     int       `json:"sort"`
	Status   int       `json:"status"`
	Children *[]DeptV0 `json:"children,omitempty"` // 使用指针和omitempty
}
