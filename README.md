# Label Printing Log System

A comprehensive label printing and logging system built with Angular 18, Go, and PostgreSQL.

## 🏗 Architecture

### Frontend (Angular 18 + Tailwind CSS)
- **Label Management**: Display and manage label data with duplicate detection
- **Print Integration**: Zebra Browser Print SDK for direct printing
- **Authentication**: JWT/OAuth with role-based access control
- **Audit Features**: CSV export and audit logging

### Backend (Go/Gin)
- **RESTful API**: Handle label data processing and print jobs
- **Database Integration**: PostgreSQL with stored procedures
- **Authentication**: JWT token management
- **Print Job Queue**: Retry mechanism for failed prints

### Database (PostgreSQL)
- **Label Logs**: Track all printed labels
- **Print Jobs**: Monitor print job status and retries
- **User Management**: Authentication and authorization
- **Audit Logs**: Complete audit trail

## 🚀 Quick Start

### Prerequisites
- Node.js 18+
- Go 1.21+
- PostgreSQL 14+
- Zebra Browser Print SDK

### Backend Setup
```bash
cd backend
go mod tidy
go run main.go
```

### Frontend Setup
```bash
cd frontend
npm install
ng serve
```

### Database Setup
```bash
# Run the SQL scripts in backend/db/
psql -U your_user -d your_database -f backend/db/procedures.sql
```

## 📁 Project Structure

```
LabelOps/
├── backend/                 # Go backend application
├── frontend/               # Angular 18 frontend application
├── docs/                   # Documentation
└── README.md              # This file
```

## 🔧 Features

- ✅ Label data management with duplicate detection
- ✅ Direct printing via Zebra Browser Print SDK
- ✅ Role-based access control (RBAC)
- ✅ Audit logging and CSV export
- ✅ Print job retry mechanism
- ✅ QR code generation support
- ✅ JWT authentication
- ✅ Responsive UI with Tailwind CSS

## 📄 License

MIT License 