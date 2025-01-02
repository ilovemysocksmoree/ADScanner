package ad

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-ldap/ldap/v3"
	LDAP "github.com/ilovemysocksmore/ADScanner/internal/server/ldap"
)

type DomainUser struct {
	DomainName        string `json:"domain_name"` // in DN format
	Username          string `json:"username"`    // middlename
	Password          string `json:"password"`
	DisplayName       string `json:"display_name"`        // full name
	GivenName         string `json:"given_name"`          // first name
	Surname           string `json:"surname"`             // last name
	SAMAccountName    string `json:"sam_account_name"`    // firstname
	UserPrincipalName string `json:"user_principal_name"` // givenname@domain
	LogonCount        uint64 `json:"logon_count"`
	LastLogon         string `json:"last_logon"`

	OtherAttrs map[string][]string `json:"other_attrs"`
}

func AddAUser(cfg *LDAP.LDAPConfig) {

}

func GetAllUsers(cfg *LDAP.LDAPConfig) ([]DomainUser, error) {
	dn := cfg.SplitDomain()
	searchReq := ldap.NewSearchRequest(
		fmt.Sprintf("DC=%s,DC=%s", dn[0], dn[1]),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=user)",
		[]string{"cn", "sAMAccountName", "displayName", "mail", "memberOf", "*"},
		nil,
	)

	resp, err := cfg.Conn.Search(searchReq)
	if err != nil {
		fmt.Println("Error while processing search request")
		fmt.Println(err.Error())
		return nil, fmt.Errorf("error processing request: %w", err)
	}

	users := make([]DomainUser, 0, len(resp.Entries))
	for _, entry := range resp.Entries {
		logonCount, _ := strconv.ParseUint(entry.GetAttributeValue("logonCount"), 10, 64)
		usr := DomainUser{
			DomainName:        entry.DN,
			Username:          entry.GetAttributeValue("initials"),
			SAMAccountName:    entry.GetAttributeValue("sAMAccountName"),
			Surname:           entry.GetAttributeValue("sn"),
			GivenName:         entry.GetAttributeValue("givenName"),
			DisplayName:       entry.GetAttributeValue("displayName"),
			UserPrincipalName: entry.GetAttributeValue("userPrincipalName"),
			LogonCount:        logonCount,
			LastLogon:         entry.GetAttributeValue("lastLogonTimestamp"),
		}

		if len(entry.Attributes) > 0 {
			usr.OtherAttrs = make(map[string][]string)
			for _, attr := range entry.Attributes {
				if attr.Name == "" {
					continue
				}

				usr.OtherAttrs[attr.Name] = attr.Values
			}
		}

		// err := usr.ValidateDomainUser()
		users = append(users, usr)
	}

	return users, nil
}

func (DU *DomainUser) ValidateDomainUser() error {
	if DU == nil {
		return fmt.Errorf("User cannot be null")
	}

	required := map[string]string{
		"SAMAccountName":    DU.SAMAccountName,
		"DisplayName":       DU.DisplayName,
		"UserPrincipalName": DU.UserPrincipalName,
		"GivenName":         DU.GivenName,
		"Surname":           DU.Surname,
	}

	for x, y := range required {
		if strings.TrimSpace(y) == "" {
			return fmt.Errorf("Required field %s is empty", x)
		}
	}

	return nil
}

func (DU *DomainUser) CheckIfUserExist() error {

	return nil
}
