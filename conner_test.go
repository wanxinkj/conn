package conn

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewExternalProcedure(t *testing.T) {
	type args struct {
		conners []Conner
	}
	tests := []struct {
		name string
		args args
		want *ExternalProcedure
	}{
		{
			args:args{conners: []Conner{NewMysqlConfig("127.0.0.1", "3306", "uws_officebos", "f9nSziCuY5dan0AR", "qa_bos")}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewExternalProcedure(tt.args.conners...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewExternalProcedure() = %v, want %v", got, tt.want)
			}
			fmt.Println(GDB())
		})
	}
}