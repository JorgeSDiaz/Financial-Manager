package list_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/financial-manager/api/cmd/api/handlers/account/list"
	"github.com/financial-manager/api/cmd/api/handlers/account/response"
)

// listResponse mirrors the handler's unexported listResponse for test decoding.
type listResponse struct {
	Accounts      []response.Account `json:"accounts"`
	GlobalBalance float64            `json:"global_balance"`
}

func TestHandler_Handle(t *testing.T) {
	t.Parallel()

	account := buildDomainAccount("a1", "Cash")
	accountResp := response.ToAccount(account)

	tests := []struct {
		name         string
		fakeLister   *fakeLister
		fakeBalancer *fakeBalanceGetter
		wantStatus   int
		wantBody     any
	}{
		{
			name:         "returns 200 with accounts and global balance",
			fakeLister:   &fakeLister{out: buildListOutput(account)},
			fakeBalancer: &fakeBalanceGetter{out: buildBalanceOutput(1000.0)},
			wantStatus:   http.StatusOK,
			wantBody: listResponse{
				Accounts:      []response.Account{accountResp},
				GlobalBalance: 1000.0,
			},
		},
		{
			name:         "empty list returns 200 with empty accounts and zero balance",
			fakeLister:   &fakeLister{out: buildListOutput()},
			fakeBalancer: &fakeBalanceGetter{out: buildBalanceOutput(0.0)},
			wantStatus:   http.StatusOK,
			wantBody: listResponse{
				Accounts:      []response.Account{},
				GlobalBalance: 0.0,
			},
		},
		{
			name:         "lister error returns 500",
			fakeLister:   &fakeLister{err: errors.New("db error")},
			fakeBalancer: &fakeBalanceGetter{},
			wantStatus:   http.StatusInternalServerError,
			wantBody:     response.Error{Error: "internal server error"},
		},
		{
			name:         "balance getter error returns 500",
			fakeLister:   &fakeLister{out: buildListOutput()},
			fakeBalancer: &fakeBalanceGetter{err: errors.New("db error")},
			wantStatus:   http.StatusInternalServerError,
			wantBody:     response.Error{Error: "internal server error"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := list.New(tc.fakeLister, tc.fakeBalancer)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/accounts", nil)
			rec := httptest.NewRecorder()

			h.Handle(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Equal(t, tc.wantBody, decodeAs(t, rec, tc.wantBody))
		})
	}
}

// decodeAs decodes the recorder body into a new value of the same type as want.
func decodeAs(t *testing.T, rec *httptest.ResponseRecorder, want any) any {
	t.Helper()
	ptr := reflect.New(reflect.TypeOf(want))
	require.NoError(t, json.NewDecoder(rec.Body).Decode(ptr.Interface()))
	return ptr.Elem().Interface()
}
