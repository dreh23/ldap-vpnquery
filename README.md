# ldap-vpnquery
This is a tool to query a Samba 4 Active Directory Server. We use it to see if certain flags are set so our VPN Frontend can decide to let users in or not.


    Active Directory Attributes:

    msNPAllowDialin: [TRUE]
    msNPAllowDialin: [FALSE]
    msNPAllowDialin: Not available

    userAccountControl: 512 (user enabled)
    userAccountControl: 514 (user disabled)


You need to install [http://golang.org/](go) and build the program. Once it's compiled use:

    vpnquery --help 
