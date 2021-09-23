package format

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"time"
)

const 	JwtClaimExpiry  = "exp"

func HumanizeTokenExpiry(claims map[string]interface{}) (expiryInfo string, err error) {
	expires, ok := claims[JwtClaimExpiry].(float64)
	if !ok {
		return "",fmt.Errorf("exp claim %v for expiry does not exist or cannot be cast to expected float64", claims["exp"])
	} else {
		severity := "👍"
		diff := int64(expires) - time.Now().Unix()
		if diff < 604800 { // less than a week
			severity = "⚠️"
		}
		expiryInfo = fmt.Sprintf("%s %s", humanize.Time(time.Unix(int64(expires), 0)), severity)
	}
	return expiryInfo,nil
}
