package auth

import (
	log "github.com/cihub/seelog"
	"github.com/mysterium/node/openvpn"
	"github.com/mysterium/node/openvpn/middlewares"
	"net"
	"regexp"
)

type middleware struct {
	authenticator Authenticator
	connection    net.Conn
	username      string
	password      string
	state         middlewares.State
}

func NewMiddleware(authenticator Authenticator) openvpn.ManagementMiddleware {
	return &middleware{
		authenticator: authenticator,
		connection:    nil,
	}
}

func (m *middleware) Start(connection net.Conn) error {
	m.connection = connection
	log.Info("starting client user-pass authenticator middleware")

	return nil
}

func (m *middleware) Stop() error {
	return nil
}

func (m *middleware) State() middlewares.State {
	return m.state
}

func (m *middleware) ConsumeLine(line string) (consumed bool, err error) {
	//>PASSWORD:Need 'Auth' username/password
	rule, err := regexp.Compile("^>PASSWORD:Need 'Auth' username/password$")
	if err != nil {
		return false, err
	}

	match := rule.FindStringSubmatch(line)
	if len(match) > 0 {
		m.Reset()
		m.state = middlewares.STATE_AUTH
		username, password, err := m.authenticator()
		log.Info("authenticating user ", username, " with pass: ", password)

		_, err = m.connection.Write([]byte("password 'Auth' " + password + "\n"))
		if err != nil {
			return false, err
		}

		_, err = m.connection.Write([]byte("username 'Auth' " + username + "\n"))
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (m *middleware) Reset() {
	m.username = ""
	m.password = ""
	m.state = middlewares.STATE_WAIT
}