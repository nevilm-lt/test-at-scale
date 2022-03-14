package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/LambdaTest/synapse/pkg/lumber"
	"github.com/LambdaTest/synapse/pkg/service/teststats"
	"github.com/LambdaTest/synapse/testutils"
	"github.com/gin-gonic/gin"
)

// NOTE: Tests in this package are meant to be run in a Linux environment

func TestNewRouter(t *testing.T) {
	logger, _ := testutils.GetLogger()
	cfg, _ := testutils.GetConfig()
	ts, err := teststats.New(cfg, logger)
	if err != nil {
		t.Errorf("Error creating teststats service: %v", err)
	}
	type args struct {
		logger lumber.Logger
		ts     *teststats.ProcStats
	}
	tests := []struct {
		name string
		args args
		want Router
	}{
		{"TestNewRouter", args{logger, ts}, Router{logger, ts}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRouter(tt.args.logger, tt.args.ts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestRouter_Handler(t *testing.T) {
	logger, _ := testutils.GetLogger()
	cfg, _ := testutils.GetConfig()
	ts, err := teststats.New(cfg, logger)
	if err != nil {
		t.Errorf("Error creating teststats service: %v", err)
	}
	tests := []struct {
		name             string
		httpRequest      *http.Request
		wantResponseCode int
		wantStatusText   string
	}{
		{"Test handler health route for success", httptest.NewRequest(http.MethodGet, "/health", nil), 200, http.StatusText(http.StatusOK)},
		{"Test handler result route", httptest.NewRequest(http.MethodPost, "/results", bytes.NewBuffer([]byte(`{"TaskID" : "123"}`))), 200, http.StatusText(http.StatusOK)},
		{"Test handler result route for error in jsonBinding and hence http.StatusBadRequest", httptest.NewRequest(http.MethodPost, "/results", nil), http.StatusBadRequest, `{"message":"EOF"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newRouter := NewRouter(logger, ts)
			resp := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)
			c, _ := gin.CreateTestContext(resp)
			c.Request = tt.httpRequest
			newRouter.Handler().ServeHTTP(resp, c.Request)
			if resp.Code != tt.wantResponseCode {
				t.Errorf("Router.Handler() responseCode = %v, want = %v\n", resp.Code, tt.wantResponseCode)
				return
			}
			if !reflect.DeepEqual(resp.Body.String(), tt.wantStatusText) {
				t.Errorf("Router.Handler() statusText = %v, want = %v\n", resp.Body.String(), tt.wantStatusText)
			}
		})
	}
}