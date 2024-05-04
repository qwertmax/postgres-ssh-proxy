package tunnel

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

type SSHTunnelConfig struct {
	SSHAddress        string
	SSHUser           string
	SSHPrivateKeyPath string
	LocalEndpoint     string
	RemoteEndpoint    string
}

func SetupSSHTunnel(config SSHTunnelConfig) (net.Listener, error) {
	key, err := os.ReadFile(config.SSHPrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: config.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", config.SSHAddress, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to SSH server: %v", err)
	}

	localListener, err := net.Listen("tcp", config.LocalEndpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to listen on local endpoint: %v", err)
	}

	go func() {
		defer sshClient.Close()
		for {
			localConn, err := localListener.Accept()
			if err != nil {
				log.Printf("Failed to accept local connection: %v", err)
				continue
			}

			remoteConn, err := sshClient.Dial("tcp", config.RemoteEndpoint)
			if err != nil {
				log.Printf("Failed to dial remote endpoint: %v", err)
				localConn.Close()
				continue
			}

			go func() {
				defer localConn.Close()
				defer remoteConn.Close()
				copyConn(localConn, remoteConn)
			}()
		}
	}()

	return localListener, nil
}

func copyConn(localConn, remoteConn net.Conn) {
	// Use io.Copy to duplex communication
	go func() { _, _ = io.Copy(localConn, remoteConn) }()
	_, _ = io.Copy(remoteConn, localConn)
}
