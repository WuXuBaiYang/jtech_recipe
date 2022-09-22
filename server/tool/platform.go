package tool

// 平台标签/下标对照表
var platformMap = map[string]int{
	"ios":      1,
	"android":  2,
	"windows":  3,
	"xos":      4,
	"web":      5,
	"mini_web": 6,
	"linux":    7,
}

// Platform2Int 平台信息标签转下标
func Platform2Int(tag string) *int {
	if v, ok := platformMap[tag]; ok {
		return &v
	}
	return nil
}

// Platform2Tag 平台信息下标转标签
func Platform2Tag(i int) *string {
	for k, v := range platformMap {
		if v == i {
			return &k
		}
	}
	return nil
}

// PlatformVerify 验证platform是否存在
func PlatformVerify(v any) bool {
	switch v.(type) {
	case string:
		return Platform2Int(v.(string)) != nil
	case int:
		return Platform2Tag(v.(int)) != nil
	}
	return false
}
