package fsdump

// FormFiled ...
type FormFiled struct {
	FID   string   `json:"fid"`
	Value []string `json:"c"`
}

// Form ...
type Form struct {
	ID      string      `json:"_id"`
	Created Time        `json:"createdAt"`
	FormID  string      `json:"form_id"`
	Data    []FormFiled `json:"data"`
}
