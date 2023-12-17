package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const (
	successfulAuthResponse = "auth_ok:1\nuid:%d\ngid:%d\ndir:%s\nend\n"
	failedAuthNotFound     = "auth_ok:0\nend\n"
	failedAuthFatalError   = "auth_ok:-1\nend\n"
)

var debugMode = false

func debugPrint(message string) {
	if debugMode {
		fmt.Println("DEBUG:", message)
	}
}

func main() {
	// Read user credentials from environment variables
	username := os.Getenv("AUTHD_ACCOUNT")
	password := os.Getenv("AUTHD_PASSWORD")

	if username == "" || password == "" {
		fmt.Print(failedAuthFatalError)
		os.Exit(1)
	}

	// Read DEBUG environment variable
	debugEnv := os.Getenv("DEBUG")
	if debugEnv == "true" {
		debugMode = true
	}

	// Read database configuration from environment variables
	dbHost, dbPort, dbName, dbUser, dbPassword, dbEngine := getDatabaseConfig()
	if dbHost == "" || dbPort == "" || dbName == "" || dbUser == "" || dbPassword == "" || dbEngine == "" {
		debugPrint("Database configuration is incomplete.")
		fmt.Print(failedAuthFatalError)
		os.Exit(1)
	}

	// Open database connection
	db, err := openDatabaseConnection(dbEngine, dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		debugPrint(fmt.Sprintf("Failed to open database connection: %v", err))
		fmt.Print(failedAuthFatalError)
		os.Exit(1)
	}
	defer db.Close()

	// Query the database for the stored credentials
	query, queryArgs := getDatabaseQuery(dbEngine, username)
	rows, err := db.Query(query, queryArgs...)
	if err != nil {
		debugPrint(fmt.Sprintf("Failed to execute database query: %v", err))
		fmt.Print(failedAuthFatalError)
		os.Exit(1)
	}
	defer rows.Close()

	// Check if the user exists and verify the password
	if rows.Next() {
		var dbUsername, dbPassword string
		err := rows.Scan(&dbUsername, &dbPassword)
		if err != nil {
			handleError("Failed to scan database row", failedAuthFatalError)
		}

		// Extract parameters from the stored hash
		parts := strings.Split(dbPassword, "$")

		// Extract iterations from the stored hash
		iterations, err := strconv.Atoi(parts[1])
		if err != nil {
			handleError("Failed to convert iterations to integer", failedAuthFatalError)
		}

		salt := parts[2]
		hashedPassword := parts[3]

		// Hash the entered password using PBKDF2
		hashedInput := pbkdf2.Key([]byte(password), []byte(salt), iterations, sha256.Size, sha256.New)

		// Encode the hashed input to base64 for comparison
		encodedHashedInput := base64.StdEncoding.EncodeToString(hashedInput)

		// Compare the generated hash with the stored hash
		if encodedHashedInput == hashedPassword {
			// Get the current user's UID and GID
			uid, gid := getCurrentUserIDs()

			// Get the consumption directory from the environment variable
			consumptionDir := os.Getenv("PAPERLESS_CONSUMPTION_DIR")
			if consumptionDir == "" {
				handleError("PAPERLESS_CONSUMPTION_DIR is not set", failedAuthFatalError)
			}

			// Print the successful authentication response
			fmt.Printf(successfulAuthResponse, uid, gid, consumptionDir)
			os.Exit(0)
		} else {
			handleError("Password verification failed", failedAuthFatalError)
		}
	} else {
		handleError("User not found in the database", failedAuthFatalError)
	}
}

func getDatabaseConfig() (string, string, string, string, string, string) {
	return os.Getenv("PAPERLESS_DBHOST"), os.Getenv("PAPERLESS_DBPORT"),
		os.Getenv("PAPERLESS_DBNAME"), os.Getenv("DB_USER"),
		os.Getenv("PAPERLESS_DBPASS"), os.Getenv("PAPERLESS_DBENGINE")
}

func openDatabaseConnection(engine, host, port, user, password, name string) (*sql.DB, error) {
	var connStr string
	switch engine {
	case "postgres", "postgresql":
		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)
		engine = "postgres"
	case "mysql", "mariadb":
		connStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, name)
	default:
		debugPrint(fmt.Sprintf("Unsupported database engine: %s", engine))
		fmt.Print(failedAuthFatalError)
		os.Exit(1)
	}

	return sql.Open(engine, connStr)
}

func getDatabaseQuery(engine, username string) (string, []interface{}) {
	var query string
	var queryArgs []interface{}

	switch engine {
	case "postgres", "postgresql":
		query = "SELECT username, password FROM auth_user WHERE username = $1"
		queryArgs = []interface{}{username}
	case "mysql", "mariadb":
		query = "SELECT username, password FROM auth_user WHERE username = ?"
		queryArgs = []interface{}{username}
	default:
		debugPrint(fmt.Sprintf("Unsupported database engine: %s", engine))
		fmt.Print(failedAuthFatalError)
		os.Exit(1)
	}

	return query, queryArgs
}

func getCurrentUserIDs() (int, int) {
	currentUser, err := user.Current()
	if err != nil {
		handleError("Failed to get current user information", failedAuthFatalError)
	}

	uid, err := strconv.Atoi(currentUser.Uid)
	if err != nil {
		handleError("Failed to convert UID to integer", failedAuthFatalError)
	}

	gid, err := strconv.Atoi(currentUser.Gid)
	if err != nil {
		handleError("Failed to convert GID to integer", failedAuthFatalError)
	}

	return uid, gid
}

func handleError(errorMessage, exitMessage string) {
	debugPrint(errorMessage)
	fmt.Print(exitMessage)
	os.Exit(1)
}
