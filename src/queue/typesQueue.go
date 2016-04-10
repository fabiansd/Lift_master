package queue

import(
	"time"
)


type orderStatus struct {
	active bool
	addr   string      `json:"-"`
	timer  *time.Timer `json:"-"`
}

var inactive = orderStatus{active: false, addr: "", timer: nil}
