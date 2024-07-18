package service

import (
	"reflect"
	"testing"

	"github.com/belmadge/freteRapido/domain"
)

func TestCreateQuote(t *testing.T) {
	type args struct {
		input domain.QuoteRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.QuoteResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateQuote(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateQuote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateQuote() got = %v, want %v", got, tt.want)
			}
		})
	}
}
