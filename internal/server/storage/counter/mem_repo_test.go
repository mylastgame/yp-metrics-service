package counter

import (
	"github.com/mylastgame/yp-metrics-service/internal/server/domain/counter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemRepo_Get(t *testing.T) {
	type fields struct {
		storage map[string]counter.ValType
	}
	tests := []struct {
		name   string
		fields fields
		arg    string
		want   *counter.Counter
		want1  bool
	}{
		{
			name: "Get cnt1",
			fields: fields{
				storage: map[string]counter.ValType{
					"cnt1": 1,
					"cnt2": 4,
				},
			},
			arg:   "cnt1",
			want:  &counter.Counter{Title: "cnt1", Val: 1},
			want1: true,
		},
		{
			name: "Get cnt1",
			fields: fields{
				storage: map[string]counter.ValType{
					"cnt1": 1,
					"cnt2": 4,
				},
			},
			arg:   "cnt3",
			want:  &counter.Counter{},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MemRepo{
				storage: tt.fields.storage,
			}
			got, got1 := r.Get(tt.arg)

			require.Equal(t, tt.want1, got1)

			if tt.want1 {
				assert.Equal(t, tt.want, got)
				assert.Equal(t, tt.arg, got.Title)
				assert.Equal(t, counter.ValType(1), got.Val)
			}

		})
	}
}

func TestMemRepo_Save(t *testing.T) {
	type fields struct {
		storage map[string]counter.ValType
	}
	type args struct {
		v []counter.Counter
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantVal []counter.ValType
		wantErr bool
	}{
		{
			name: "Save cnt1, cnt2",
			fields: fields{
				storage: map[string]counter.ValType{},
			},
			args: args{
				v: []counter.Counter{
					{Title: "cnt1", Val: 1},
					{Title: "cnt2", Val: 1},
				},
			},
			wantVal: []counter.ValType{counter.ValType(1), counter.ValType(1)},
			wantErr: false,
		},
		{
			name: "Increment only cnt1",
			fields: fields{
				storage: map[string]counter.ValType{},
			},
			args: args{
				v: []counter.Counter{
					{Title: "cnt1", Val: 1},
					{Title: "cnt1", Val: 1},
					{Title: "cnt2", Val: 1},
				},
			},
			wantVal: []counter.ValType{counter.ValType(1), counter.ValType(2), counter.ValType(2)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MemRepo{
				storage: tt.fields.storage,
			}

			for i, a := range tt.args.v {
				err := r.Add(&a)
				require.NoError(t, err)
				assert.Equal(t, tt.wantVal[i], r.storage["cnt1"])
			}
		})
	}
}
