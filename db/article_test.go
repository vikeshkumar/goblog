package db

import "testing"

func Test_createUrlTitle(t *testing.T) {

	tests := []struct {
		name string
		title string
		want string
	}{
		{"simple word", "hello", "hello"},
		{"2 words", "Hello World", "hello-world"},
		{"words with special characters", "How to@    do it without any external libraries?!!", "how-to-do-it-without-any-external-libraries"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createUrlFromTitle(tt.title); got != tt.want {
				t.Errorf("createUrlFromTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}
