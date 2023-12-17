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
	dbHost := os.Getenv("PAPERLESS_DBHOST")
	dbPort := os.Getenv("PAPERLESS_DBPORT")
	dbName := os.Getenv("PAPERLESS_DBNAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("PAPERLESS_DBPASS")
	dbEngine := os.Getenv("PAPERLESS_DBENGINE")

	if dbHost == "" || dbPort == "" || dbName == "" || dbUser == "" || dbPassword == "" || dbEngine == "" {
		debugPrint("Database configuration is incomplete.")
		fmt.Print(failedAuthFatalError)
		os.Exit(1)
	}

	var connStr string
	var db *sql.DB
	var err error

	switch dbEngine {
	case "postgres":
		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName)
		db, err = sql.Open("postgres", connStr)
	case "mysql", "mariadb":
		connStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
		db, err = sql.Open("mysql", connStr)
	default:
		debugPrint(fmt.Sprintf("Unsupported database engine: %s", dbEngine))
		fmt.Print(failedAuthFatalError)
		os.Exit(1)
	}

	if err != nil {
		debugPrint(fmt.Sprintf("Failed to open database connection: %v", err))
		fmt.Print(failedAuthFatalError)
		os.Exit(1)
	}
	defer db.Close()

	// Query the database for the stored credentials
	var query string
	var queryArgs []interface{}

	switch dbEngine {
	case "postgres":
		query = "SELECT username, password FROM auth_user WHERE username = $1"
		queryArgs = []interface{}{username}
	case "mysql", "mariadb":
		query = "SELECT username, password FROM auth_user WHERE username = ?"
		queryArgs = []interface{}{username}
	default:
		debugPrint(fmt.Sprintf("Unsupported database engine: %s", dbEngine))
		fmt.Print(failedAuthFatalError)
		os.Exit(1) // Exit with code 1 for fatal error
	}

	rows, err := db.Query(query, queryArgs...)
	if err != nil {
		debugPrint(fmt.Sprintf("Failed to execute database query: %v", err))
		fmt.Print(failedAuthFatalError)
		os.Exit(1) // Exit with code 1 for fatal error
	}
	defer rows.Close()

	// Check if the user exists and verify the password
	if rows.Next() {
		var dbUsername, dbPassword string
		err := rows.Scan(&dbUsername, &dbPassword)
		if err != nil {
			debugPrint(fmt.Sprintf("Failed to scan database row: %v", err))
			fmt.Print(failedAuthFatalError)
			os.Exit(1)
		}

		// Extract parameters from the stored hash
		parts := strings.Split(dbPassword, "$")

		// Extract iterations from the stored hash
		iterations, err := strconv.Atoi(parts[1])
		if err != nil {
			debugPrint(fmt.Sprintf("Failed to convert iterations to integer: %v", err))
			fmt.Print(failedAuthFatalError)
			os.Exit(1)
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
			currentUser, err := user.Current()
			if err != nil {
				debugPrint(fmt.Sprintf("Failed to get current user information: %v", err))
				fmt.Print(failedAuthFatalError)
				os.Exit(1)
			}

			uid, err := strconv.Atoi(currentUser.Uid)
			if err != nil {
				debugPrint(fmt.Sprintf("Failed to convert UID to integer: %v", err))
				fmt.Print(failedAuthFatalError)
				os.Exit(1)
			}

			gid, err := strconv.Atoi(currentUser.Gid)
			if err != nil {
				debugPrint(fmt.Sprintf("Failed to convert GID to integer: %v", err))
				fmt.Print(failedAuthFatalError)
				os.Exit(1)
			}

			// Get the consumption directory from the environment variable
			consumptionDir := os.Getenv("PAPERLESS_CONSUMPTION_DIR")
			if consumptionDir == "" {
				debugPrint("PAPERLESS_CONSUMPTION_DIR is not set.")
				fmt.Print(failedAuthFatalError)
				os.Exit(1)
			}

			// Print the successful authentication response
			fmt.Printf(successfulAuthResponse, uid, gid, consumptionDir)
			os.Exit(0)
		} else {
			debugPrint("Password verification failed.")
			os.Exit(1)
		}
	} else {
		debugPrint("User not found in the database.")
		os.Exit(1)
	}
}
