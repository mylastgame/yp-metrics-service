package counter

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFromString(t *testing.T) {
	type args struct {
		title string
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    *Counter
		wantErr bool
	}{
		{
			name:    "Int to counter",
			args:    args{"cnt1", "100"},
			want:    &Counter{"cnt1", ValType(100)},
			wantErr: false,
		},
		{
			name:    "Float to counter",
			args:    args{"cnt1", "100.0001"},
			want:    &Counter{},
			wantErr: true,
		},
		{
			name:    "Symbol to counter",
			args:    args{"cnt1", "100a"},
			want:    &Counter{},
			wantErr: true,
		},
		{
			name:    "Empty string to counter",
			args:    args{"cnt1", ""},
			want:    &Counter{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromString(tt.args.title, tt.args.value)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		title string
		val   ValType
	}
	tests := []struct {
		name string
		args args
		want *Counter
	}{
		{
			name: "New cnt1",
			args: args{"cnt1", ValType(100)},
			want: &Counter{"cnt1", ValType(100)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.title, tt.args.val)
			require.New(t)
			assert.Equal(t, tt.want, got)
		})
	}
}
