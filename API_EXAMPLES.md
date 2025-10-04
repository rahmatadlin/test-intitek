# API Examples

This document provides practical examples for testing the Warehouse Management API.

## Setup

First, get your authentication token:

```bash
# Login and save the token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.token'
```

Set the token as an environment variable:

```bash
export TOKEN="your-token-here"
```

## Authentication Examples

### 1. Register a New User

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123"
  }' | jq
```

**Expected Response:**

```json
{
  "message": "User registered successfully",
  "user": {
    "id": 2,
    "username": "john_doe",
    "email": "john@example.com",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }' | jq
```

**Expected Response:**

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

## Product Management Examples

### 3. Create Products

**Example 1: Laptop**

```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Laptop Dell XPS 15",
    "sku": "LAPTOP-001",
    "quantity": 25,
    "location": "Warehouse A, Shelf 12",
    "status": "in_stock"
  }' | jq
```

**Example 2: Wireless Mouse (Low Stock)**

```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Wireless Mouse Logitech MX Master 3",
    "sku": "MOUSE-001",
    "quantity": 8,
    "location": "Warehouse A, Shelf 3",
    "status": "low_stock"
  }' | jq
```

**Example 3: Keyboard**

```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Mechanical Keyboard RGB",
    "sku": "KEYB-001",
    "quantity": 50,
    "location": "Warehouse B, Shelf 7",
    "status": "in_stock"
  }' | jq
```

**Example 4: Monitor (Out of Stock)**

```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Monitor 27 inch 4K",
    "sku": "MON-001",
    "quantity": 0,
    "location": "Warehouse A, Shelf 15",
    "status": "out_of_stock"
  }' | jq
```

### 4. Get All Products

```bash
curl -X GET http://localhost:8080/api/products \
  -H "Authorization: Bearer $TOKEN" | jq
```

### 5. Get Products with Filters

**Get only low stock products:**

```bash
curl -X GET "http://localhost:8080/api/products?low_stock=true" \
  -H "Authorization: Bearer $TOKEN" | jq
```

**Get products by status:**

```bash
# In Stock
curl -X GET "http://localhost:8080/api/products?status=in_stock" \
  -H "Authorization: Bearer $TOKEN" | jq

# Out of Stock
curl -X GET "http://localhost:8080/api/products?status=out_of_stock" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### 6. Get Single Product

```bash
curl -X GET http://localhost:8080/api/products/1 \
  -H "Authorization: Bearer $TOKEN" | jq
```

### 7. Update Product

```bash
curl -X PUT http://localhost:8080/api/products/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Laptop Dell XPS 15 (Updated)",
    "sku": "LAPTOP-001",
    "quantity": 30,
    "location": "Warehouse B, Shelf 5",
    "status": "in_stock"
  }' | jq
```

### 8. Delete Product

```bash
curl -X DELETE http://localhost:8080/api/products/1 \
  -H "Authorization: Bearer $TOKEN" | jq
```

## Dashboard Examples

### 9. Get Dashboard Statistics

```bash
curl -X GET http://localhost:8080/api/dashboard/stats \
  -H "Authorization: Bearer $TOKEN" | jq
```

**Expected Response:**

```json
{
  "total_products": 4,
  "total_stock": 83,
  "low_stock_count": 1,
  "low_stock_products": [
    {
      "id": 2,
      "name": "Wireless Mouse Logitech MX Master 3",
      "sku": "MOUSE-001",
      "quantity": 8,
      "location": "Warehouse A, Shelf 3",
      "status": "low_stock"
    }
  ]
}
```

## Export & Barcode Examples

### 10. Export Products as CSV

```bash
curl -X GET http://localhost:8080/api/export/csv \
  -H "Authorization: Bearer $TOKEN" \
  --output products.csv

# View the CSV
cat products.csv
```

### 11. Generate Barcode

```bash
# Download barcode as PNG
curl -X GET http://localhost:8080/api/barcode/LAPTOP-001 \
  -H "Authorization: Bearer $TOKEN" \
  --output barcode.png

# View the image (Linux)
xdg-open barcode.png

# View the image (Mac)
open barcode.png

# View the image (Windows)
start barcode.png
```

