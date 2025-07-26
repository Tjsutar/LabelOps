# LabelOps Setup Guide

This guide will help you set up the complete LabelOps system with Angular 18 frontend, Go backend, and PostgreSQL database.

## Prerequisites

- **Node.js 18+** and npm
- **Go 1.21+**
- **PostgreSQL 14+**
- **Git**

## Quick Start

### 1. Clone and Setup

```bash
# Clone the repository
git clone <your-repo-url>
cd LabelOps

# Copy environment files
cp backend/env.example backend/.env
```

### 2. Database Setup

```bash
# Create PostgreSQL database
createdb labelops

# Update backend/.env with your database credentials
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=labelops
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

### 3. Backend Setup

```bash
cd backend

# Install Go dependencies
go mod tidy

# Run the backend
go run main.go
```

The backend will automatically:
- Connect to PostgreSQL
- Create all necessary tables
- Set up stored procedures
- Start the API server on port 8080

### 4. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start the development server
npm start
```

The frontend will be available at `http://localhost:4200`

## API Endpoints

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration

### Labels (Protected)
- `POST /api/v1/labels/batch` - Process TMT Bar batch
- `GET /api/v1/labels` - Get labels with filters
- `GET /api/v1/labels/:id` - Get specific label
- `POST /api/v1/labels/:id/print` - Print label
- `GET /api/v1/labels/export/csv` - Export labels as CSV

### Print Jobs (Protected)
- `GET /api/v1/print-jobs` - Get print jobs
- `GET /api/v1/print-jobs/:id` - Get specific print job
- `POST /api/v1/print-jobs/:id/retry` - Retry failed print job

### Audit Logs (Protected)
- `GET /api/v1/audit-logs` - Get audit logs
- `GET /api/v1/audit-logs/export/csv` - Export audit logs as CSV

### Admin (Admin Only)
- `GET /api/v1/admin/users` - Get all users
- `POST /api/v1/admin/users` - Create user
- `PUT /api/v1/admin/users/:id` - Update user
- `DELETE /api/v1/admin/users/:id` - Delete user
- `GET /api/v1/admin/stats` - Get system statistics

## Sample TMT Bar Data

The system is designed to handle TMT Bar data in this format:

```json
{
  "LOCATION": null,
  "BUNDLE_NOS": 6,
  "PQD": "101520002123005267",
  "UNIT": "SAIL-BSP",
  "TIME1": "14:21",
  "LENGTH": "STD",
  "HEAT_NO": "C075400",
  "PRODUCT_HEADING": "TMT BAR",
  "ISI_BOTTOM": "CML 187244",
  "ISI_TOP": "IS 1786:2008",
  "CHARGE_DTM": "202403041421",
  "MILL": "MM",
  "GRADE": "IS 1786 FE550D",
  "URL_APIKEY": "4700b47719c54d979af7f7a2d1dc67db",
  "ID": null,
  "WEIGHT": null,
  "SECTION": "TMT BAR 25",
  "DATE1": "04-MAR-24"
}
```

## Features

### ✅ Implemented
- **Authentication**: JWT-based login/logout with role-based access
- **TMT Bar Processing**: Batch processing with duplicate detection
- **Label Management**: CRUD operations for TMT Bar labels
- **Print Integration**: ZPL generation and print job management
- **Audit Logging**: Complete audit trail for all actions
- **CSV Export**: Export labels and audit logs
- **Responsive UI**: Modern Angular 18 + Tailwind CSS interface
- **Database**: PostgreSQL with stored procedures for batch processing

### 🔧 Configuration

#### Environment Variables (Backend)
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=labelops

# JWT
JWT_SECRET=your-super-secret-jwt-key

# Server
PORT=8080
GIN_MODE=debug
```

#### Environment Variables (Frontend)
```typescript
// src/environments/environment.ts
export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080/api/v1'
};
```

## Development

### Backend Development
```bash
cd backend
go run main.go
```

### Frontend Development
```bash
cd frontend
npm start
```

### Database Migrations
The system automatically creates tables and stored procedures on startup. For production, consider using a proper migration tool.

## Production Deployment

### Backend
1. Build the Go application
2. Set up PostgreSQL with proper credentials
3. Configure environment variables
4. Use a process manager like PM2 or systemd

### Frontend
1. Build the Angular application: `npm run build`
2. Serve the built files with a web server like nginx
3. Configure reverse proxy to the backend API

### Database
1. Set up PostgreSQL with proper security
2. Create dedicated database user
3. Configure connection pooling
4. Set up regular backups

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check PostgreSQL is running
   - Verify credentials in `.env`
   - Ensure database exists

2. **Frontend Can't Connect to Backend**
   - Check backend is running on port 8080
   - Verify CORS configuration
   - Check network connectivity

3. **Print Jobs Not Working**
   - Ensure Zebra Browser Print SDK is loaded
   - Check printer connectivity
   - Verify ZPL content generation

### Logs
- Backend logs are printed to console
- Frontend errors appear in browser console
- Database errors are logged in PostgreSQL logs

## Support

For issues and questions:
1. Check the logs for error messages
2. Verify all prerequisites are installed
3. Ensure environment variables are correctly set
4. Test database connectivity

## License

MIT License - see LICENSE file for details. 