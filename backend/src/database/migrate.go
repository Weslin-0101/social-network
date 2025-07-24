package database

import (
    "database/sql"
    "fmt"
    "os"
    "path/filepath"
    
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(db *sql.DB, migrationPath string) error {
    fmt.Printf("🔍 Checking migration path: %s\n", migrationPath)
    
    // Verificar se o diretório existe
    if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
        fmt.Printf("❌ Migration directory does not exist: %s\n", migrationPath)
        return fmt.Errorf("migration directory does not exist: %s", migrationPath)
    }
    
    // Listar arquivos na pasta
    files, err := filepath.Glob(filepath.Join(migrationPath, "*.sql"))
    if err != nil {
        fmt.Printf("❌ Error reading migration files: %v\n", err)
        return fmt.Errorf("error reading migration files: %w", err)
    }
    
    fmt.Printf("📁 Found %d migration files:\n", len(files))
    for _, file := range files {
        fmt.Printf("   - %s\n", filepath.Base(file))
    }
    
    if len(files) == 0 {
        fmt.Println("⚠️  No migration files found - skipping migrations")
        return nil
    }

    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        fmt.Printf("❌ Could not create migrate driver: %v\n", err)
        return fmt.Errorf("could not create migrate driver: %w", err)
    }

    sourceURL := fmt.Sprintf("file://%s", migrationPath)
    fmt.Printf("🚀 Migration source URL: %s\n", sourceURL)
    
    m, err := migrate.NewWithDatabaseInstance(
        sourceURL,
        "postgres",
        driver,
    )
    if err != nil {
        fmt.Printf("❌ Could not create migrate instance: %v\n", err)
        return fmt.Errorf("could not create migrate instance: %w", err)
    }

    // Verificar versão atual
    version, dirty, err := m.Version()
    if err != nil && err != migrate.ErrNilVersion {
        fmt.Printf("⚠️  Could not get current version: %v\n", err)
    } else {
        fmt.Printf("📊 Current migration version: %d (dirty: %v)\n", version, dirty)
    }

    fmt.Println("🔄 Running migrations...")
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        fmt.Printf("❌ Could not run migrations: %v\n", err)
        return fmt.Errorf("could not run migrations: %w", err)
    }
    
    // Verificar nova versão
    newVersion, _, _ := m.Version()
    fmt.Printf("✅ Migration completed successfully! New version: %d\n", newVersion)

    return nil
}
