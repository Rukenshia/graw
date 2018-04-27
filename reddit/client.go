package reddit

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Rukenshia/graw/metrics"
	log "github.com/sirupsen/logrus"
)

// tokenURL is the url of reddit's oauth2 authorization service.
const tokenURL = "https://www.reddit.com/api/v1/access_token"

// clientConfig holds all the information needed to define Client behavior, such
// as who the client will identify as externally and where to authorize.
type clientConfig struct {
	// Agent is the user agent set in all requests made by the Client.
	agent string

	// If all fields in App are set, this client will attempt to identify as
	// a registered Reddit app using the credentials.
	app App
}

// client executes http Requests and invisibly handles OAuth2 authorization.
type client interface {
	Do(*http.Request) ([]byte, error)
}

type baseClient struct {
	cli *http.Client
}

func (b *baseClient) Do(req *http.Request) ([]byte, error) {
	resp, err := b.cli.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	metrics.Requests.WithLabelValues(req.URL.Path, fmt.Sprintf("%d", resp.StatusCode)).Inc()

	used := resp.Header.Get("X-Ratelimit-Used")
	remaining := resp.Header.Get("X-Ratelimit-Remaining")
	reset := resp.Header.Get("X-Ratelimit-Reset")

	if used != "" {
		n, err := strconv.Atoi(used)
		if err == nil {
			metrics.RateLimitUsed.Set(float64(n))
		}
	}
	if remaining != "" {
		n, err := strconv.Atoi(remaining)
		if err == nil {
			metrics.RateLimitRemaining.Set(float64(n))
		}
	}
	if reset != "" {
		n, err := strconv.Atoi(reset)
		if err == nil {
			metrics.RateLimitReset.Set(float64(n))
		}
	}

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusForbidden:
		return nil, PermissionDeniedErr
	case http.StatusServiceUnavailable:
		return nil, BusyErr
	case http.StatusTooManyRequests:
		return nil, RateLimitErr
	case http.StatusBadGateway:
		return nil, GatewayErr
	case http.StatusGatewayTimeout:
		return nil, GatewayTimeoutErr
	default:
		return nil, fmt.Errorf("bad response code: %d", resp.StatusCode)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"component": "graw",
		"body":      buf.String(),
	}).Debug("received response")

	metrics.ResponseSize.WithLabelValues(req.URL.Path, fmt.Sprintf("%d", resp.StatusCode)).Observe(float64(buf.Len()))

	return buf.Bytes(), nil
}

// newClient returns a new client using the given user to make requests.
func newClient(c clientConfig) (client, error) {
	if c.app.tokenURL == "" {
		c.app.tokenURL = tokenURL
	}

	if c.app.configured() {
		return newAppClient(c)
	}

	return &baseClient{clientWithAgent(c.agent)}, nil
}
