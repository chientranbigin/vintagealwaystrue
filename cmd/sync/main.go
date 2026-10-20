package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "24jHtiebmCYdexH.root:SaleDbPassword@2026@tcp(gateway01.ap-southeast-1.prod.aws.tidbcloud.com:4000)/sale_db?tls=true&multiStatements=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	// Use single connection to preserve session variables
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		log.Fatalf("Failed to get dedicated conn: %v", err)
	}
	defer conn.Close()

	fmt.Println("Disabling foreign key checks and dropping all existing tables in TiDB sale_db...")
	_, err = conn.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 0;")
	if err != nil {
		log.Printf("Warning setting FK checks: %v", err)
	}

	rows, err := conn.QueryContext(ctx, "SHOW TABLES;")
	if err != nil {
		log.Fatalf("Failed to show tables: %v", err)
	}
	var tables []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err == nil {
			tables = append(tables, t)
		}
	}
	rows.Close()

	for _, t := range tables {
		fmt.Printf("Dropping table `%s`...\n", t)
		if _, err := conn.ExecContext(ctx, fmt.Sprintf("DROP TABLE IF EXISTS `%s`;", t)); err != nil {
			log.Printf("Warning dropping %s: %v", t, err)
		}
	}

	dumpPath := "/Applications/MAMP/htdocs/sale/database/backups/fresh_db.sql"
	file, err := os.Open(dumpPath)
	if err != nil {
		log.Fatalf("Failed to open dump: %v", err)
	}
	defer file.Close()

	fmt.Println("Executing fresh_db.sql into TiDB sale_db...")
	scanner := bufio.NewScanner(file)
	// increase buffer capacity for large INSERT statements
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 64*1024*1024)

	var currentStmt strings.Builder
	executed := 0

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "--") || strings.HasPrefix(trimmed, "/*") || trimmed == "" {
			continue
		}

		currentStmt.WriteString(line)
		currentStmt.WriteString("\n")

		if strings.HasSuffix(trimmed, ";") {
			stmt := strings.TrimSpace(currentStmt.String())
			currentStmt.Reset()

			if strings.HasPrefix(strings.ToUpper(stmt), "LOCK TABLES") || strings.HasPrefix(strings.ToUpper(stmt), "UNLOCK TABLES") {
				continue
			}

			if _, err := conn.ExecContext(ctx, stmt); err != nil {
				log.Printf("Error executing stmt (first 100 chars: %s): %v", stmt[:min(100, len(stmt))], err)
			} else {
				executed++
				if executed%10 == 0 {
					fmt.Printf("Executed %d statements...\n", executed)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner error: %v", err)
	}

	fmt.Printf("Finished! Executed %d statements successfully.\n", executed)

	// Verify table counts
	fmt.Println("\n--- Final Row Counts in TiDB Cloud sale_db ---")
	rows, err = conn.QueryContext(ctx, "SHOW TABLES;")
	if err != nil {
		log.Fatalf("Failed to show tables: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err == nil {
			var count int
			if err := conn.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM `%s`", t)).Scan(&count); err == nil {
				fmt.Printf(" - %s: %d rows\n", t, count)
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
