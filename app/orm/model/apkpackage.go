package model
type ApkPackage struct {
	Id int64 `json:"id"`
	Guid string `json:"guid"`
	Version string `json:"version"`
	Ts string `json:"ts"`
	Link string `json:"link"`
}


func (ApkPackage) TableName()  string {
	return "bsmi_apk_package"
	
}