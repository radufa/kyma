package testsuite

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const (
	// Should match event type in tests/function-controller/pkg/subscription/subscription.go
	eventType = "sap.kyma.custom.something.order.created.v1"
)

func CreateEvent(url string) error {

	payload := fmt.Sprintf(`{ "%s": "%s" }`, TestDataKey, EventPing)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(payload))
	if err != nil {
		return fmt.Errorf("while creating new request: method %s, url %s, payload %s", http.MethodPost, url, payload)
	}

	// headers taken from example from documentation
	req.Header.Add("x-b3-flags", "1")
	req.Header.Add("ce-specversion", "1.0")
	req.Header.Add("ce-type", eventType)
	req.Header.Add("ce-time", "2018-04-05T03:56:24Z")
	req.Header.Add("ce-id", "45a8b444-3213-4758-be3f-540bf93f85ff")
	req.Header.Add("ce-source", "kyma")
	req.Header.Add("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "while making request to NATS publisher %s", url)
	}

	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("invalid response status %s while making a request to %s", resp.Status, url)
	}
	return nil
}
