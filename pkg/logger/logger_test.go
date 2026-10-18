package logger

import (
	"fmt"
	"strings"
	"testing"
)

type capturePrinter struct {
	message string
}

func (p *capturePrinter) Log(_ Category, _ Level, _ int, format string, args ...interface{}) {
	p.message = fmt.Sprintf(format, args...)
}

func TestWithFieldPercentSignsRemainLiteral(t *testing.T) {
	oldPrinter := printer
	capture := &capturePrinter{}
	printer = capture
	defer func() {
		printer = oldPrinter
	}()

	var assigned string
	fieldValue := "http://example.com/login?bk_token=bkcrypt%gxxx%D"
	Logger{Depth: 1}.Sys().With("url", fieldValue).AssignWhenLogging(&assigned).Info("get response data")

	for _, message := range []string{capture.message, assigned} {
		if !strings.Contains(message, fieldValue) {
			t.Fatalf("expected message to contain literal field value %q, got %q", fieldValue, message)
		}
		if strings.Contains(message, "MISSING") || strings.Contains(message, "%!") {
			t.Fatalf("expected message not to contain fmt missing marker, got %q", message)
		}
	}
}
