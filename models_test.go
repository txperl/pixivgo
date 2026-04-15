package pixivgo

import (
	"encoding/json"
	"testing"
)

func TestFlexInt_UnmarshalJSON_Number(t *testing.T) {
	var fi FlexInt
	if err := json.Unmarshal([]byte("12345"), &fi); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fi.Int() != 12345 {
		t.Errorf("got %d, want 12345", fi.Int())
	}
}

func TestFlexInt_UnmarshalJSON_String(t *testing.T) {
	var fi FlexInt
	if err := json.Unmarshal([]byte(`"12345"`), &fi); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fi.Int() != 12345 {
		t.Errorf("got %d, want 12345", fi.Int())
	}
}

func TestFlexInt_UnmarshalJSON_Negative(t *testing.T) {
	var fi FlexInt
	if err := json.Unmarshal([]byte("-42"), &fi); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fi.Int() != -42 {
		t.Errorf("got %d, want -42", fi.Int())
	}
}

func TestFlexInt_UnmarshalJSON_Zero(t *testing.T) {
	var fi FlexInt
	if err := json.Unmarshal([]byte("0"), &fi); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fi.Int() != 0 {
		t.Errorf("got %d, want 0", fi.Int())
	}
}

func TestFlexInt_UnmarshalJSON_Error(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"invalid string", `"abc"`},
		{"null", "null"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fi FlexInt
			if err := fi.UnmarshalJSON([]byte(tt.input)); err == nil {
				t.Errorf("expected error for input %q, got nil", tt.input)
			}
		})
	}
}

func TestFlexInt_MarshalJSON(t *testing.T) {
	fi := FlexInt(12345)
	data, err := json.Marshal(fi)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "12345" {
		t.Errorf("got %s, want 12345", string(data))
	}
}

func TestFlexInt_Int(t *testing.T) {
	fi := FlexInt(42)
	if fi.Int() != 42 {
		t.Errorf("got %d, want 42", fi.Int())
	}
}

func TestFlexInt_InStruct(t *testing.T) {
	type s struct {
		ID FlexInt `json:"id"`
	}

	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"number", `{"id": 123}`, 123},
		{"string", `{"id": "123"}`, 123},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v s
			if err := json.Unmarshal([]byte(tt.input), &v); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if v.ID.Int() != tt.want {
				t.Errorf("got %d, want %d", v.ID.Int(), tt.want)
			}
		})
	}
}
