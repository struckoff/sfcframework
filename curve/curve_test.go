package curve

import (
	"testing"
)

func TestDrawCurve(t *testing.T) {
	type args struct {
		cType CurveType
		bits  uint64
		op    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Hilbert 5 bits",
			args: args{
				cType: Hilbert,
				bits:  5,
				op:    "hilbert.png",
			},
			wantErr: false,
		},
		{
			name: "Morton 5 bits",
			args: args{
				cType: Morton,
				bits:  5,
				op:    "morton.png",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DrawCurve(tt.args.cType, tt.args.bits, tt.args.op); (err != nil) != tt.wantErr {
				t.Errorf("DrawCurve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDrawSplitCurve(t *testing.T) {
	type args struct {
		cType  CurveType
		bits   uint64
		splits []float64
		op     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Hilbert 8 bits",
			args: args{
				cType:  Hilbert,
				bits:   8,
				splits: []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0},
				op:     "hilbert-split.png",
			},
			wantErr: false,
		},
		{
			name: "Morton 8 bits",
			args: args{
				cType:  Morton,
				bits:   8,
				splits: []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0},
				op:     "morton-split.png",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DrawSplitCurve(tt.args.cType, tt.args.bits, tt.args.splits, tt.args.op); (err != nil) != tt.wantErr {
				t.Errorf("DrawSplitCurve() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
