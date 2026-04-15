package pixivgo

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestErrAuthRequired(t *testing.T) {
	var err error = ErrAuthRequired
	if err == nil {
		t.Fatal("ErrAuthRequired should not be nil")
	}
	if !strings.Contains(err.Error(), "authentication required") {
		t.Errorf("ErrAuthRequired message = %q, want to contain 'authentication required'", err.Error())
	}
}

func TestPixivError_Error_WithStatusCode(t *testing.T) {
	pe := &PixivError{Message: "forbidden", StatusCode: 403}
	msg := pe.Error()
	if !strings.Contains(msg, "forbidden") {
		t.Errorf("Error() = %q, want to contain 'forbidden'", msg)
	}
	if !strings.Contains(msg, "HTTP 403") {
		t.Errorf("Error() = %q, want to contain 'HTTP 403'", msg)
	}
}

func TestPixivError_Error_WithoutStatusCode(t *testing.T) {
	pe := &PixivError{Message: "connection failed"}
	msg := pe.Error()
	if !strings.Contains(msg, "connection failed") {
		t.Errorf("Error() = %q, want to contain 'connection failed'", msg)
	}
	if strings.Contains(msg, "HTTP 0") {
		t.Errorf("Error() = %q, should not contain 'HTTP 0'", msg)
	}
}

func TestPixivError_Unwrap(t *testing.T) {
	pe := &PixivError{Message: "wrapped", Err: io.EOF}
	if !errors.Is(pe, io.EOF) {
		t.Error("errors.Is(pe, io.EOF) should be true")
	}
}

func TestPixivError_Unwrap_Nil(t *testing.T) {
	pe := &PixivError{Message: "no inner"}
	if pe.Unwrap() != nil {
		t.Error("Unwrap() should return nil when Err is nil")
	}
}

func TestPixivError_ErrorsAs(t *testing.T) {
	var err error = &PixivError{
		Message:    "bad request",
		StatusCode: 400,
		Body:       `{"error":"invalid"}`,
	}

	var pe *PixivError
	if !errors.As(err, &pe) {
		t.Fatal("errors.As should match *PixivError")
	}
	if pe.StatusCode != 400 {
		t.Errorf("StatusCode = %d, want 400", pe.StatusCode)
	}
	if pe.Body != `{"error":"invalid"}` {
		t.Errorf("Body = %q, want %q", pe.Body, `{"error":"invalid"}`)
	}
}
