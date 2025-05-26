package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenbupe/recipe-video-parser/internal/shared/platform/storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nCancelando operación...")
		cancel()
		os.Exit(1)
	}()

	if len(os.Args) < 3 {
		fmt.Println("Se requiere un comando y el nombre de la base de datos: db-create <ruta/nombre.db>")
		fmt.Println("Uso: dev db-create <ruta/nombre.db>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "db-create":
		dbCreateCmd(ctx, os.Args[2])
	default:
		fmt.Printf("Comando desconocido: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func dbCreateCmd(ctx context.Context, database string) {
	// Quitar extensión .db si la tiene para usar storage.CreateConnection
	if len(database) > 3 && database[len(database)-3:] == ".db" {
		database = database[:len(database)-3]
	}

	dbFile := database + ".db"
	dir := filepath.Dir(database)
	if dir != "." && dir != "" {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("Error creando el directorio para la base de datos: %v\n", err)
			os.Exit(1)
		}
	}

	// Verificar si la base de datos ya existe
	if _, err := os.Stat(dbFile); err == nil {
		fmt.Printf("Error: la base de datos '%s' ya existe.\n", dbFile)
		os.Exit(1)
	}

	cfg := &storage.Dbconfig{Database: database}
	conn, err := storage.CreateConnection("cli-dev", cfg)
	if err != nil {
		fmt.Printf("Error abriendo la base de datos: %v\n", err)
		os.Exit(1)
	}
	defer conn.Db.Close()

	// Leer el esquema desde el archivo sql/schema.sql
	schemaBytes, err := os.ReadFile("sql/schema.sql")
	if err != nil {
		fmt.Printf("Error leyendo el archivo de esquema: %v\n", err)
		os.Exit(1)
	}

	_, err = conn.Db.ExecContext(ctx, string(schemaBytes))
	if err != nil {
		fmt.Printf("Error creando la base de datos: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Base de datos SQLite inicializada en %s\n", dbFile)
}
