# Database Seeder

This seeder tool allows you to populate the database with dummy data for testing and development purposes. It can also handle database reconstruction when schema changes occur.

## Features

- 🔄 **Database Migrations**: Automatically runs migrations before seeding
- 🗑️ **Fresh Start**: Drop and recreate all tables/collections

## Usage

### Basic Seeding (with migrations)
Run migrations and seed data:
```bash
go run seeders/seeder.go
```

### Fresh Database Start
Drop all tables and collections, run migrations, then seed:
```bash
go run seeders/seeder.go --fresh
```

### Drop Tables Only
Drop tables but keep the database structure:
```bash
go run seeders/seeder.go --drop
```

### Seed Data Only
Skip migrations and only seed data (useful when tables already exist):
```bash
go run seeders/seeder.go --seed-only
```

## Command Line Flags

| Flag | Description |
|------|-------------|
| `--fresh` | Drop all tables/collections, run migrations, then seed data |
| `--drop` | Drop all tables before running migrations and seeding |
| `--seed-only` | Only seed data without running migrations |

## Prerequisites

Make sure you have:
- ✅ MySQL database running and configured
- ✅ MongoDB running and configured
- ✅ Redis running and configured
- ✅ `.env` file with correct database credentials

## Environment Variables

The seeder uses the same environment configuration as the main application. Ensure your `.env` file is properly configured in the `backend` directory.

## Notes

- The seeder is completely independent from the main application
- Safe to run multiple times (especially with `--fresh` flag)
- MongoDB indexes are created automatically

## Troubleshooting

If you encounter errors:

1. **Connection refused**: Ensure all databases (MySQL, MongoDB, Redis) are running
2. **Migration errors**: Try running with `--fresh` flag to start clean
3. **Duplicate key errors**: Run with `--fresh` to clear existing data

## Example Workflow

```bash
# Fresh start - recommended for development
go run seeders/seeder.go --fresh

# After schema changes in models
go run seeders/seeder.go --fresh

# Just add more data (tables exist)
go run seeders/seeder.go --seed-only

# Regular seed with migrations
go run seeders/seeder.go
```
