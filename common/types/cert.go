package types

//CertPolicy - certificate validation policy
//Possible values are: none,required,verify
type CertPolicy string

const (
	//CertPolicyNone - Certificate is not required
	CertPolicyNone CertPolicy = "none"
	//CertPolicyRequired - Certificate is required but if it is invalid is still accepted
	CertPolicyRequired CertPolicy = "required"
	//CertPolicyVerify - Certificate is required and checked if it is valid
	CertPolicyVerify CertPolicy = "verify"
)

/*ConnectionSecurityLevel - general system-wide,security policy.
Both sides of connection should set at least the same level, wherein
actual level is determinated by value of a server-side configuration.

There are three levels that can be set:

none - means that connection is not secured

serveronly - on server side, means that only server uses certificate to secure connection however,
client should checks if server sends a valid certificate

clientandserver - both client and server use certificates to establish connection

*/
type ConnectionSecurityLevel string

const (
	//ConnectionSecurityLevelNone  - do not use certificates on both sides
	ConnectionSecurityLevelNone ConnectionSecurityLevel = "none"
	//ConnectionSecurityLevelServeOnly - only server sends certificate
	ConnectionSecurityLevelServeOnly ConnectionSecurityLevel = "serveronly"
	//ConnectionSecurityLevelClientAndServer - both client and server use certificates
	ConnectionSecurityLevelClientAndServer ConnectionSecurityLevel = "clientandserver"
)
