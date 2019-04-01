package pact

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/devenney/form3/common"
	"github.com/devenney/form3/server"
	"github.com/devenney/form3/utils"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
)

func TestMain(m *testing.M) {
	common.InitConfig()

	m.Run()
}

/*** A mock client to produce a pact. Run once, output treated as third-party. ***/
/*
func TestMockConsumerGetPaymentPact(t *testing.T) {
	// Create Pact connecting to local Daemon
	pact := &dsl.Pact{
		Consumer: "MockConsumer",
		Provider: "Payments",
		Host:     "localhost",
	}
	defer pact.Teardown()

	payment, err := payments.Get("09a8fe0d-e239-4aff-8098-7923eadd0b98")
	if err != nil {
		t.Fatal(err)
	}

	// Pass in test case. This is the component that makes the external HTTP call
	var test = func() (err error) {
		u := fmt.Sprintf("http://localhost:%d/getPayment/%s", pact.Server.Port, payment.ID)
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return err
		}

		// NOTE: by default, request bodies are expected to be sent with a Content-Type
		// of application/json. If you don't explicitly set the content-type, you
		// will get a mismatch during Verification.
		req.Header.Set("Content-Type", "application/json")

		if _, err = http.DefaultClient.Do(req); err != nil {
			return err
		}
		return nil
	}

	// Set up our expected interactions.
	pact.
		AddInteraction().
		Given("Payment 09a8fe0d-e239-4aff-8098-7923eadd0b98 exists").
		UponReceiving("A request to get 09a8fe0d-e239-4aff-8098-7923eadd0b98").
		WithRequest(dsl.Request{
			Method:  "GET",
			Path:    dsl.String("/getPayment/09a8fe0d-e239-4aff-8098-7923eadd0b98"),
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
		}).
		WillRespondWith(dsl.Response{
			Status:  200,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
			Body:    dsl.Match(&payments.Payment{}),
		})

	// Run the test, verify it did what we expected and capture the contract
	if err := pact.Verify(test); err != nil {
		log.Fatalf("Error on Verify: %v", err)
	}
}
*/

func TestProviderGetPayment(t *testing.T) {
	// Create Pact connecting to local Daemon
	pact := &dsl.Pact{
		Consumer: "MockConsumer",
		Provider: "Payments",
	}

	// Start provider API in the background
	go server.StartServer()

	// Ensure our mock data is populated
	utils.InsertMockData("../../payments/test_data/mock_data.json")

	// Verify the Provider using the locally saved Pact Files
	pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL: "http://localhost:8000",
		PactURLs:        []string{filepath.ToSlash(fmt.Sprintf("%s/mockconsumer-payments.json", "pacts"))},
	})
}
