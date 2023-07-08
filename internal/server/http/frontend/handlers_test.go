package frontend

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akrillis/k8storage/internal/entities"
	"github.com/akrillis/k8storage/internal/service"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_putFile(t *testing.T) {
	tests := []struct {
		name      string
		body      interface{}
		wantCode  int
		setupMock func(
			mReceiver *service.MockReceiver,
		)
	}{
		{
			name:     "parse error",
			body:     123,
			wantCode: http.StatusBadRequest,
			setupMock: func(
				mReceiver *service.MockReceiver,
			) {
				mReceiver.EXPECT().Put(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "client_id is empty",
			body: entities.PutFileRequest{
				Name: "test",
				Data: []byte("test"),
			},
			wantCode: http.StatusBadRequest,
			setupMock: func(
				mReceiver *service.MockReceiver,
			) {
				mReceiver.EXPECT().Put(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "name is empty",
			body: entities.PutFileRequest{
				ClientID: "1",
				Data:     []byte("test"),
			},
			wantCode: http.StatusBadRequest,
			setupMock: func(
				mReceiver *service.MockReceiver,
			) {
				mReceiver.EXPECT().Put(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "data is empty",
			body: entities.PutFileRequest{
				ClientID: "1",
				Name:     "test",
			},
			wantCode: http.StatusBadRequest,
			setupMock: func(
				mReceiver *service.MockReceiver,
			) {
				mReceiver.EXPECT().Put(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "service level error",
			body: entities.PutFileRequest{
				ClientID: "1",
				Name:     "test",
				Data:     []byte("test"),
			},
			wantCode: http.StatusInternalServerError,
			setupMock: func(
				mReceiver *service.MockReceiver,
			) {
				mReceiver.
					EXPECT().
					Put(gomock.Any(), gomock.Any()).
					Return(errors.New("test"))
			},
		},
		{
			name: "ok",
			body: entities.PutFileRequest{
				ClientID: "1",
				Name:     "test",
				Data:     []byte("test"),
			},
			wantCode: http.StatusOK,
			setupMock: func(
				mReceiver *service.MockReceiver,
			) {
				mReceiver.
					EXPECT().
					Put(gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			receiver := service.NewMockReceiver(ctrl)
			test.setupMock(receiver)

			h := &handlers{receiver: receiver}

			curBody, err := json.Marshal(test.body)
			require.NoError(t, err)
			request := httptest.NewRequest(http.MethodPost, urlFiles, bytes.NewBuffer(curBody))

			response := httptest.NewRecorder()
			http.HandlerFunc(h.putFile).ServeHTTP(response, request)
			require.Equal(t, test.wantCode, response.Code)
		})
	}
}

func Test_getFile(t *testing.T) {
	tests := []struct {
		name      string
		request   *http.Request
		wantCode  int
		setupMock func(
			mRestorer *service.MockRestorer,
		)
	}{
		{
			name:     "parse error",
			request:  httptest.NewRequest(http.MethodGet, urlFiles, nil),
			wantCode: http.StatusBadRequest,
			setupMock: func(
				mRestorer *service.MockRestorer,
			) {
				mRestorer.EXPECT().Get(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:     "client_id is empty",
			request:  httptest.NewRequest(http.MethodGet, urlFiles+"?name=test", nil),
			wantCode: http.StatusBadRequest,
			setupMock: func(
				mRestorer *service.MockRestorer,
			) {
				mRestorer.EXPECT().Get(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:     "name is empty",
			request:  httptest.NewRequest(http.MethodGet, urlFiles+"?client_id=1", nil),
			wantCode: http.StatusBadRequest,
			setupMock: func(
				mRestorer *service.MockRestorer,
			) {
				mRestorer.EXPECT().Get(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:     "service level error",
			request:  httptest.NewRequest(http.MethodGet, urlFiles+"?client_id=1&name=test", nil),
			wantCode: http.StatusInternalServerError,
			setupMock: func(
				mRestorer *service.MockRestorer,
			) {
				mRestorer.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("test"))
			},
		},
		{
			name:     "no content",
			request:  httptest.NewRequest(http.MethodGet, urlFiles+"?client_id=1&name=test", nil),
			wantCode: http.StatusNoContent,
			setupMock: func(
				mRestorer *service.MockRestorer,
			) {
				mRestorer.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return([]byte(nil), nil)
			},
		},
		{
			name:     "ok",
			request:  httptest.NewRequest(http.MethodGet, urlFiles+"?client_id=1&name=test", nil),
			wantCode: http.StatusOK,
			setupMock: func(
				mRestorer *service.MockRestorer,
			) {
				mRestorer.
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return([]byte("ksdjcbsdkjc"), nil)
			},
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			restorer := service.NewMockRestorer(ctrl)
			test.setupMock(restorer)

			h := &handlers{restorer: restorer}

			response := httptest.NewRecorder()
			http.HandlerFunc(h.getFile).ServeHTTP(response, test.request)
			require.Equal(t, test.wantCode, response.Code)
		})
	}
}
