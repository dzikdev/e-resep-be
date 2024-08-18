package helper

import (
	"time"
	_ "time/tzdata"
)

var TimezoneJakarta, _ = time.LoadLocation("Asia/Jakarta")
