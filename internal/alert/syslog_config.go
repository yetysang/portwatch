package alert

import (
	"fmt"
	"log/syslog"
)

// facilityPriority maps facility name strings to syslog.Priority values.
var facilityPriority = map[string]syslog.Priority{
	"kern":   syslog.LOG_KERN,
	"user":   syslog.LOG_USER,
	"mail":   syslog.LOG_MAIL,
	"daemon": syslog.LOG_DAEMON,
	"auth":   syslog.LOG_AUTH,
	"syslog": syslog.LOG_SYSLOG,
	"lpr":    syslog.LOG_LPR,
	"news":   syslog.LOG_NEWS,
	"local0": syslog.LOG_LOCAL0,
	"local1": syslog.LOG_LOCAL1,
	"local2": syslog.LOG_LOCAL2,
	"local3": syslog.LOG_LOCAL3,
	"local4": syslog.LOG_LOCAL4,
	"local5": syslog.LOG_LOCAL5,
	"local6": syslog.LOG_LOCAL6,
	"local7": syslog.LOG_LOCAL7,
}

// resolveFacility converts a facility name to a syslog.Priority.
func resolveFacility(name string) (syslog.Priority, error) {
	p, ok := facilityPriority[name]
	if !ok {
		return 0, fmt.Errorf("unknown syslog facility: %q", name)
	}
	return p, nil
}
