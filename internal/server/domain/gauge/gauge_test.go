package gauge

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToGauge(t *testing.T) {
	type args struct {
		title string
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    *Gauge
		wantErr bool
	}{
		{
			name:    "Int to gauge",
			args:    args{"cnt1", "100"},
			want:    &Gauge{"cnt1", ValType(100)},
			wantErr: false,
		},
		{
			name:    "Float to gauge",
			args:    args{"cnt1", "100.0001"},
			want:    &Gauge{"cnt1", ValType(100.0001)},
			wantErr: false,
		},
		{
			name:    "Symbol to gauge",
			args:    args{"cnt1", "100a"},
			want:    &Gauge{},
			wantErr: true,
		},
		{
			name:    "Empty string to gauge",
			args:    args{"cnt1", ""},
			want:    &Gauge{},
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
		want *Gauge
	}{
		{
			name: "New cnt1",
			args: args{"cnt1", ValType(100.99)},
			want: &Gauge{"cnt1", ValType(100.99)},
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
