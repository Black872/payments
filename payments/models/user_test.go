package models

import "testing"

func TestSignUpValidation(t *testing.T) {
	fname := "(*User).SignUpValidation()"
	tests := []struct {
		in   User
		want error
	}{
		{
			in: User{
				Name:         "Mike Jordan",
				Email:        "mike@example.com",
				PasswordHash: "a6b7c86e3",
			},
			want: nil,
		},
		{
			in: User{
				Name:         "",
				Email:        "mike@example.com",
				PasswordHash: "a6b7c86e3",
			},
			want: errNameIsEmpty,
		},
		{
			in: User{
				Name:         "Mike Jordan",
				Email:        "",
				PasswordHash: "a6b7c86e3",
			},
			want: errEmailIsEmpty,
		},
		{
			in: User{
				Name:         "Mike Jordan",
				Email:        "mike@example",
				PasswordHash: "a6b7c86e3",
			},
			want: errEmailIsNotValid,
		},
		{
			in: User{
				Name:         "Mike Jordan",
				Email:        "mike@example.com",
				PasswordHash: "",
			},
			want: errPasswordIsEmpty,
		},
	}
	for _, tt := range tests {
		if got := tt.in.SignUpValidation(); got != tt.want {
			t.Errorf("%v: got = %#v; want = %#v\n", fname, got, tt.want)
		}
	}
}

func TestLoginValidation(t *testing.T) {
	fname := "(*User).LoginValidation()"
	tests := []struct {
		in   User
		want error
	}{
		{
			in: User{
				Email:        "mike@example.com",
				PasswordHash: "a6b7c86e3",
			},
			want: nil,
		},
		{
			in: User{
				Email:        "",
				PasswordHash: "a6b7c86e3",
			},
			want: errEmailIsEmpty,
		},
		{
			in: User{
				Email:        "mike@example",
				PasswordHash: "a6b7c86e3",
			},
			want: errEmailIsNotValid,
		},
		{
			in: User{
				Email:        "mike@example.com",
				PasswordHash: "",
			},
			want: errPasswordIsEmpty,
		},
	}
	for _, tt := range tests {
		if got := tt.in.LoginValidation(); got != tt.want {
			t.Errorf("%v: got = %#v; want = %#v\n", fname, got, tt.want)
		}
	}
}

func TestIsEmailValid(t *testing.T) {
	fname := "isEmailValid"
	tests := []struct {
		in   string
		want bool
	}{
		{in: "", want: false},
		{in: "@this.is.a.tag", want: false},
		{in: "this.is.email@example.com", want: true},
		{in: "abc@xy.z", want: false},
		{in: "abc@xy.zz", want: true},
		{in: "th-i5@ex-ample.com", want: true},
		{in: "th$4@example.com.zz", want: false},
		{in: "some+email@example.co.nz", want: true},
	}
	for _, tt := range tests {
		if got := isEmailValid(tt.in); got != tt.want {
			t.Errorf("%v(%#v): got = %#v want = %v\n", fname, tt.in, got, tt.want)
		}
	}
}
