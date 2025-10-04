# Quick Setup Guide

This guide will help you get the Warehouse Management System up and running quickly.

## Prerequisites Check

Before you begin, make sure you have:

- [ ] Go 1.21+ installed (`go version`)
- [ ] Node.js 18+ installed (`node --version`)
- [ ] MySQL 8.0+ installed and running
- [ ] Git installed

## Step-by-Step Setup

### 1. Database Setup (5 minutes)

Open MySQL and create the database:

```sql
CREATE DATABASE warehouse_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

Verify it was created:

```sql
SHOW DATABASES;
```

### 2. Backend Setup (5 minutes)

```bash
# Navigate to backend
cd backend

# Copy environment file
cp env.example .env

# Edit .env with your MySQL credentials
# nano .env (Linux/Mac) or notepad .env (Windows)
# Update DB_PASSWORD with your MySQL password

# Install dependencies
go mod tidy
go mod download

# Run the backend
go run main.go
```

You should see:

```
Database connected successfully
Database migration completed
Default admin user created (username: admin, password: admin123)
Server starting on port 8080...
```

### 3. Frontend Setup (5 minutes)

Open a **NEW terminal** window:

```bash
# Navigate to frontend
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

You should see:

```
VITE v5.x.x  ready in xxx ms

➜  Local:   http://localhost:5173/
```

### 4. Access the Application

1. Open your browser
2. Go to: http://localhost:5173
3. Login with:
   - Username: `admin`
   - Password: `admin123`

## Verification Checklist

- [ ] Backend running on http://localhost:8080
- [ ] Frontend running on http://localhost:5173
- [ ] Can login to the application
- [ ] Dashboard shows statistics
- [ ] Can add a test product
- [ ] Can view, edit, and delete products

## Test the API

You can test the API using curl:

```bash
# Health check
curl http://localhost:8080/health

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

## Common Issues

### Backend won't start

**Error**: `Failed to connect to database`

**Solutions**:

1. Verify MySQL is running: `mysql --version`
2. Check MySQL credentials in `.env`
3. Verify database exists: `SHOW DATABASES;`

### Frontend shows connection error

**Solutions**:

1. Verify backend is running on port 8080
2. Check browser console for error messages
3. Clear browser cache and reload

### Port already in use

**Backend (8080)**:

```bash
# Linux/Mac
lsof -ti:8080 | xargs kill -9

# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

**Frontend (5173)**:

```bash
# Change port in vite.config.js
server: {
  port: 3000  // or any available port
}
```

## Next Steps

After successful setup:

1. Create some test products
2. Test filtering and search
3. Try exporting to CSV
4. Generate barcodes
5. Explore the dashboard

## Development Tips

### Hot Reload

Both frontend and backend support hot reload:

- **Frontend**: Automatically reloads on file changes
- **Backend**: Use `air` for hot reload (optional)

Install air for Go:

```bash
go install github.com/cosmtrek/air@latest
cd backend
air
```

### Database Management

To reset the database:

```sql
DROP DATABASE warehouse_db;
CREATE DATABASE warehouse_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

Then restart the backend to run migrations again.

### API Testing

Use Postman or Insomnia to test the API:

1. Import the endpoints from README.md
2. Set up environment variables for the token
3. Test all CRUD operations

## Production Deployment

For production deployment:

1. **Backend**:

   ```bash
   cd backend
   go build -o warehouse-api
   ./warehouse-api
   ```

2. **Frontend**:

   ```bash
   cd frontend
   npm run build
   # Serve the dist/ folder with nginx or similar
   ```

3. Update environment variables for production
4. Use a process manager (PM2, systemd)
5. Set up reverse proxy (nginx)
6. Enable HTTPS

## Need Help?

- Check the main README.md for detailed documentation
- Review API documentation for endpoint details
- Check GitHub issues for known problems

---

Happy coding! 🚀
