package gvalstrings

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/PaesslerAG/gval"
)

func TestSingleQuoted(t *testing.T) {
	tests := []struct {
		expression string
		want       interface{}
		wantError  string
	}{
		{
			expression: `'Hello World'`,
			want:       `Hello World`,
		},
		{
			expression: `'Hello World' + 10`,
			want:       `Hello World10`,
		},
		{
			expression: `'Hello' World'`,
			wantError:  `unexpected Ident while scanning operator`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			got, err := gval.Evaluate(tt.expression, nil, SingleQuoted())
			if tt.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantError) {
					t.Errorf("got error %v, want error %v", err, tt.wantError)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("gval.Evaluate(%s, SingleQuoted()) = %v, %v, want %v", tt.expression, got, err, tt.want)
			}
		})
	}
}

func ExampleSingleQuoted() {
	value, err := gval.Evaluate(`'Hello `+"World!", nil, SingleQuoted())
	if err != nil {
		panic(err)
	}

	fmt.Println(value)

	// Output
	// Hello World!
}
