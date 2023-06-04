package tokenutil

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDecodeUserIDFromToken(t *testing.T) {
	userID := uuid.New().String()

	type args struct {
		token string
	}
	tests := []struct {
		name       string
		args       args
		want       string
		wantErr    bool
		errMessage string
	}{
		{
			name:       "positive",
			args:       args{token: userID + "." + "0000000"},
			want:       userID,
			wantErr:    false,
			errMessage: "",
		},
		{
			name:       "empty token",
			args:       args{token: ""},
			want:       "",
			wantErr:    true,
			errMessage: "empty token",
		},
		{
			name:       "cant decode user id from token",
			args:       args{token: "12345678"},
			want:       "",
			wantErr:    true,
			errMessage: "cant decode user id from token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeUserIDFromToken(tt.args.token)

			if err != nil {
				assert.Equal(t, tt.wantErr, true)
				assert.Equal(t, tt.errMessage, err.Error())
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsTokenValid(t *testing.T) {
	userID := uuid.New().String()
	tokenValid := GenerateTokenForUser(userID)
	tokenInvalid := userID + "12345"

	type args struct {
		token string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid token",
			args: args{token: tokenValid},
			want: true,
		},
		{
			name: "invalid token",
			args: args{token: tokenInvalid},
			want: false,
		},
		{
			name: "empty token",
			args: args{token: ""},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsTokenValid(tt.args.token)
			assert.Equal(t, tt.want, got)
		})
	}
}
