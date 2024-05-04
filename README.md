# postgres-ssh-proxy

```golang
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/ssh"
)

type SSHTunnelConfig struct {
	SSHAddress        string
	SSHUser           string
	SSHPrivateKeyPath string
	LocalEndpoint     string
	RemoteEndpoint    string
}

func main() {
	config := SSHTunnelConfig{
		SSHAddress:        "postgresHost:22",
		SSHUser:           "sshUser",
		SSHPrivateKeyPath: "key.pem",
		LocalEndpoint:     "localhost:0",
		RemoteEndpoint:    "remoteHost:5432",
	}

	listener, err := setupSSHTunnel(config)
	if err != nil {
		log.Fatalf("Failed to start SSH tunnel: %v", err)
	}
	defer listener.Close()

	// Construct DSN for PostgreSQL connection
	localAddr := listener.Addr().(*net.TCPAddr)
	dsn := fmt.Sprintf("host=%s port=%d user=postgres password=123 dbname=db_name sslmode=disable", localAddr.IP, localAddr.Port)

	ctx := context.TODO()
	// Connect to PostgreSQL via SSH tunnel
	db, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Use db to perform SQL operations
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL through SSH tunnel!")
}

func setupSSHTunnel(config SSHTunnelConfig) (net.Listener, error) {
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
```