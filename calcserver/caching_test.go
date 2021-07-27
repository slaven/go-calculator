package calcserver

import (
	"reflect"
	"testing"
)

func TestNewCache(t *testing.T) {
	got := reflect.TypeOf(NewCache()).String()
	want := "*calcserver.CalcCache"
	if got != want {
		t.Errorf("NewCache() type got: %v, type want: %v", got, want)
	}
}

func TestGetSet(t *testing.T) {
	type argsSet struct {
		key string
		val float64
	}

	tablesSet := []struct {
		args argsSet
		want bool
	}{
		{
			args: argsSet{
				key: "",
				val: 0,
			},
			want: true,
		},
		{
			args: argsSet{
				key: "test1",
				val: 2,
			},
			want: true,
		},
		{
			args: argsSet{
				key: "test1",
				val: 3,
			},
			want: true,
		},
		{
			args: argsSet{
				key: "test2",
				val: 4,
			},
			want: true,
		},
	}

	type argsGet struct {
		key string
	}
	tablesGet := []struct {
		args       argsGet
		wantResult float64
		wantBool   bool
	}{
		{
			args: argsGet{
				key: "x",
			},
			wantResult: 0,
			wantBool:   false,
		},
		{
			args: argsGet{
				key: "random",
			},
			wantResult: 0,
			wantBool:   false,
		},
		{
			args: argsGet{
				key: "test1",
			},
			wantResult: 3,
			wantBool:   true,
		},
		{
			args: argsGet{
				key: "test2",
			},
			wantResult: 4,
			wantBool:   true,
		},
	}

	cache := NewCache()

	for _, table := range tablesSet {
		got, _ := cache.SetOrUpdate(table.args.key, table.args.val)
		if got != table.want {
			t.Errorf("SetOrUpdate(%s, %v) got: %v, want: %v.", table.args.key, table.args.val, got, table.want)
		}
	}

	for _, table := range tablesGet {
		gotResult, gotBool := cache.Get(table.args.key)
		if gotResult != table.wantResult || gotBool != table.wantBool {
			t.Errorf("Get(%s) gotResult: %v, wantResult: %v, gotBool: %v, wantBool:%v.", table.args.key, gotResult, table.wantResult, gotBool, table.wantBool)
		}
	}
}

func TestBuildCacheKey(t *testing.T) {
	type args struct {
		calculation CalcOperation
		x           float64
		y           float64
	}

	tables := []struct {
		args args
		want string
	}{
		{
			args: args{
				calculation: addCalc,
				x:           1,
				y:           2},
			want: "add|1|2",
		},
		{
			args: args{
				calculation: addCalc,
				x:           2,
				y:           1},
			want: "add|1|2",
		},
		{
			args: args{
				calculation: multiplyCalc,
				x:           2,
				y:           3},
			want: "multiply|2|3",
		},
		{
			args: args{
				calculation: multiplyCalc,
				x:           3,
				y:           2},
			want: "multiply|2|3",
		},
		{
			args: args{
				calculation: subtractCalc,
				x:           2,
				y:           3},
			want: "subtract|2|3",
		},
		{
			args: args{
				calculation: subtractCalc,
				x:           3,
				y:           2},
			want: "subtract|3|2",
		},
		{
			args: args{
				calculation: divideCalc,
				x:           2,
				y:           4},
			want: "divide|2|4",
		},
		{
			args: args{
				calculation: divideCalc,
				x:           4,
				y:           2},
			want: "divide|4|2",
		},
	}

	for _, table := range tables {
		got := buildCacheKey(table.args.calculation, table.args.x, table.args.y)
		if got != table.want {
			t.Errorf("buildCacheKey(%s, %v, %v) got: %v, want: %v.", table.args.calculation, table.args.x, table.args.y, got, table.want)
		}
	}
}
