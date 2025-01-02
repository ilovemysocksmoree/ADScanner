package ldap

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

type LDAPConfig struct {
	ServerAddr   string
	ServerPort   uint64
	BindUser     string
	BindPassword string
	Domain       string
	Conn         *ldap.Conn
}

func GetLDAPConfig() *LDAPConfig {
	return &LDAPConfig{
		ServerAddr:   "192.168.254.16",
		ServerPort:   389,
		BindUser:     "Administrator",
		BindPassword: "admin@123",
		Domain:       "baiman.local",
	}
}

func (lc *LDAPConfig) ConnectToServer() bool {
	fullAddr := fmt.Sprintf("%s:%d", lc.ServerAddr, lc.ServerPort)
	conn, err := ldap.Dial("tcp", fullAddr)
	if err != nil {
		fmt.Printf("Error while making connection request to LDAP server at: %s \n", fullAddr)
		fmt.Println(err.Error())

		return false
	}

	lc.Conn = conn
	return true
}

func (lc *LDAPConfig) Authenticate() error {
	fullUser := fmt.Sprintf("%s@%s", lc.BindUser, lc.Domain)

	err := lc.Conn.Bind(fullUser, lc.BindPassword)
	if err != nil {
		fmt.Printf("Error while authenticating user: %s \n", fullUser)
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (lc *LDAPConfig) SplitDomain() []string {
	return strings.Split(lc.Domain, ".")
}
