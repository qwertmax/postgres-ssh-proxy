# postgres-ssh-proxy

example of how to use SSH proxy for Postgres connection

```golang
func main(){
	listener, err := tunnel.SetupSSHTunnel(tunnel.SSHTunnelConfig{
		SSHAddress:        "sshServer:22",
		SSHUser:           "user",
		SSHPrivateKeyPath: "key.pem",
		LocalEndpoint:     "localhost:0",
		RemoteEndpoint:    "remoteHost:5432",
	})
	if err != nil {
		log.Fatalf("Failed to start SSH tunnel: %v", err)
	}
	defer listener.Close()
	localAddr := listener.Addr().(*net.TCPAddr)

	ctx := context.TODO()
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", "pgUser", "pgPassword", localAddr.IP, localAddr.Port, "dbName")
	pgConn, err := pgx.Connect(ctx, connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pgConn.Close(ctx)
}
```