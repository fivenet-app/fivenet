package auth

import (
	"context"
	"net"

	"google.golang.org/grpc/credentials"
)

// Copied from google.golang.org/grpc/credentials/insecure
// Modified to return PrivacyAndIntegrity instead of NoSecurity as security level

// NewCredentials returns a TransportCredentials implementation which does
// not perform any authentication. It should be used only for testing or
// development purposes.
func NewCredentials() credentials.TransportCredentials {
	return insecureTC{}
}

// insecureTC implements the insecure transport credentials. The handshake
// methods simply return the passed in net.Conn and set the security level to
// NoSecurity.
type insecureTC struct{}

func (insecureTC) ClientHandshake(
	_ context.Context,
	_ string,
	conn net.Conn,
) (net.Conn, credentials.AuthInfo, error) {
	return conn, info{
		credentials.CommonAuthInfo{SecurityLevel: credentials.PrivacyAndIntegrity},
	}, nil
}

func (insecureTC) ServerHandshake(conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	return conn, info{
		credentials.CommonAuthInfo{SecurityLevel: credentials.PrivacyAndIntegrity},
	}, nil
}

func (insecureTC) Info() credentials.ProtocolInfo {
	return credentials.ProtocolInfo{SecurityProtocol: "secure"}
}

func (insecureTC) Clone() credentials.TransportCredentials {
	return insecureTC{}
}

func (insecureTC) OverrideServerName(string) error {
	return nil
}

// info contains the auth information for an insecure connection.
// It implements the AuthInfo interface.
type info struct {
	credentials.CommonAuthInfo
}

// AuthType returns the type of info as a string.
func (info) AuthType() string {
	return "insecure"
}

// ValidateAuthority allows any value to be overridden for the :authority
// header.
func (info) ValidateAuthority(string) error {
	return nil
}
