package db

import (
	"reflect"
	"testing"

	"github.com/nimitsarup/restserver/model"
)

func TestInMemoryDB_GetUser(t *testing.T) {
	testUser := model.User{Name: "expected_name"}
	testUserId := "user-id"
	type fields struct {
		users map[string]model.User
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.User
		wantErr bool
	}{
		{
			name:    "WHENUserDoesNotExistsTHENCheckErr",
			wantErr: true,
		},
		{
			name:    "WHENUserExistsTHENCheckRetrival",
			wantErr: false,
			fields:  fields{users: map[string]model.User{testUserId: testUser}},
			want:    testUser,
			args:    args{id: testUserId},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &InMemoryDB{
				users: tt.fields.users,
			}
			got, err := d.GetUser(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("InMemoryDB.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InMemoryDB.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
