/**
 * @author: D-S
 * @date: 2020/5/11 4:32 下午
 */

package sign

import "game-test/constant"

func GetKey(string2 string) string {
	switch string2 {
	case "1":
		return constant.Private_key
	case "18":
		return constant.Private_key18
	case "19":
		return constant.Private_key19
	default:
		return constant.Private_key
	}
}
