package model
type ApkPackage struct {
	Id int64
	GappId string
	Name string
	Version string
	PublishTime string
	Views int64
}


func (ApkPackage) TableName()  string {
	return "bsmi_apk_package"
	
}