package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPlainText(t *testing.T) {
	tableTests := []struct {
		name   string
		arg    []byte
		result string
	}{
		{
			name:   "test1",
			arg:    []byte("http://htgfnn.yandex/nubcnadqasd321"),
			result: "ee136682",
		},
		{
			name:   "test2",
			arg:    []byte("4c2432bf"),
			result: "e3d70a5b",
		},
	}

	for _, test := range tableTests {
		t.Run(test.name, func(t *testing.T) {
			hashResult, _ := HashPlainText(test.arg)
			assert.Equal(t, test.result, hashResult)

		})

	}

}

func TestDeHashText(t *testing.T) {

	tests := []struct {
		name string
		args string
		want []byte
	}{
		{
			name: "test1",
			args: "ee136682",
			want: []byte("http://htgfnn.yandex/nubcnadqasd321"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _ = HashPlainText(tt.want)
			got, _ := DeHashText(tt.args)

			assert.Equalf(t, tt.want, got, "DeHashText(%v)", tt.args)
		})
	}
}
