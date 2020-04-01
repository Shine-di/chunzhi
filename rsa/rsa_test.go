package rsa

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifySign(t *testing.T) {
	sign, errSign := Sign("game_id=651841432110167848&request_time=1573186162&tenant_id=2", private_key)
	errVerify := VerifySign("tenant_id=2&game_id=0&league_id=0&area_id=0", "qg8Vp1i4yN070mM9UFkDd6HYIpPZEfsP8j8qOh7W+tSOhJY/MTVGawQOAS/EHUQa6hM4ODQfNMw9af4OebckOROfSqpkxtM98i0AtiRbaKHy9CACX5LDGFfFqgivgYc0/B2IUNH4lAxdqnx7JkE4oNVGid2kdArzTXorq/P1KmMDqs/VR2xLZDX0n4/k3WRf4IbKHnAvr6oDlAvUgpGBihT56ayUmU1T0HFsFYuiNBZ80Zmb44oseSD3K20dnG8x4fl0xGAPsI1UDaReZZQh4mStLGh8lU5vc+HNFXTZuWgD3GZsP2h/JvjVdUWQN36wfvbKJNymUcuUzYyV9CM3dg==", public_key)

	fmt.Println(sign)
	fmt.Println(errSign)
	fmt.Println(errVerify)

	ass := assert.New(t)
	ass.True(len(sign) > 0 && errVerify == nil)

} //620568336701.dkr.ecr.ap-southeast-1.amazonaws.com/risewinter/data-result-statistics:master-14-25aa0c5d058cb9b7b721f68c6693d8a6606a5da5
//620568336701.dkr.ecr.ap-southeast-1.amazonaws.com/risewinter/data-b-api:master-62-a8cf3d8eec4c211c53800f99464dd89213c15af3
const private_key = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDop2sdYJQ4yCEg
Qu5Lp32Ou82zBTs6ty62qgMy4oF21hyZQq7iKmuZZY4gBsz5GShgncTqpLD0TxiZ
qOYpARyLGibsQPZ3wIJnPPlbtCuFEDZdEhfLo08aGv01TgnhIcoE7OwxR3dV9jeJ
N77qw8ZJzSR2UlAZ0y3QvNxPWqJBoWap5d8UAU2jQrx2d+uj7EVUTVtRsCE4QPeU
z385e0c1VOu8qDYz/j/rDrS2BlsmeDXq9ZFobmvoO+pZ2lHjFTRCeaLCVKnpX8d4
+mZK3bTpXm9Fdu8Nx9jZrC3vykirqUPnnZ7TB+pBvRt3Vz9SEpk1yXD/8lNyBN62
4IJLSXdDAgMBAAECggEAfJJgLUuwMbMe4ZpU49dbyFhQrMFpVGgPMClaKx3S+mFs
0Lc+0sSp9mnFLurVR6+rygfQD199jGLppiUkj+ITeXvYSXoDPl2qtUKVtf+Dqezj
XvQ4H4Zi7XR0Dd2qNoyUEg0V7tD4WePLGsLpi+SlwJCCLISodRt5FaJ6SFccOA0B
ZKz9JX4hlEOVkjKgFIoovM9G9F+rvGheoDdEpyRRXBVI9G5HhRW5cqB0sNK6A9gX
i/Uv3jeno4W4QO5FboR+S0snyoPOMmQwkhYpJAWblIldEUfmuagnmp0a6sXQd2t5
Wfr7A4P8JtS+u507r6Biquikob+XJApidIPKniUiKQKBgQD7xc05fWcbXv5z77aM
4Te+ZOviQhNGuwz0ngArSBcOVUiA4CgfBQrle6KnZyomUmhGAq6TPjU2iQTtDlZx
UY8dUvZ2ENpuRhNmS8i7Rh+hSrNUVCSK1v1n4+asiSfYYyJws6bCgw7pa0YcszWu
GNvef1gctchLxDCs+rHne5stdwKBgQDsj3BJam1Z/J5D1G1a8aQ6MiPdm4wD3DlS
g+jN2hzCyG+nxaKXKZnx/FE3rV4z62nt6HYuItCd1u+GsLRlTuzlhqVuyo/MzlSj
X0r9Cdf/81h1oEVmAS2puCSNY1TqtE7g/r/0wWDegAI+fPVPyvWJRwa9n7q2KBKc
Abge2YRHlQKBgQD1DLnJqdewGU5aM0efaRmjc4DvMFaosihS8nHBrqHaLoGqBgKm
5naLk0Fl5BBvSif5dGTMJXEPil9EB391Peeop/YARjkDuarqFvrh48enahiPDHKg
u83az0PWTIx+nUaJISI/EeZypBmSl464y7M8pP9yui+gJu0lf7+mSXVo0wKBgF+k
itCUAAxO76oa+++2HSEOXqPdnNl+s4piHMEFu3UhVsttQ5R8VGqbCjdJl/nD53sx
7n4uw0vdt9AsJ3OCWpNeQgquST+T+HJpN8dgsH0iZRSBrS1VsqGY+uZTT+To669a
MEAD42dyN/YNzZzqQSW0mswWBYZaY1PB+jA2352VAoGARurIZOma2sXn1bJHmELH
EbTUJFDPpcR+ZHPwgi+dfzJrGssLSN9fJP3h1AZcKtSS6R9NqkMzYSZBRpnoaJHM
1Sl++/mYa+w+4PsULz7c5RRb4KF3bqouwRAlkcnuAwyt4MaFI+SAZNH8rF2EFxto
ZxKO0EJ8PhrqcHHYjUgHaqw=
-----END PRIVATE KEY-----`

const public_key = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA6KdrHWCUOMghIELuS6d9
jrvNswU7OrcutqoDMuKBdtYcmUKu4iprmWWOIAbM+RkoYJ3E6qSw9E8YmajmKQEc
ixom7ED2d8CCZzz5W7QrhRA2XRIXy6NPGhr9NU4J4SHKBOzsMUd3VfY3iTe+6sPG
Sc0kdlJQGdMt0LzcT1qiQaFmqeXfFAFNo0K8dnfro+xFVE1bUbAhOED3lM9/OXtH
NVTrvKg2M/4/6w60tgZbJng16vWRaG5r6DvqWdpR4xU0QnmiwlSp6V/HePpmSt20
6V5vRXbvDcfY2awt78pIq6lD552e0wfqQb0bd1c/UhKZNclw//JTcgTetuCCS0l3
QwIDAQAB
-----END PUBLIC KEY-----`

//620568336701.dkr.ecr.ap-southeast-1.amazonaws.com/risewinter/data-result-statistics:master-14-25aa0c5d058cb9b7b721f68c6693d8a6606a5da5
