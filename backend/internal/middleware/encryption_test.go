package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShouldDecryptRequestBody(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		want        bool
	}{
		{
			name:        "encrypted JSON request",
			contentType: "text/plain",
			want:        true,
		},
		{
			name:        "multipart upload with boundary",
			contentType: "multipart/form-data; boundary=upload-boundary",
			want:        false,
		},
		{
			name:        "multipart upload case insensitive",
			contentType: "Multipart/Form-Data; boundary=upload-boundary",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/upload/image", nil)
			req.Header.Set("Content-Type", tt.contentType)

			if got := shouldDecryptRequestBody(req); got != tt.want {
				t.Fatalf("shouldDecryptRequestBody() = %v, want %v", got, tt.want)
			}
		})
	}
}
