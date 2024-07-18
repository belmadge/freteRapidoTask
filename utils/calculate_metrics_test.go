package utils

import (
	"reflect"
	"testing"

	"github.com/belmadge/freteRapido/domain"
)

func TestCalculateMetrics(t *testing.T) {
	type args struct {
		quotes []domain.Quote
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateMetrics(tt.args.quotes)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateMetrics() got = %v, want %v", got, tt.want)
			}
		})
	}
}
