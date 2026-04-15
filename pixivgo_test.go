package pixivgo

import "testing"

func TestInt(t *testing.T) {
	p := Int(42)
	if p == nil || *p != 42 {
		t.Errorf("Int(42) = %v, want pointer to 42", p)
	}
}

func TestInt_ZeroValue(t *testing.T) {
	p := Int(0)
	if p == nil || *p != 0 {
		t.Errorf("Int(0) = %v, want non-nil pointer to 0", p)
	}
}

func TestString(t *testing.T) {
	p := String("hello")
	if p == nil || *p != "hello" {
		t.Errorf("String(\"hello\") = %v, want pointer to \"hello\"", p)
	}
}

func TestString_Empty(t *testing.T) {
	p := String("")
	if p == nil || *p != "" {
		t.Errorf("String(\"\") = %v, want non-nil pointer to \"\"", p)
	}
}

func TestBool(t *testing.T) {
	pt := Bool(true)
	if pt == nil || *pt != true {
		t.Errorf("Bool(true) = %v, want pointer to true", pt)
	}

	pf := Bool(false)
	if pf == nil || *pf != false {
		t.Errorf("Bool(false) = %v, want pointer to false", pf)
	}
}
