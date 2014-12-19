/*
Copyright 2014 Celluloid VFX, Berlin and Johannes Amorosa

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

/*
This is a commandline tool written for Celluloid VFX. We are using this to query
the LDAP server for two attributes. First attribute verifies if a user account
is enabled or disabled. The Second one describes if a Dial in is set and true.

If the conditions are met we return 0. If the conditions ar not met we return 1
*/

/*
Active Directory Attributes:

msNPAllowDialin: [TRUE]
msNPAllowDialin: [FALSE]
msNPAllowDialin: Not available

userAccountControl: 512 (user enabled)
userAccountControl: 514 (user disabled)

*/

/*
TODO:
	- do connections tls encrypted
*/

package main

import (
	"flag"
	"fmt"
	"github.com/nmcclain/ldap"
	"log"
	"os"
	"strconv"
)

const AppVersion = "0.0.1"

/*
---------------------------------
Version 0.0.1 20141219

Init version -- JA

*/

var (
	ldapserver string
	ldapport   string
	basedn     string
	Attributes []string = []string{"msNPAllowDialin", "userAccountControl"}
	queryuser  string
	passwd     string
	user       string
	filter     string
	rawoutput  bool
	account    string
	vpn        string
)

func init() {
	flag.StringVar(&user, "user", "testuser", "Username to query")
	flag.StringVar(&ldapserver, "ldaphost", "cell-dc-03", "Ldap server URL")
	flag.StringVar(&ldapport, "ldapport", "389", "Ldap Server PORT")
	flag.StringVar(&queryuser, "ldapuser", "cn=cellquery,cn=Users,dc=celluloidvfx,dc=inc", "User for authentification")
	flag.StringVar(&passwd, "ldappasswd", "cellquery123", "Password for authentification")
	flag.StringVar(&basedn, "ldapbase", "cn=Users,dc=celluloidvfx,dc=inc", "base DN for search")
	flag.BoolVar(&rawoutput, "raw", false, "Switch for displaying raw output")
}

func main() {

	version := flag.Bool("version", false, "prints current app version and exits")
	license := flag.Bool("license", false, "dumps the license and exits")

	greeter := "VPN Query " + AppVersion + " Copyright 2014 Celluloid VFX, Berlin and Johannes Amorosa"

	flag.Parse()

	if *version {
		fmt.Println(AppVersion)
		os.Exit(2)
	}

	log.Printf(greeter)

	// License flag was set
	if *license {
		printLicenseText()
		os.Exit(2)
	}

	// Need portnumber as int
	port, _ := strconv.Atoi(ldapport)

	// Dial
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapserver, port))
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
		os.Exit(2)
	}
	defer l.Close()
	// l.Debug = true

	// Bind
	err = l.Bind(queryuser, passwd)
	if err != nil {
		log.Printf("ERROR: Cannot bind: %s\n", err.Error())
		os.Exit(2)
	}
	// Set filter to user
	filter := "(cn=" + user + ")"

	// Build Search
	search := ldap.NewSearchRequest(
		basedn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		Attributes,
		nil)

	// Do search
	sr, err := l.Search(search)

	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
		os.Exit(2)
	}
	if len(sr.Entries) > 0 {
		// This display a "raw" output for debugging
		if rawoutput {
			log.Printf("Search: %s -> num of entries = %d\n", search.Filter, len(sr.Entries))
			sr.PrettyPrint(0)
			os.Exit(2)
		}

		// Renice data
		account = sr.Entries[0].GetAttributeValue("userAccountControl")
		vpn = sr.Entries[0].GetAttributeValue("msNPAllowDialin")

		// String compare
		if account == "512" && vpn == "TRUE" {
			log.Printf("VPN access for user " + user + " allowed")
			os.Exit(0)
		} else {
			log.Printf("VPN access for user " + user + " declined")
			os.Exit(1)
		}

	} else {
		// User doesn't exist or something funky happend
		log.Printf("VPN access for user " + user + " declined")
		os.Exit(2)
	}
}

func printLicenseText() {
	fmt.Println(`
Copyright 2014 Celluloid VFX, Berlin and Johannes Amorosa

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
`)
}
