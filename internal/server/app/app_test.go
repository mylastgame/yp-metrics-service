package app

import (
	"github.com/mylastgame/yp-metrics-service/internal/domain/counter"
	"github.com/mylastgame/yp-metrics-service/internal/domain/gauge"
	counterRepo "github.com/mylastgame/yp-metrics-service/internal/storage/counter"
	gaugeRepo "github.com/mylastgame/yp-metrics-service/internal/storage/gauge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestApp_SaveCounter(t *testing.T) {
	type fields struct {
		gaugeRepo   gauge.Repo
		counterRepo counter.Repo
	}
	type args struct {
		title string
		val   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    []args
		want    []counter.ValType
		wantErr bool
	}{
		{
			name: "Save cnt1",
			fields: fields{
				gaugeRepo:   gaugeRepo.NewMemRepo(),
				counterRepo: counterRepo.NewMemRepo(),
			},
			args: []args{
				{title: "cnt1", val: "1"},
				{title: "cnt1", val: "2"},
				{title: "cnt2", val: "2"},
			},
			want:    []counter.ValType{counter.ValType(1), counter.ValType(3), counter.ValType(2)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				gaugeRepo:   tt.fields.gaugeRepo,
				counterRepo: tt.fields.counterRepo,
			}
			for i, arg := range tt.args {
				err := app.SaveCounter(arg.title, arg.val)
				require.NoError(t, err)
				c, ok := app.counterRepo.Get(arg.title)
				require.True(t, ok)
				assert.Equal(t, tt.want[i], c.Val)
			}

			c, ok := app.counterRepo.Get("cnt1")
			require.True(t, ok)
			assert.Equal(t, counter.ValType(3), c.Val)

		})
	}
}

func TestApp_SaveGauge(t *testing.T) {
	type fields struct {
		gaugeRepo   gauge.Repo
		counterRepo counter.Repo
	}
	type args struct {
		title string
		val   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    []args
		want    []gauge.ValType
		wantErr bool
	}{
		{
			name: "Save cnt1",
			fields: fields{
				gaugeRepo:   gaugeRepo.NewMemRepo(),
				counterRepo: counterRepo.NewMemRepo(),
			},
			args: []args{
				{title: "cnt1", val: "1"},
				{title: "cnt1", val: "2.00000099"},
				{title: "cnt2", val: "10.1"},
			},
			want:    []gauge.ValType{gauge.ValType(1), gauge.ValType(2.00000099), gauge.ValType(10.1)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				gaugeRepo:   tt.fields.gaugeRepo,
				counterRepo: tt.fields.counterRepo,
			}
			for i, arg := range tt.args {
				err := app.SaveGauge(arg.title, arg.val)
				require.NoError(t, err)
				c, ok := app.gaugeRepo.Get(arg.title)
				require.True(t, ok)
				assert.Equal(t, tt.want[i], c.Val)
			}

			c, ok := app.gaugeRepo.Get("cnt1")
			require.True(t, ok)
			assert.Equal(t, gauge.ValType(2.00000099), c.Val)

		})
	}
}
