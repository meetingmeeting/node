package auth

import (
	"github.com/mysterium/node/openvpn/middlewares"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

var currentState middlewares.State

func NewFakeAuthenticator() Authenticator {
	return func(username, password string) (bool, error) {
		if username == "bad" {
			return false, nil
		}

		return true, nil
	}
}

type fakeAuthenticator struct {
	username string
	password string
}

func (a *fakeAuthenticator) auth() (username string, password string, err error) {
	return
}

func (a *fakeAuthenticator) authWithValid() (username string, password string, err error) {
	username = "valid_username"
	password = "valid_password"
	return
}

type fakeConnection struct {
	lastDataWritten []byte
	net.Conn
}

func (conn *fakeConnection) Read(b []byte) (int, error) {
	return 0, nil
}

func (conn *fakeConnection) Write(b []byte) (n int, err error) {
	conn.lastDataWritten = b
	return 0, nil
}

func Test_Factory(t *testing.T) {
	authenticator := NewFakeAuthenticator()
	middleware := NewMiddleware(authenticator)
	assert.NotNil(t, middleware)
}

func Test_ConsumeLineSkips(t *testing.T) {
	var tests = []struct {
		line string
	}{
		{">SOME_LINE_DELIVERED"},
		{">ANOTHER_LINE_DELIVERED"},
	}
	authenticator := NewFakeAuthenticator()
	middleware := NewMiddleware(authenticator)

	for _, test := range tests {
		consumed, err := middleware.ConsumeLine(test.line)
		assert.NoError(t, err, test.line)
		assert.False(t, consumed, test.line)
	}
}

func Test_ConsumeLineTakes(t *testing.T) {
	var tests = []struct {
		line          string
		expectedState middlewares.State
	}{
		{">CLIENT:REAUTH,0,0", middlewares.STATE_AUTH},
		{">CLIENT:CONNECT,0,0", middlewares.STATE_AUTH},
		{">CLIENT:ENV,password=(.*)", middlewares.STATE_AUTH},
		{">CLIENT:ENV,username=(.*)", middlewares.STATE_AUTH},
	}

	authenticator := NewFakeAuthenticator()
	middleware := NewMiddleware(authenticator)
	connection := &fakeConnection{}
	middleware.Start(connection)

	for _, test := range tests {
		consumed, err := middleware.ConsumeLine(test.line)
		assert.NoError(t, err, test.line)
		assert.True(t, consumed, test.line)
		assert.Equal(t, test.expectedState, middleware.State())
	}
}