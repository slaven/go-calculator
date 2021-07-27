package calcserver

import (
	"net/url"
	"reflect"
	"testing"
)

func TestCreate(t *testing.T) {
	got := reflect.TypeOf(Create()).String()
	want := "*calcserver.CalcServer"
	if got != want {
		t.Errorf("Create() type got: %v, type want: %v", got, want)
	}
}

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

func TestGetSupportedRoutes(t *testing.T) {

	got := getSupportedRoutes()
	want := []CalcOperation{addCalc, subtractCalc, multiplyCalc, divideCalc}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("getSupportedRoutes() does not return all required calc operations")
	}
}

func TestGetQueryValues(t *testing.T) {
	urlVals := func(x string, y string) url.Values {
		v := url.Values{}
		v.Set("x", x)
		v.Set("y", y)
		return v
	}

	type args struct {
		vals url.Values
	}
	tables := []struct {
		args    args
		wantX   float64
		wantY   float64
		wantErr bool
	}{
		{
			args: args{
				vals: urlVals("1", "2"),
			},
			wantX:   1,
			wantY:   2,
			wantErr: false,
		},
		{
			args: args{
				vals: urlVals("111", ""),
			},
			wantX:   0,
			wantY:   0,
			wantErr: true,
		},

		{
			args: args{
				vals: urlVals("aaa", ""),
			},
			wantX:   0,
			wantY:   0,
			wantErr: true,
		},
		{
			args: args{
				vals: urlVals("", "222"),
			},
			wantX:   0,
			wantY:   0,
			wantErr: true,
		},
		{
			args: args{
				vals: urlVals("", "bbb"),
			},
			wantX:   0,
			wantY:   0,
			wantErr: true,
		},
		{
			args: args{
				vals: urlVals("2", "bbb"),
			},
			wantX:   0,
			wantY:   0,
			wantErr: true,
		},
	}

	for _, table := range tables {
		gotX, gotY, gotErrRaw := getQueryValues(table.args.vals)
		gotErr := gotErrRaw != nil
		if gotX != table.wantX || gotY != table.wantY || table.wantErr != gotErr {
			t.Errorf("getQueryValues(%v) gotX: %v, wantX: %v, gotY: %v, wantY: %v, gotErr: %v, wantErr: %v", table.args, gotX, table.wantX, gotY, table.wantY, gotErr, table.wantErr)
		}
	}
}
