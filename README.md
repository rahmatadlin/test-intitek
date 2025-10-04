# 📦 Warehouse Management System

A full-stack warehouse management system built with Vue.js, Golang, and MySQL. This application provides comprehensive inventory tracking, user authentication, and powerful features for managing warehouse operations.

## ✨ Features

### Core Features

- **Product Management**: Full CRUD operations for inventory items
- **Dashboard**: Real-time statistics and insights
  - Total products count
  - Total stock quantity
  - Low stock alerts
  - Low stock products list
- **Filtering**: Filter products by status and low stock items
- **Responsive UI**: Modern, clean interface with mobile support

### Bonus Features

- **User Authentication**: JWT-based secure authentication system
- **CSV Export**: Download complete product inventory as CSV
- **Barcode Generation**: Generate Code128 barcodes for each product SKU
- **Auto Status Updates**: Automatically updates product status based on quantity

## 🛠️ Tech Stack

### Backend

- **Language**: Go 1.21+
- **Framework**: Gin (HTTP web framework)
- **ORM**: GORM
- **Database**: MySQL
- **Authentication**: JWT (golang-jwt/jwt)
- **Barcode**: boombuler/barcode

### Frontend

- **Framework**: Vue.js 3 (Composition API)
- **Build Tool**: Vite
- **State Management**: Pinia
- **HTTP Client**: Axios
- **Routing**: Vue Router

## 📋 Prerequisites

Before running this application, ensure you have:

