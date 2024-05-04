package tunnel

type SSHTunnelConfig struct {
	SSHAddress        string
	SSHUser           string
	SSHPrivateKeyPath string
	LocalEndpoint     string
	RemoteEndpoint    string
}

// func main() {
// 	config := SSHTunnelConfig{
// 		SSHAddress:        "postgresHost:22",
// 		SSHUser:           "sshUser",
// 		SSHPrivateKeyPath: "key.pem",
// 		LocalEndpoint:     "localhost:0",
// 		RemoteEndpoint:    "remoteHost:5432",
// 	}

// 	listener, err := SetupSSHTunnel(config)
// 	if err != nil {
// 		log.Fatalf("Failed to start SSH tunnel: %v", err)
// 	}
// 	defer listener.Close()

// 	// Construct DSN for PostgreSQL connection
// 	localAddr := listener.Addr().(*net.TCPAddr)
// 	dsn := fmt.Sprintf("host=%s port=%d user=postgres password=123 dbname=db_name sslmode=disable", localAddr.IP, localAddr.Port)

// 	ctx := context.TODO()
// 	// Connect to PostgreSQL via SSH tunnel
// 	db, err := pgx.Connect(ctx, dsn)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
// 	}

// 	// Use db to perform SQL operations
// 	if err := db.Ping(ctx); err != nil {
// 		log.Fatalf("Failed to ping database: %v", err)
// 	}

// 	fmt.Println("Successfully connected to PostgreSQL through SSH tunnel!")
// }
