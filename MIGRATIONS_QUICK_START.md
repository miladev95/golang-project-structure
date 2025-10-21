# Migrations Quick Start

Fast reference for database migrations.

## 🚀 One-Minute Setup

```bash
# 1. Set environment variables
export DB_DRIVER=postgres
export DB_HOST=localhost
export DB_USER=postgres
export DB_NAME=myapp

# 2. Create database
createdb myapp

# 3. Start server (migrations run automatically)
go run cmd/server/main.go
```

**That's it!** Server starts and creates tables automatically. ✅

---

## ⚡ Commands

### Check Database

```bash
# List all tables
psql -U postgres -d myapp -c "\dt"

# Check users table
psql -U postgres -d myapp -c "SELECT * FROM users;"

# List indexes
psql -U postgres -d myapp -c "\di"

# Check table structure
psql -U postgres -d myapp -c "\d users"
```

### Run Migrations Manually

```bash
# Forward
psql -U postgres -d myapp -f migrations/001_create_users_table.up.sql

# Rollback
psql -U postgres -d myapp -f migrations/001_create_users_table.down.sql
```

### Test APIs

```bash
# Health check
curl http://localhost:8080/health

# Get all users
curl http://localhost:8080/api/v1/users

# Create user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

---

## 📁 Migration Files

| File | Purpose |
|------|---------|
| `migrations/001_create_users_table.up.sql` | Creates users table |
| `migrations/001_create_users_table.down.sql` | Drops users table |
| `migrations/README.md` | Detailed migration guide |
| `internal/config/migrations.go` | Migration runner |

---

## 🔄 Workflow

### Adding New Migration

1. **Create SQL files:**
   ```bash
   # Up migration
   vim migrations/002_add_password_column.up.sql
   
   # Down migration  
   vim migrations/002_add_password_column.down.sql
   ```

2. **Update `internal/config/migrations.go`:**
   ```go
   func RunMigrations(db *gorm.DB) error {
       // Add your migration code
   }
   ```

3. **Test:**
   ```bash
   go run cmd/server/main.go
   ```

4. **Verify:**
   ```bash
   psql -U postgres -d myapp -c "\d users"
   ```

---

## 🐳 Docker Quick Start

```bash
# Start PostgreSQL
docker run --name postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=myapp \
  -p 5432:5432 \
  -d postgres:15

# Wait a moment for startup
sleep 2

# Run server (auto-migrates)
export DB_HOST=localhost
go run cmd/server/main.go
```

---

## 📊 Current Schema

### Users Table

```sql
-- Primary table
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for fast email lookups
CREATE INDEX idx_users_email ON users(email);

-- Auto-update trigger
CREATE TRIGGER trigger_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_users_updated_at();
```

---

## ✅ Verify Installation

```bash
# Check database created
psql -U postgres -c "\l" | grep myapp

# Check tables exist
psql -U postgres -d myapp -c "\dt"

# Check indexes
psql -U postgres -d myapp -c "\di"

# Insert test data
psql -U postgres -d myapp -c \
  "INSERT INTO users (name, email) VALUES ('Test', 'test@example.com');"

# Verify insert
psql -U postgres -d myapp -c "SELECT * FROM users;"
```

---

## 🔧 Troubleshooting

| Problem | Solution |
|---------|----------|
| `FATAL: database does not exist` | Run `createdb myapp` |
| `permission denied` | Check DB_USER and DB_PASSWORD |
| `connection refused` | Check DB_HOST and DB_PORT |
| `already exists` | Add `IF NOT EXISTS` to SQL |
| Port already in use | Change DB_PORT or kill process |

---

## 📚 Full Documentation

See [docs/MIGRATIONS_GUIDE.md](docs/MIGRATIONS_GUIDE.md) for complete guide including:
- Multiple database setup (PostgreSQL, MySQL)
- Docker Compose setup
- golang-migrate tool integration
- Production deployment
- Troubleshooting

---

**Status:** ✅ Ready to use  
**Time to setup:** ~5 minutes  
**Database:** PostgreSQL or MySQL