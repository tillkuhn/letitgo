package format

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
)

const JwtClaimExpiry = "exp"

func HumanizeTokenExpiry(claims map[string]interface{}) (expiryInfo string, err error) {
	expires, ok := claims[JwtClaimExpiry].(float64)
	if !ok {
		return "", fmt.Errorf("exp claim %v for expiry does not exist or cannot be cast to expected float64", claims["exp"])
	} else {
		severity := "👍"
		diff := int64(expires) - time.Now().Unix()
		if diff < 604800 { // less than a week
			severity = "⚠️"
		}
		expiryInfo = fmt.Sprintf("%s %s", humanize.Time(time.Unix(int64(expires), 0)), severity)
	}
	return expiryInfo, nil
}

// FunWithStringSlices based on Slice unicode/ascii strings in golang? https://stackoverflow.com/a/31418646/4292075
func FunWithStringSlices() {
	s1 := "tillkuhn"
	s2 := "tüllkühn!!!"
	fmt.Println(s1[0:8])                 // prints tillkuhn
	fmt.Println(s2[0:8])                 // prints tüllkü
	fmt.Println(string([]rune(s2)[0:8])) // prints tüllkühn

	s3 := "维基百科:关于中文维基百科"
	fmt.Println(string([]rune(s3)[2:9])) // prints 百科:关于中文
}
