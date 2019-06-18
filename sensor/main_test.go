package main

import (
	"testing"
	. "github.com/metalblueberry/PeePooMonitor/sensor/mocks"
	. "github.com/metalblueberry/PeePooMonitor/sensor/hcsr51"
)

func Test_watchInputChanges(t *testing.T) {
	type args struct {
		pinNumber int
		cmd       MockCommander
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c Commander 
			c = &tt.args.cmd
			got, err := WatchInputChanges(tt.args.pinNumber, c)
			if (err != nil) != tt.wantErr {
				t.Errorf("watchInputChanges() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("watchInputChanges() = %v, want %v", got, tt.want)
			}
		})
	}
}
