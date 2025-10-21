# Database Migrations Guide

Complete guide to managing database migrations in this Go project.

## 📚 Quick Start

### Option 1: Automatic Migrations (Development) ✅ **Current Setup**

Migrations run automatically when the server starts:

```bash
go run cmd/server/main.go
# 🔄 Running database migrations...
# ✅ Created users table
# ✅ Created index on users.email
# ✅ All migrations completed successfully
```

**Pros:**
- ✅ No additional tools needed
- ✅ Simple for development
- ✅ Idempotent (safe to run multiple times)

**Cons:**
- ❌ Harder to track versions in production
- ❌ Less control over rollbacks

### Option 2: SQL Migration Files (Production) 🏗️

Manual SQL files in `migrations/` directory for production use:

```bash
# PostgreSQL
psql -U postgres -d myapp -f migrations/001_create_users_table.up.sql

# MySQL
mysql -u root -p myapp < migrations/001_create_users_table.up.sql

# Rollback
psql -U postgres -d myapp -f migrations/001_create_users_table.down.sql
```

### Option 3: golang-migrate Tool (Recommended for Production) 🚀

Professional migration management with version tracking:

```bash
# Install
go install -tags 'postgres,mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path migrations/ -database "postgres://user:pass@localhost/myapp?sslmode=disable" up

# Check version
migrate -path migrations/ -database "postgres://user:pass@localhost/myapp?sslmode=disable" version

# Rollback last migration
migrate -path migrations/ -database "postgres://user:pass@localhost/myapp?sslmode=disable" down 1
```

---

## 🗂️ File Structure

```
migrations/
├── 001_create_users_table.up.sql      # Forward migration
├── 001_create_users_table.down.sql    # Rollback migration
├── README.md                           # Migration guide
└── [Future migrations...]
```

**Naming Convention:** `{VERSION}_{DESCRIPTION}.{DIRECTION}.sql`

- `VERSION`: 3-digit number (001, 002, 003...)
- `DIRECTION`: `up` (forward) or `down` (rollback)

---

## 🔄 Current Migrations

### Migration 001: Create Users Table

**File:** `001_create_users_table.up.sql`

Creates the users table with:
- `id`: BIGSERIAL primary key
- `name`: VARCHAR(255) not null
- `email`: VARCHAR(255) unique
- `created_at`: Auto-timestamp
- `updated_at`: Auto-updated on record changes
- Auto-update trigger for `updated_at`
- Index on `email` for fast lookups

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);

-- Auto-update trigger
CREATE TRIGGER trigger_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_users_updated_at();
```

---

## 🛠️ Working with Migrations

### Running Migrations Programmatically

```go
import "github.com/miladev95/golang-project-structure/internal/config"

// In your application startup
db, err := config.NewDatabase(cfg)
if err != nil {
    log.Fatal(err)
}

// Run all pending migrations
if err := config.RunMigrations(db); err != nil {
    log.Fatal(err)
}
```

### Checking Migration Status

```go
import "github.com/miladev95/golang-project-structure/internal/config"

db, err := config.NewDatabase(cfg)
status := config.CheckMigrationStatus(db)

// Output:
// map[string]bool{
//     "users_table": true,
//     "users_email_index": true,
// }
```

### Rolling Back Migrations

```go
// ⚠️ WARNING: This will delete all data!
if err := config.RollbackMigrations(db); err != nil {
    log.Fatal(err)
}
```

---

## 📝 Adding New Migrations

### Step 1: Create Migration Files

Create two files in `migrations/`:

**002_add_password_to_users.up.sql:**
```sql
ALTER TABLE users ADD COLUMN password VARCHAR(255);
```

**002_add_password_to_users.down.sql:**
```sql
ALTER TABLE users DROP COLUMN password;
```

### Step 2: Update Migration Handler (if using programmatic approach)

Edit `internal/config/migrations.go`:

```go
func RunMigrations(db *gorm.DB) error {
    // ... existing code ...

    // Add new password column
    if !db.Migrator().HasColumn("users", "password") {
        if err := db.Migrator().AddColumn("users", "password VARCHAR(255)"); err != nil {
            return fmt.Errorf("failed to add password column: %w", err)
        }
        log.Println("✅ Added password column to users")
    }

    return nil
}
```

### Step 3: Test

```bash
# Start server - runs migrations automatically
go run cmd/server/main.go

# Verify in database
psql -U postgres -d myapp -c "SELECT * FROM users;"
```

---

## 🗄️ Database Setup

### PostgreSQL

```bash
# Create database
createdb myapp

# Set environment
export DB_DRIVER=postgres
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=yourpassword
export DB_NAME=myapp

# Start server (runs migrations automatically)
go run cmd/server/main.go
```

### MySQL

```bash
# Create database
mysql -u root -p -e "CREATE DATABASE myapp;"

