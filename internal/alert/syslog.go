package alert

import (
	"fmt"
	"log/syslog"
	"strings"

	"github.com/user/portwatch/internal/monitor"
)

// SyslogHandler sends change alerts to the system syslog daemon.
type SyslogHandler struct {
	writer *syslog.Writer
	tag    string
}

// NewSyslogHandler creates a SyslogHandler that writes to syslog under the
// given tag (e.g. "portwatch"). It uses LOG_DAEMON facility by default.
func NewSyslogHandler(tag string) (*SyslogHandler, error) {
	if tag == "" {
		tag = "portwatch"
	}
	w, err := syslog.New(syslog.LOG_DAEMON|syslog.LOG_INFO, tag)
	if err != nil {
		return nil, fmt.Errorf("syslog: open writer: %w", err)
	}
	return &SyslogHandler{writer: w, tag: tag}, nil
}

// Handle writes each change to syslog. Added bindings are logged at INFO
// level; removed bindings are logged at WARNING level.
func (h *SyslogHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}
	var errs []string
	for _, c := range changes {
		msg := formatSyslogMsg(c)
		var err error
		switch c.Kind {
		case monitor.Added:
			err = h.writer.Info(msg)
		case monitor.Removed:
			err = h.writer.Warning(msg)
		default:
			err = h.writer.Info(msg)
		}
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("syslog: write errors: %s", strings.Join(errs, "; "))
	}
	return nil
}

// Drain is a no-op for SyslogHandler (no buffering).
func (h *SyslogHandler) Drain() error {
	return nil
}

// Close releases the underlying syslog connection.
func (h *SyslogHandler) Close() error {
	return h.writer.Close()
}

func formatSyslogMsg(c monitor.Change) string {
	b := c.Binding
	proto := strings.ToUpper(b.Proto)
	host := b.Host
	if b.Hostname != "" {
		host = b.Hostname
	}
	service := b.Service
	if service == "" {
		service = fmt.Sprintf("%d", b.Port)
	}
	pid := ""
	if b.PID > 0 {
		pid = fmt.Sprintf(" pid=%d", b.PID)
	}
	process := ""
	if b.Process != "" {
		process = fmt.Sprintf(" process=%s", b.Process)
	}
	return fmt.Sprintf("[%s] %s %s:%s%s%s", c.Kind, proto, host, service, pid, process)
}
