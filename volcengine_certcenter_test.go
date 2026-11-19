package volcenginecertcenter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/volcengine/volcengine-go-sdk/service/certificateservice"
	ve "github.com/volcengine/volcengine-go-sdk/volcengine"
	vesession "github.com/volcengine/volcengine-go-sdk/volcengine/session"
)

func TestUploadCertificate(t *testing.T) {
	for _, scenario := range []struct {
		name   string
		result map[string]string
		wantID string
	}{
		{"new certificate", map[string]string{"InstanceId": "new"}, "new"},
		{"existing certificate", map[string]string{"RepeatId": "existing"}, "existing"},
		{"missing certificate ID", map[string]string{}, ""},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				query := r.URL.Query()
				if r.Method != http.MethodPost || query.Get("Action") != "ImportCertificate" || query.Get("Version") != "2024-10-01" {
					t.Error("upload did not use the Certificate Center import endpoint")
				}
				var request certificateservice.ImportCertificateInput
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					t.Errorf("decode upload request: %v", err)
				}
				if request.Repeatable == nil || *request.Repeatable ||
					ve.StringValue(request.ProjectName) != "test-project" || request.CertificateInfo == nil ||
					ve.StringValue(request.CertificateInfo.CertificateChain) != "test-cert" ||
					ve.StringValue(request.CertificateInfo.PrivateKey) != "test-key" {
					t.Error("unexpected upload payload or duplicate policy")
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(map[string]any{"ResponseMetadata": map[string]string{"RequestId": "test"}, "Result": scenario.result}); err != nil {
					t.Errorf("respond to import request: %v", err)
				}
			}))
			defer server.Close()

			manager, err := NewCertmgr(&CertmgrConfig{ProjectName: "test-project"})
			if err != nil {
				t.Fatal(err)
			}
			config := ve.NewConfig().WithAkSk("test", "test").WithRegion("cn-beijing").WithEndpoint(server.URL)
			session, err := vesession.NewSession(config)
			if err != nil {
				t.Fatal(err)
			}
			manager.sdkClient = certificateservice.New(session, config)
			result, err := manager.Upload(t.Context(), "test-cert", "test-key")
			if scenario.wantID == "" {
				if err == nil || !strings.Contains(err.Error(), "empty certificate id") || result != nil {
					t.Fatalf("expected an empty ID to fail, got %v, %v", result, err)
				}
				return
			}
			if err != nil || result == nil || result.CertId != scenario.wantID {
				t.Fatalf("upload returned %v, %v, want %q", result, err, scenario.wantID)
			}
		})
	}
}
