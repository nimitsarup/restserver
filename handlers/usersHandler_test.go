package handlers

import (
	"net/http"
	"testing"

	"github.com/nimitsarup/restserver/db"
	"github.com/stretchr/testify/assert"
)

func TestHandlers_AddUser(t *testing.T) {
	thePassword := "plainTextPassword"
	type fields struct {
		DB db.UsersInMemoryDB
	}
	type args struct {
		body          string
		checkPassword bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "password is stored in encrypted format",
			args:    args{body: `{"name": "John Doe", "email": "johndoe@postoffice.com", "password": "plainTextPassword"}`, checkPassword: true},
			wantErr: false,
			want:    http.StatusCreated,
			fields:  fields{DB: db.New()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				DB: tt.fields.DB,
			}
			got, err := h.AddUser(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handlers.AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Handlers.AddUser() = %v, want %v", got, tt.want)
			}
			if tt.args.checkPassword {
				_, users := h.GetAllUsers()
				assert.Equal(t, len(users), 1)
				assert.Equal(t, db.CheckPasswordHash(thePassword, users[0].Password), true)
			}
		})
	}
}
