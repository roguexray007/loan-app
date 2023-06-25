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

	"github.com/roguexray007/loan-app/internal/app/loans"
)

type LoanSuite struct {
	BaseSuite
}

func TestLoanSuite(t *testing.T) {
	suite.Run(t, new(LoanSuite))
}

func (s *LoanSuite) TestLoanFlow() {
	// create Loan
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
	var resp loans.Loan
	json.Unmarshal(response, &resp)
	assert.Len(s.T(), resp.LoanPayments, 3)

	// approve loan
	input = map[string]interface{}{
		"loan_id": resp.ID,
	}
	b, _ = json.Marshal(input)
	req = httptest.NewRequest(http.MethodPost, "/v1/loans/approve", bytes.NewReader(b))
	req.SetBasicAuth("admin", "adminpass")
	recorder = httptest.NewRecorder()
	s.Server.ServeHTTP(recorder, req)
	response, _ = ioutil.ReadAll(recorder.Body)
	expected = map[string]interface{}{
		"status": "approved",
	}
	AssertJSONSelectiveMatch(s.T(), response, expected)

	// pay loan 3 times
	for i := 1; i <= 3; i++ {
		input = map[string]interface{}{
			"loan_id": resp.ID,
			"amount":  500,
		}
		b, _ = json.Marshal(input)
		req = httptest.NewRequest(http.MethodPost, "/v1/loans/pay", bytes.NewReader(b))
		req.SetBasicAuth("loanuser", "loanpass")
		recorder = httptest.NewRecorder()
		s.Server.ServeHTTP(recorder, req)
		response, _ = ioutil.ReadAll(recorder.Body)
	}
	json.Unmarshal(response, &resp)
	expected = map[string]interface{}{
		"status": "paid",
	}
	AssertJSONSelectiveMatch(s.T(), response, expected)
}

/*
	Can add further tests. not adding right now. Sharing list of possible tests.
	1. db update failure,
    2. db txn failure
    3. loan already paid
    4. fail when loan is in pending state and user is trying to pay.
	5. validation failures
    6. loan_id not found in db
*/
