package libplumraw

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetLogicalLoadLevel(t *testing.T) {
	expPost := `{"level":123,"llid":"load-uuid"}`
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		bod, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, expPost, string(bod))
		w.WriteHeader(204)
	})
	pad := newMockLightpad(hf)
	pad.LLID = "load-uuid"
	err := pad.SetLogicalLoadLevel(123)
	assert.NoError(t, err)

}

func newMockLightpad(handler http.Handler) *DefaultLightpad {
	ts := httptest.NewTLSServer(handler)
	ipPort := strings.Split(strings.TrimPrefix(ts.URL, "https://"), ":")
	ip := net.ParseIP(ipPort[0])
	port, _ := strconv.Atoi(ipPort[1])
	pad := &DefaultLightpad{
		IP:   ip,
		Port: port,
	}
	return pad
}