# Set environment
export DB_DRIVER=mysql
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=yourpassword
export DB_NAME=myapp

# Start server (runs migrations automatically)
go run cmd/server/main.go
```

### Docker (PostgreSQL)

```bash
# Start PostgreSQL container
docker run --name postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=myapp \
  -p 5432:5432 \
  -d postgres:15

# Set environment
export DB_DRIVER=postgres
export DB_HOST=localhost
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=myapp

# Start server
go run cmd/server/main.go
```

---

## 🧪 Testing Migrations

### Test Forward Migration

```bash
# Start fresh
dropdb myapp 2>/dev/null || true
createdb myapp

# Run forward
psql -U postgres -d myapp -f migrations/001_create_users_table.up.sql

# Verify
psql -U postgres -d myapp -c "\dt"  # List tables
psql -U postgres -d myapp -c "\di"  # List indexes
```

### Test Rollback

```bash
# Run rollback
psql -U postgres -d myapp -f migrations/001_create_users_table.down.sql

# Verify table is gone
psql -U postgres -d myapp -c "\dt"
```

### Test Both Directions

```bash
# Forward
psql -U postgres -d myapp -f migrations/001_create_users_table.up.sql

# Verify created
psql -U postgres -d myapp -c "\dt"

# Backward
psql -U postgres -d myapp -f migrations/001_create_users_table.down.sql

# Verify dropped
psql -U postgres -d myapp -c "\dt"

# Forward again
psql -U postgres -d myapp -f migrations/001_create_users_table.up.sql

# Verify recreated
psql -U postgres -d myapp -c "\dt"
```

---

## 🐳 Docker Compose Setup

Create `docker-compose.yml` for easy local development:

```yaml
version: '3.8'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: myapp
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  app:
    build: .
    environment:
      DB_DRIVER: postgres
      DB_HOST: postgres
      DB_USER: postgres
      DB_PASSWORD: password
      DB_NAME: myapp
      SERVER_PORT: 8080
    ports:
      - "8080:8080"
    depends_on:
      - postgres

volumes:
  postgres_data:
```

Run with:
```bash
docker-compose up
```

---

## 📋 Checklist for Production

- [ ] All migrations have both `.up.sql` and `.down.sql` files
- [ ] Tested forward migration on fresh database
- [ ] Tested rollback migration
- [ ] Verified data integrity after rollback
- [ ] Used constraints (PRIMARY KEY, UNIQUE, NOT NULL)
- [ ] Created indexes for frequently queried columns
- [ ] Added comments for complex migrations
- [ ] Documented any breaking changes
- [ ] Backed up database before applying
- [ ] Tested in staging environment first
- [ ] Ran migrations outside of peak hours
- [ ] Monitored application after deployment
- [ ] Saved rollback plan before deploying

---

## 🚨 Troubleshooting

### Migration Fails with "Column Already Exists"

**Problem:** Running migration twice causes error

**Solution:** Add `IF NOT EXISTS` or `IF NOT` clauses:

```sql
-- ✅ Good
ALTER TABLE users ADD COLUMN IF NOT EXISTS password VARCHAR(255);

-- ❌ Bad
ALTER TABLE users ADD COLUMN password VARCHAR(255);
```

### Migration Fails with "Table Doesn't Exist"

**Problem:** Running rollback first or wrong order

**Solution:** Check migration files exist and database exists:

```bash
# Check tables
psql -U postgres -d myapp -c "\dt"

# Recreate database if needed
dropdb myapp
createdb myapp

# Run migrations fresh
psql -U postgres -d myapp -f migrations/001_create_users_table.up.sql
```

### Migration Lock During Deployment

**Problem:** Migration locks the table and blocks application

**Solution:** Run migrations before deploying new code:

1. Apply migrations
2. Verify success
3. Deploy new application code

### Rollback Too Slow on Large Tables

**Problem:** Dropping/recreating large tables takes long time

**Solution:** Add `CONCURRENTLY` (PostgreSQL) or plan for downtime:

```sql
-- PostgreSQL: Create index without blocking writes
CREATE INDEX CONCURRENTLY idx_column ON table_name(column);

-- For large operations, plan maintenance window
```

---

## 🔗 Related Documentation

- [README.md](../README.md) - Project overview
- [DI_ARCHITECTURE.md](./DI_ARCHITECTURE.md) - Dependency injection
- [migrations/README.md](../migrations/README.md) - Migration files reference
- [.env.example](../.env.example) - Environment configuration

---

## 📚 External Resources

- **golang-migrate**: https://github.com/golang-migrate/migrate
- **GORM Migrations**: https://gorm.io/docs/migration.html
- **PostgreSQL Docs**: https://www.postgresql.org/docs/current/sql-syntax.html
- **MySQL Docs**: https://dev.mysql.com/doc/

---

**Last Updated:** 2024
**Status:** ✅ Production Ready