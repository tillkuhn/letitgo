package ioutil

import (
	"io"

	"github.com/rs/zerolog/log"
)

// SafeClose useful in defer
func SafeClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Err(err).Msg(err.Error())
	}
}