- **Go** 1.21 or higher ([Download](https://golang.org/dl/))
- **Node.js** 18.x or higher ([Download](https://nodejs.org/))
- **MySQL** 8.0 or higher ([Download](https://dev.mysql.com/downloads/))
- **Git** (for cloning the repository)

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd test-intitek
```

### 2. Database Setup

Create a MySQL database:

```sql
CREATE DATABASE warehouse_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. Backend Setup

Navigate to the backend directory:

```bash
cd backend
```

Create a `.env` file (copy from `env.example`):

```bash
cp env.example .env
```

Edit `.env` with your MySQL credentials:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=warehouse_db

JWT_SECRET=your-secret-key-change-this-in-production
PORT=8080
```

Install Go dependencies:

```bash
go mod tidy
go mod download
```

Run the backend server:

```bash
go run main.go
```

The backend will start on `http://localhost:8080`

**Note**: A default admin user will be created automatically:

- Username: `admin`
- Password: `admin123`

### 4. Frontend Setup

Open a new terminal and navigate to the frontend directory:

```bash
cd frontend
```

Install dependencies:

```bash
npm install
```

Run the development server:

```bash
npm run dev
```

The frontend will start on `http://localhost:5173`

### 5. Access the Application

Open your browser and visit: `http://localhost:5173`

Login with the default credentials:

- **Username**: `admin`
- **Password**: `admin123`

## 📡 API Documentation

Base URL: `http://localhost:8080/api`

### Authentication Endpoints

#### Register User

```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (201 Created):**

```json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com"
  }
}
```

#### Login

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

**Response (200 OK):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "admin",
    "email": "admin@warehouse.com"
  }
}
```

### Product Endpoints

**Note**: All product endpoints require authentication. Include the JWT token in the Authorization header:

```http
Authorization: Bearer <your-jwt-token>
```

#### Get All Products

```http
GET /api/products
```

**Query Parameters:**

- `status` (optional): Filter by status (`in_stock`, `low_stock`, `out_of_stock`)
- `low_stock` (optional): Set to `true` to show only low stock items

**Response (200 OK):**

```json
{
  "data": [
    {
      "id": 1,
      "name": "Laptop Dell XPS 15",
      "sku": "LAPTOP-001",
      "quantity": 15,
      "location": "Warehouse A, Shelf 12",
      "status": "in_stock",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

#### Get Single Product

```http
GET /api/products/:id
```

**Response (200 OK):**

```json
{
  "data": {
    "id": 1,
    "name": "Laptop Dell XPS 15",
    "sku": "LAPTOP-001",
    "quantity": 15,
    "location": "Warehouse A, Shelf 12",
    "status": "in_stock",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

#### Create Product

```http
POST /api/products
Content-Type: application/json

{
  "name": "Laptop Dell XPS 15",
  "sku": "LAPTOP-001",
  "quantity": 15,
  "location": "Warehouse A, Shelf 12",
  "status": "in_stock"
}
```

**Validation Rules:**

- `name`: Required, string
- `sku`: Required, unique, string
- `quantity`: Required, integer, minimum 0
- `location`: Required, string
- `status`: Required, enum (`in_stock`, `low_stock`, `out_of_stock`)

**Response (201 Created):**

```json
{
  "data": {
    "id": 1,
    "name": "Laptop Dell XPS 15",
    "sku": "LAPTOP-001",
    "quantity": 15,
    "location": "Warehouse A, Shelf 12",
    "status": "in_stock",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

#### Update Product

```http
PUT /api/products/:id
Content-Type: application/json

{
  "name": "Laptop Dell XPS 15",
  "sku": "LAPTOP-001",
  "quantity": 8,
  "location": "Warehouse B, Shelf 5",
  "status": "low_stock"
}
```

**Response (200 OK):**

```json
{
  "data": {
    "id": 1,
    "name": "Laptop Dell XPS 15",
    "sku": "LAPTOP-001",
    "quantity": 8,
    "location": "Warehouse B, Shelf 5",
    "status": "low_stock",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T14:20:00Z"
  }
}
```

#### Delete Product

```http
DELETE /api/products/:id
```

**Response (200 OK):**

```json
{
  "message": "Product deleted successfully"
}
```

### Dashboard Endpoint

#### Get Dashboard Statistics

```http
GET /api/dashboard/stats
```

**Response (200 OK):**

```json
{
  "total_products": 50,
  "total_stock": 1250,
  "low_stock_count": 5,
  "low_stock_products": [
    {
      "id": 3,
      "name": "Mouse Wireless",
      "sku": "MOUSE-003",
      "quantity": 8,
      "location": "Warehouse A, Shelf 3",
      "status": "low_stock",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### Export & Barcode Endpoints

#### Export Products as CSV

```http
GET /api/export/csv
```

**Response (200 OK):**
Returns a CSV file with all products.

**Headers:**

```
Content-Type: text/csv
Content-Disposition: attachment; filename=products.csv
```

#### Generate Barcode

```http
GET /api/barcode/:sku
```

**Example:**

```http
GET /api/barcode/LAPTOP-001
```

**Response (200 OK):**
Returns a PNG image of the barcode.

**Headers:**

```
Content-Type: image/png
Content-Disposition: inline; filename=barcode-LAPTOP-001.png
```

### Error Responses

All endpoints may return the following error responses:

**400 Bad Request:**

```json
{
  "error": "Validation error message"
}
```

**401 Unauthorized:**

```json
{
  "error": "Authorization header required"
}
```

**404 Not Found:**

```json
{
  "error": "Product not found"
}
```

**500 Internal Server Error:**

```json
{
  "error": "Internal server error message"
}
```

## 🗂️ Project Structure

```
test-intitek/
├── backend/
│   ├── config/           # Configuration management
│   │   └── config.go
│   ├── controllers/      # Request handlers
│   │   ├── auth_controller.go
│   │   └── product_controller.go
│   ├── database/         # Database connection
│   │   └── database.go
│   ├── middleware/       # HTTP middlewares
│   │   └── auth.go
│   ├── models/          # Data models
│   │   ├── product.go
│   │   └── user.go
│   ├── routes/          # Route definitions
│   │   └── routes.go
│   ├── main.go          # Application entry point
│   ├── go.mod           # Go dependencies
│   └── env.example      # Environment variables template
│
├── frontend/
│   ├── src/
│   │   ├── assets/      # Static assets
│   │   │   └── styles.css
│   │   ├── router/      # Vue Router configuration
│   │   │   └── index.js
│   │   ├── services/    # API services
│   │   │   ├── api.js
│   │   │   └── productService.js
│   │   ├── stores/      # Pinia stores
│   │   │   └── auth.js
│   │   ├── views/       # Vue components (pages)
│   │   │   ├── Dashboard.vue
│   │   │   ├── Layout.vue
│   │   │   ├── Login.vue
│   │   │   └── Products.vue
│   │   ├── App.vue      # Root component
│   │   └── main.js      # Application entry point
│   ├── index.html       # HTML template
│   ├── package.json     # Node dependencies
│   └── vite.config.js   # Vite configuration
│
└── README.md            # This file
```

## 🎨 Features Walkthrough

### 1. Authentication

- Secure JWT-based authentication
- Login and registration functionality
- Automatic token refresh
- Protected routes

### 2. Dashboard

- Overview of warehouse statistics
- Total products and stock count
- Low stock alerts with product list
- Real-time data updates

### 3. Product Management

- **List View**: Display all products in a clean table format
- **Add Product**: Create new inventory items with validation
- **Edit Product**: Update existing product details
- **Delete Product**: Remove products with confirmation
- **Filter**: Filter by status or show only low stock items

### 4. CSV Export

- Export all products to CSV format
- Includes all product fields
- One-click download

### 5. Barcode Generation

- Generate Code128 barcodes for each product
- View barcode in modal
- Download barcode image

## 🔧 Product Status Logic

Products automatically update their status based on quantity:

- **In Stock**: Quantity > 10
- **Low Stock**: Quantity ≤ 10 and > 0
- **Out of Stock**: Quantity = 0

This logic is implemented in the backend (`models/product.go`) and automatically applied when creating or updating products.

## 🐛 Troubleshooting

### Backend Issues

**Problem**: `Failed to connect to database`

- **Solution**: Verify MySQL is running and credentials in `.env` are correct

**Problem**: `Port 8080 already in use`

- **Solution**: Change the `PORT` in `.env` file or stop the process using port 8080

### Frontend Issues

**Problem**: `CORS errors`

- **Solution**: Ensure the backend is running and CORS is properly configured

**Problem**: `401 Unauthorized`

- **Solution**: Login again to refresh your authentication token

## 🚢 Building for Production

### Backend

```bash
cd backend
go build -o warehouse-api main.go
./warehouse-api
```

### Frontend

```bash
cd frontend
npm run build
```

The built files will be in the `frontend/dist` directory. Serve them using any static file server.

## 📝 Environment Variables

### Backend (.env)

| Variable    | Description        | Default            |
| ----------- | ------------------ | ------------------ |
| DB_HOST     | MySQL host         | localhost          |
| DB_PORT     | MySQL port         | 3306               |
| DB_USER     | MySQL username     | root               |
| DB_PASSWORD | MySQL password     | -                  |
| DB_NAME     | Database name      | warehouse_db       |
| JWT_SECRET  | JWT signing secret | default-secret-key |
| PORT        | Server port        | 8080               |

## 🔒 Security Considerations

- Passwords are hashed using bcrypt
- JWT tokens expire after 24 hours
- Input validation on both frontend and backend
- SQL injection protection via GORM
- CORS configured for specific origins
- Authentication required for all product operations

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 👨‍💻 Author

Created as a technical assessment for warehouse management system development.

## 📞 Support

For issues or questions, please open an issue in the GitHub repository.

---

**Happy Coding! 🚀**
