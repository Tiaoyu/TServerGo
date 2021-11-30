package main

import (
	"reflect"
	"testing"
)

func TestMain(t *testing.M) {
	initNotify()
	initMatch()
	initRoom()
	initUser()
}

func TestGetPlayerByOpenId(t *testing.T) {
	type args struct {
		openId string
	}
	tests := []struct {
		name string
		args args
		want *Player
	}{
		{
			name: "GetPlayerByOpenId",
			args: args{
				openId: "1",
			},
			want: &Player{
				OpenId: "1",
				Sess: &UserSession{
					Conn:        nil,
					OpenId:      "",
					SendChannel: nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPlayerByOpenId(tt.args.openId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPlayerByOpenId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlayerLogin(t *testing.T) {
	type args struct {
		u *Player
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PlayerLogin(tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("PlayerLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_initUser(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_onPlayerLogout(t *testing.T) {
	type args struct {
		params []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
