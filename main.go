package main

import (
	"fmt"

	ad "github.com/ilovemysocksmore/ADScanner/internal/server/AD"
	"github.com/ilovemysocksmore/ADScanner/internal/server/ldap"
)

func main() {
	ldapCfg := ldap.GetLDAPConfig()
	resp := ldapCfg.ConnectToServer()
	if !resp {
		fmt.Println("Unable to connect to LDAP server, check logs for more information")
		return
	}

	err := ldapCfg.Authenticate()
	if err != nil {
		fmt.Println("Unable to authenticate user")
		return
	}

	users, err := ad.GetAllUsers(ldapCfg)
	if err != nil {
		fmt.Println("Error while fetching all users...")
		return
	}

	fmt.Println(users)
}
