package balancer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange_Fits(t *testing.T) {
	type fields struct {
		Min uint64
		Max uint64
		Len uint64
	}
	type args struct {
		index uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "true",
			fields: fields{
				Min: 10,
				Max: 15,
				Len: 15 - 10,
			},
			args: args{
				index: 10,
			},
			want: true,
		},
		{
			name: "false",
			fields: fields{
				Min: 10,
				Max: 15,
				Len: 15 - 10,
			},
			args: args{
				index: 15,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Range{
				Min: tt.fields.Min,
				Max: tt.fields.Max,
				Len: tt.fields.Len,
			}
			assert.Equal(t, tt.want, r.Fits(tt.args.index))
		})
	}
}
