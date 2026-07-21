package request

// ProfitQuery 利润查询
type ProfitQuery struct {
	Period string `json:"period" form:"period"`
	Months int    `json:"months" form:"months"`
}
