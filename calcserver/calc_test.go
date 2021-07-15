package calcserver

import "testing"

func TestCalculate(t *testing.T) {
	type args struct {
		calculation CalcOperation
		x           float64
		y           float64
	}

	tables := []struct {
		args    args
		want    float64
		wantErr bool
	}{
		{
			args: args{
				calculation: addCalc,
				x:           1,
				y:           1},
			want:    2,
			wantErr: false,
		},
		{
			args: args{
				calculation: subtractCalc,
				x:           2,
				y:           2},
			want:    0,
			wantErr: false,
		},
		{
			args: args{
				calculation: multiplyCalc,
				x:           2,
				y:           2},
			want:    4,
			wantErr: false,
		},
		{
			args: args{
				calculation: divideCalc,
				x:           9,
				y:           3},
			want:    3,
			wantErr: false,
		},
		{
			args: args{
				calculation: divideCalc,
				x:           2,
				y:           0},
			want:    0,
			wantErr: true,
		},
		{
			args: args{
				calculation: "unknown",
				x:           2,
				y:           0},
			want:    0,
			wantErr: true,
		},
	}

	for _, table := range tables {
		got, gotErrRaw := calculate(table.args.calculation, table.args.x, table.args.y)
		gotErr := gotErrRaw != nil
		if got != table.want || gotErr != table.wantErr {
			t.Errorf("calculate(%s, %v, %v) got: %v, gotErr: %v, want: %v, wantErr: %v.", table.args.calculation, table.args.x, table.args.y, got, gotErr, table.want, table.wantErr)
		}
	}
}
