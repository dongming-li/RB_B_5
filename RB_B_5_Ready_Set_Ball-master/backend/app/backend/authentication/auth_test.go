package authentication

import (
	"reflect"
	"testing"

	tr "git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/operations/transaction"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/model"
)

func Test_result_GetMeta(t *testing.T) {
	type fields struct {
		meta map[string]string
		data interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		// Test cases.
		{
			name:   "Metadata is returned",
			fields: fields{meta: map[string]string{"meta": "test", "data": "second"}, data: nil},
			want:   map[string]string{"meta": "test", "data": "second"},
		},
		// { //TODO: I don't know if I actually want this, because of nil checks
		// 	name:   "Nil metadat returns empty map",
		// 	fields: fields{meta: nil, data: nil},
		// 	want:   map[string]string{},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := result{
				meta: tt.fields.meta,
				data: tt.fields.data,
			}
			if got := r.GetMeta(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("result.GetMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_result_GetData(t *testing.T) {
	type fields struct {
		meta map[string]string
		data interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		// Test cases.
		{
			name:   "Data is returned",
			fields: fields{data: map[string]string{"meta": "test", "data": "second"}, meta: nil},
			want:   map[string]string{"meta": "test", "data": "second"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := result{
				meta: tt.fields.meta,
				data: tt.fields.data,
			}
			if got := r.GetData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("result.GetData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	type args struct {
		c    model.Collection
		cred map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    tr.Result
		wantErr bool
	}{
	//TODO: Test cases.
	// {},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Login(tt.args.c, tt.args.cred)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Login() = %v, want %v", got, tt.want)
			}
		})
	}
}
