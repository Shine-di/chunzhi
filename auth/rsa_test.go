/**
 * @author: D-S
 * @date: 2020/5/12 2:34 下午
 */

package auth

import "testing"

func TestGenRsaKey(t *testing.T) {
	GetRsaKey(2048, "996")
}
