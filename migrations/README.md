# Database Migrations

This directory contains SQL migration files for database schema management.

## Naming Convention

Migrations follow the standard pattern:
```
{VERSION}_{DESCRIPTION}.{DIRECTION}.sql
```

- **VERSION**: 3-digit number (001, 002, etc.)
- **DESCRIPTION**: What the migration does (snake_case)
- **DIRECTION**: Either `up` (forward) or `down` (rollback)

## Examples

- `001_create_users_table.up.sql` - Creates users table
- `001_create_users_table.down.sql` - Drops users table
- `002_add_password_to_users.up.sql` - Adds password column
- `002_add_password_to_users.down.sql` - Removes password column

## Running Migrations

### Using a Migration Tool

For production, use a Go migration tool:

**Option 1: golang-migrate/migrate**
```bash
# Install
go install -tags 'postgres,mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run up
migrate -path migrations/ -database "postgres://user:pass@localhost/db?sslmode=disable" up

# Run down
migrate -path migrations/ -database "postgres://user:pass@localhost/db?sslmode=disable" down 1

# Check current version
migrate -path migrations/ -database "postgres://user:pass@localhost/db?sslmode=disable" version
```

**Option 2: GORM Auto-Migration (Development)**
```go
// In repository initialization
func (r *userRepository) Init() error {
    return r.db.AutoMigrate(&models.User{})
}
```

### Manual Execution (MySQL)

```bash
# Create database
mysql -u root -p -e "CREATE DATABASE myapp;"

# Run migration
mysql -u root -p myapp < migrations/001_create_users_table.up.sql

# Rollback
mysql -u root -p myapp < migrations/001_create_users_table.down.sql
```

### Manual Execution (PostgreSQL)

```bash
# Create database
createdb myapp

# Run migration
psql -U postgres -d myapp -f migrations/001_create_users_table.up.sql

# Rollback
psql -U postgres -d myapp -f migrations/001_create_users_table.down.sql
```

## Current Migrations

| # | Description | Status |
|---|-------------|--------|
| 001 | Create users table | ✅ Initial |

## Best Practices

1. **Always provide down migration** - Makes rollbacks safe
2. **Test both directions** - Run up, then down, then up again
3. **Keep migrations small** - One logical change per migration
4. **Version sequentially** - Don't skip numbers
5. **Document complex migrations** - Add comments in SQL
6. **Test in development** - Before applying to production
7. **Backup before running** - On production databases

## Migration Strategies

### PostgreSQL with golang-migrate (Recommended)

1. Install golang-migrate
2. Set database URL: `DATABASE_URL=postgres://user:pass@localhost/db`
3. Run: `migrate -path migrations/ -database $DATABASE_URL up`

### GORM AutoMigration (Quick Development)

```go
// In repository or initialization code
type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    // Auto-migrate on init
    db.AutoMigrate(&models.User{})
    return &UserRepository{db: db}
}
```

### Docker + PostgreSQL

```bash
# Run migration in Docker
docker exec -it postgres psql -U postgres -d myapp -f migrations/001_create_users_table.up.sql
```

## Troubleshooting

**Q: Migration already exists error?**
A: The migration tool tracks applied migrations. Recreate the database if needed.

**Q: Column already exists?**
A: Use `IF NOT EXISTS` for safety, or manually check and adjust migration.

**Q: Constraint violation?**
A: Ensure down migration removes all dependent objects (indexes, triggers, functions).

---

**For detailed setup instructions, see [DI_ARCHITECTURE.md](../docs/DI_ARCHITECTURE.md) or [README.md](../README.md)**