// +build !windows

package main

// the idea here is to use something like https://github.com/apcera/gssapi/blob/master/test/client_access_test.go
// to get a kerberos ticket (given that the current session already has a TGT)
type CurrentOsGssImplementation struct {
}

// GetTicket returns a ticket for the given spn, using the credentials of the current user
func (t CurrentOsGssImplementation) GetTicket(spn string) []byte {
    return []byte{}
}