## Batch Operations

### 12. Create Multiple Products at Once

Create a script `create_products.sh`:

```bash
#!/bin/bash

TOKEN="your-token-here"

# Array of products
products=(
  '{"name":"USB Cable Type-C","sku":"USB-001","quantity":100,"location":"Warehouse A, Shelf 1","status":"in_stock"}'
  '{"name":"HDMI Cable 2m","sku":"HDMI-001","quantity":75,"location":"Warehouse A, Shelf 2","status":"in_stock"}'
  '{"name":"Ethernet Cable Cat6","sku":"ETH-001","quantity":5,"location":"Warehouse B, Shelf 4","status":"low_stock"}'
  '{"name":"Power Adapter 65W","sku":"PWR-001","quantity":40,"location":"Warehouse A, Shelf 8","status":"in_stock"}'
  '{"name":"Laptop Stand","sku":"STAND-001","quantity":15,"location":"Warehouse B, Shelf 1","status":"in_stock"}'
)

# Create each product
for product in "${products[@]}"
do
  echo "Creating product..."
  curl -X POST http://localhost:8080/api/products \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "$product" | jq
  echo "---"
done
```

Make it executable and run:

```bash
chmod +x create_products.sh
./create_products.sh
```

## Error Handling Examples

### 13. Invalid Token

```bash
curl -X GET http://localhost:8080/api/products \
  -H "Authorization: Bearer invalid-token" | jq
```

**Expected Response (401):**

```json
{
  "error": "Invalid or expired token"
}
```

### 14. Duplicate SKU

```bash
# Try to create a product with existing SKU
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Another Laptop",
    "sku": "LAPTOP-001",
    "quantity": 10,
    "location": "Warehouse C",
    "status": "in_stock"
  }' | jq
```

**Expected Response (400):**

```json
{
  "error": "Failed to create product. SKU might already exist."
}
```

### 15. Invalid Product Data

```bash
# Missing required field
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Test Product",
    "quantity": 10
  }' | jq
```

**Expected Response (400):**

```json
{
  "error": "Key: 'Product.SKU' Error:Field validation for 'SKU' failed on the 'required' tag"
}
```

### 16. Product Not Found

```bash
curl -X GET http://localhost:8080/api/products/9999 \
  -H "Authorization: Bearer $TOKEN" | jq
```

**Expected Response (404):**

```json
{
  "error": "Product not found"
}
```

## Testing with Postman

### Import to Postman

1. Create a new collection "Warehouse Management API"
2. Add environment variables:

   - `baseUrl`: http://localhost:8080
   - `token`: (will be set after login)

3. Create folders:

   - Authentication
   - Products
   - Dashboard
   - Export

4. Add requests as per examples above

### Postman Pre-request Script

Add this to your collection to auto-update the token:

```javascript
// If token is not set, skip
if (!pm.collectionVariables.get("token")) {
  console.log("No token found");
}
```

### Postman Test Script

Add this to verify responses:

```javascript
pm.test("Status code is 200", function () {
  pm.response.to.have.status(200);
});

pm.test("Response has data", function () {
  var jsonData = pm.response.json();
  pm.expect(jsonData).to.have.property("data");
});
```

## Performance Testing

### Load Test with Apache Bench

```bash
# Test product listing endpoint
ab -n 1000 -c 10 \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/products
```

## Monitoring

### Health Check

```bash
# Simple health check
curl http://localhost:8080/health | jq

# Continuous monitoring (every 5 seconds)
watch -n 5 'curl -s http://localhost:8080/health | jq'
```

## Tips

1. **Save responses to files:**

   ```bash
   curl ... > response.json
   ```

2. **Pretty print JSON:**

   ```bash
   curl ... | jq '.'
   ```

3. **Follow redirects:**

   ```bash
   curl -L ...
   ```

4. **Show response headers:**

   ```bash
   curl -i ...
   ```

5. **Verbose output for debugging:**
   ```bash
   curl -v ...
   ```

---

For more information, see the main README.md file.
