package speller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpeller(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "Correct correcting",
			text: "Масквa эттто",
			want: "Москва это",
		},
		{
			name: "With no text",
			text: "",
			want: "",
		},
		{
			name: "With url",
			text: "https://github.com/Meraiku/kode",
			want: "https://github.com/Meraiku/kode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			corrected, err := CheckText(tt.text)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(tt.want, corrected)
		})
	}
}
