package global

var GlobalData = make(map[string]interface{})

func SetRequestId(id string) {
	GlobalData["RequestId"] = id
	return
}

func GetDataByKey(key string) interface{} {
	return GlobalData[key]
}