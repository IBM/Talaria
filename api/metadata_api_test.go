package api

import (
	"reflect"
	"testing"

	"github.com/ibm/opentalaria/config"
	"github.com/ibm/opentalaria/protocol"
)

func getMockHeader(version int, requestApiKey, requestApiVersion int16, correlationId int32) protocol.RequestHeader {
	clientId := "test-client"

	return protocol.RequestHeader{
		Version:           int16(version),
		RequestApiKey:     requestApiKey,
		RequestApiVersion: requestApiVersion,
		CorrelationID:     correlationId,
		ClientID:          &clientId,
	}
}

func TestMetadataAPI_GetRequest(t *testing.T) {
	type fields struct {
		Request Request
	}
	tests := []struct {
		name      string
		fields    fields
		want      Request
		happyPath bool
	}{
		{
			name: "Happy path",
			fields: fields{
				Request: Request{
					Header:  getMockHeader(1, 0, 1, 1),
					Message: []byte{0, 1, 2},
					Conn:    nil,
				},
			},
			want: Request{
				Header:  getMockHeader(1, 0, 1, 1),
				Message: []byte{0, 1, 2},
				Conn:    nil,
			},
			happyPath: true,
		},
		{
			name: "Wrong header",
			fields: fields{
				Request: Request{
					Header:  getMockHeader(1, 0, 1, 1),
					Message: []byte{0, 1, 2},
					Conn:    nil,
				},
			},
			want: Request{
				Header:  getMockHeader(1, 0, 0, 0),
				Message: []byte{0, 1, 2},
				Conn:    nil,
			},
			happyPath: false,
		},
		{
			name: "Wrong message",
			fields: fields{
				Request: Request{
					Header:  getMockHeader(1, 0, 1, 1),
					Message: []byte{0, 1, 2},
					Conn:    nil,
				},
			},
			want: Request{
				Header:  getMockHeader(1, 0, 1, 1),
				Message: []byte{1},
				Conn:    nil,
			},
			happyPath: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MetadataAPI{
				Request: tt.fields.Request,
			}
			if got := m.GetRequest(); reflect.DeepEqual(got, tt.want) != tt.happyPath {
				t.Errorf("MetadataAPI.GetRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetadataAPI_GetHeaderVersion(t *testing.T) {
	type fields struct {
		Request Request
	}
	type args struct {
		requestVersion int16
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      int16
		happyPath bool
	}{
		{
			name: "Happy path v0",
			fields: fields{
				Request: Request{
					Header: getMockHeader(1, 0, 0, 0),
				},
			},
			args: args{
				requestVersion: 1,
			},
			want: 0,
		},
		{
			name: "Happy path v1",
			fields: fields{
				Request: Request{
					Header: getMockHeader(1, 0, 0, 0),
				},
			},
			args: args{
				requestVersion: 9,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MetadataAPI{
				Request: tt.fields.Request,
			}
			if got := m.GetHeaderVersion(tt.args.requestVersion); got != tt.want {
				t.Errorf("MetadataAPI.GetHeaderVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetadataAPI_GeneratePayload(t *testing.T) {
	// This test is not very flexible, if the mock objects are changed in the config package, the expected payload has to be regenerated.
	type fields struct {
		Request Request
	}

	config := config.MockConfig()

	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Happy path",
			fields: fields{
				Request: Request{
					Header: getMockHeader(1, 3, 0, 0),
				},
			},
			want:    []byte{0, 0, 0, 1, 0, 0, 0, 1, 0, 9, 49, 50, 55, 46, 48, 46, 48, 46, 49, 0, 0, 35, 132, 0, 0, 0, 1, 0, 0, 0, 10, 116, 101, 115, 116, 45, 116, 111, 112, 105, 99, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MetadataAPI{
				Request: tt.fields.Request,
				Config:  config,
			}
			got, err := m.GeneratePayload()
			if (err != nil) != tt.wantErr {
				t.Errorf("MetadataAPI.GeneratePayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetadataAPI.GeneratePayload() = %v, want %v", got, tt.want)
			}
		})
	}
}
