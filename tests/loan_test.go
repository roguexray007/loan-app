package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LoanSuite struct {
	BaseSuite
}

func TestLoanSuite(t *testing.T) {
	suite.Run(t, new(LoanSuite))
}

func (s *LoanSuite) TestCreateLoan() {

	input := map[string]interface{}{
		"amount": 1500,
		"terms":  3,
	}
	b, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/v1/loans", bytes.NewReader(b))
	req.SetBasicAuth("loanuser", "loanpass")
	recorder := httptest.NewRecorder()
	s.Server.ServeHTTP(recorder, req)
	response, _ := ioutil.ReadAll(recorder.Body)
	expected := map[string]interface{}{
		"amount":     1500,
		"status":     "pending",
		"terms_paid": 0,
		"loanPayments": []map[string]interface{}{
			{
				"status": "pending",
				"amount": 500,
			},
		},
	}
	AssertJSONSelectiveMatch(s.T(), response, expected)
	var resp map[string]interface{}
	json.Unmarshal(response, &resp)
	assert.Len(s.T(), resp["loanPayments"], 3)
}
