package bw

import (
	"testing"
)

func TestStatsForBW(t *testing.T) {
	type args struct {
		domain string
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
	}{
		{
			name:      "get stats for BW Dortmund",
			args:      args{
				domain: "www.boulderwelt-dortmund.de",
			},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStats, err := StatsForBW(tt.args.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("StatsForBW() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Got response: %v", gotStats)
		})
	}
}