# 🚗 Spot-Sync

**Spot-Sync** is a robust, high-performance RESTful API service built in Go, designed to manage real-time parking space availability and reservation workflows. 

Spot-Sync provides structured management of multiple parking zones (such as standard spaces, EV charging stations, and covered slots), dynamically tracks availability, and processes reservations securely with concurrency protections.

---

## 💡 Real-World Problem Solved

Traditional parking systems frequently suffer from inefficient usage, stale slot data, and booking conflicts (overbooking). Spot-Sync solves these issues by addressing three key real-world challenges:

### 1. Eliminating Overbooking (Race Conditions)
In busy parking lots, multiple drivers might try to reserve the last remaining parking space at the exact same millisecond. If not handled correctly, this creates a race condition that leads to double-booking.
*   **The Spot-Sync Solution:** During reservation creation, Spot-Sync wraps the capacity check and reservation creation inside a database transaction using an **atomic row-level lock (`SELECT ... FOR UPDATE` via GORM)** on the target parking zone. This ensures that only one request can read and decrement the capacity at a time, strictly guaranteeing that the total capacity is never breached.

### 2. Real-Time Dynamic Capacity Tracking
Drivers need to know exactly how many slots are open before driving to a lot, avoiding wasted fuel, emissions, and search time.
*   **The Spot-Sync Solution:** Spot-Sync dynamically calculates the available spots for any parking zone by subtracting the count of active reservations from the zone's total capacity. This calculation is performed efficiently on the database level via a subquery, providing instant, accurate availability figures to drivers.

### 3. Smart Resource Management & Access Control
Not all parking spots are identical, and not all users have the same privileges. EV drivers need EV spots, and lot owners need global visibility.
*   **The Spot-Sync Solution:** 
    *   **Zone Segmentation:** Supports distinct parking categories—`general`, `ev_charging`, and `covered`—each configured with its own hourly pricing and slot count.
    *   **Role-Based Access Control (RBAC):** Authenticates users via JSON Web Tokens (JWT) and restricts actions according to roles. *Drivers* can query spots and manage their own bookings. *Admins* can add zones, define pricing, and audit all reservations across the system.

---

## 🛠️ Tech Stack

*   **Language:** Go (Golang)
*   **Web Framework:** Echo v5 (High-performance, minimalist router)
*   **ORM:** GORM (v2)
*   **Database:** PostgreSQL (Cloud instance via Neon PostgreSQL)
*   **Authentication:** JWT (JSON Web Tokens)
*   **Validation:** Go-Playground Validator v10
*   **Hot Reloading:** Air

---

## ⚙️ Configuration & Setup

### 1. Prerequisites
*   Go (version 1.20 or higher)
*   PostgreSQL database (or cloud database connection string)

### 2. Environment Setup
Create a `.env` file in the root directory of the project:

```env
PORT=8080
DSN="postgresql://<username>:<password>@<host>/<database>?sslmode=require"
JWT_SECRET="your-super-secure-key"
JWT_EXPIRES_HOURS=24
```

### 3. Run the Server

#### Option A: Running with Live Reloading (Recommended for Development)
If you have [Air](https://github.com/cosmtrek/air) installed:
```bash
air
```

#### Option B: Standard Go Run
```bash
go run cmd/main.go
```
The server will start up on the port specified in your `.env` file (defaults to `8080`), migrate the database tables automatically, and output:
```text
Connected to db
⇨ http server started on [::]:8080
```

---

## 🔌 API Architecture & Response Envelopes

All API endpoints follow a consistent structure. Requests with JSON bodies require the `Content-Type: application/json` header.

### Success Response Envelope
```json
{
  "success": true,
  "message": "Action completed successfully",
  "data": { ... }
}
```

### Error Response Envelope
```json
{
  "code": 400,
  "message": "Brief error summary",
  "details": "Technical details or validation failure reasons"
}
```

---

## 📡 API Endpoint Reference

### 🔐 1. Authentication (`/api/v1/auth`)

#### **Register User**
*   **URL:** `/api/v1/auth/register`
*   **Method:** `POST`
*   **Auth Required:** No
*   **Request Body:**
    ```json
    {
      "name": "John Doe",
      "email": "john@example.com",
      "password": "securepassword",
      "role": "driver" 
    }
    ```
    *(Note: Valid values for `role` are `driver` or `admin`)*
*   **Success Response (201 Created):**
    ```json
    {
      "success": true,
      "message": "User registered successfully",
      "data": {
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com",
        "role": "driver",
        "created_at": "2026-06-28T14:00:00Z",
        "updated_at": "2026-06-28T14:00:00Z"
      }
    }
    ```

#### **Login User**
*   **URL:** `/api/v1/auth/login`
*   **Method:** `POST`
*   **Auth Required:** No
*   **Request Body:**
    ```json
    {
      "email": "john@example.com",
      "password": "securepassword"
    }
    ```
*   **Success Response (201 Created):**
    ```json
    {
      "success": true,
      "message": "Login successfully",
      "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user": {
          "id": 1,
          "name": "John Doe",
          "email": "john@example.com",
          "role": "driver"
        }
      }
    }
    ```

---

### 🅿️ 2. Parking Zones (`/api/v1/zones`)

#### **Create Parking Zone**
*   **URL:** `/api/v1/zones`
*   **Method:** `POST`
*   **Auth Required:** Yes (Admin Only)
*   **Request Body:**
    ```json
    {
      "name": "North Sector EV Bay",
      "type": "ev_charging",
      "total_capacity": 5,
      "price_per_hour": 4.50
    }
    ```
    *(Note: Valid values for `type` are `general`, `ev_charging`, or `covered`)*
*   **Success Response (201 Created):**
    ```json
    {
      "success": true,
      "message": "Zone created successfully",
      "data": {
        "id": 1,
        "name": "North Sector EV Bay",
        "type": "ev_charging",
        "total_capacity": 5,
        "price_per_hour": 4.5,
        "created_at": "2026-06-28T14:10:00Z"
      }
    }
    ```

#### **Get All Parking Zones**
*   **URL:** `/api/v1/zones`
*   **Method:** `GET`
*   **Auth Required:** No (Public)
*   **Success Response (200 OK):**
    ```json
    {
      "success": true,
      "message": "Zones retrieved successfully",
      "data": [
        {
          "id": 1,
          "name": "North Sector EV Bay",
          "type": "ev_charging",
          "total_capacity": 5,
          "available_spots": 4,
          "price_per_hour": 4.5,
          "created_at": "2026-06-28T14:10:00Z"
        }
      ]
    }
    ```

#### **Get Parking Zone by ID**
*   **URL:** `/api/v1/zones/:id`
*   **Method:** `GET`
*   **Auth Required:** No (Public)
*   **Success Response (200 OK):**
    ```json
    {
      "success": true,
      "message": "Zone retrieved successfully",
      "data": {
        "id": 1,
        "name": "North Sector EV Bay",
        "type": "ev_charging",
        "total_capacity": 5,
        "available_spots": 4,
        "price_per_hour": 4.5,
        "created_at": "2026-06-28T14:10:00Z"
      }
    }
    ```

---

### 📅 3. Reservations (`/api/v1/reservations`)

#### **Reserve a Parking Spot**
*   **URL:** `/api/v1/reservations`
*   **Method:** `POST`
*   **Auth Required:** Yes (Driver & Admin)
*   **Request Body:**
    ```json
    {
      "zone_id": 1,
      "license_plate": "NY-789-AB"
    }
    ```
*   **Success Response (201 Created):**
    ```json
    {
      "success": true,
      "message": "Reservation created successfully",
      "data": {
        "id": 12,
        "user_id": 2,
        "zone_id": 1,
        "license_plate": "NY-789-AB",
        "status": "active",
        "created_at": "2026-06-28T14:15:00Z",
        "updated_at": "2026-06-28T14:15:00Z"
      }
    }
    ```
    *If a zone reaches total capacity, concurrent booking attempts will immediately fail with a `400 Bad Request` or `500 Internal Server Error` detailing that the `parking zone is completely full`.*

#### **Get My Reservations**
*   **URL:** `/api/v1/reservations/my-reservations`
*   **Method:** `GET`
*   **Auth Required:** Yes (Driver & Admin)
*   **Success Response (200 OK):**
    ```json
    {
      "success": true,
      "message": "My reservations retrieved successfully",
      "data": [
        {
          "id": 12,
          "license_plate": "NY-789-AB",
          "status": "active",
          "created_at": "2026-06-28T14:15:00Z",
          "zone": {
            "id": 1,
            "name": "North Sector EV Bay",
            "type": "ev_charging"
          }
        }
      ]
    }
    ```

#### **Cancel Reservation**
*   **URL:** `/api/v1/reservations/:id`
*   **Method:** `DELETE`
*   **Auth Required:** Yes (Reservation Owner or Admin only)
*   **Success Response (200 OK):**
    ```json
    {
      "success": true,
      "message": "Reservation cancelled successfully",
      "data": null
    }
    ```

#### **Get All Reservations (Admin Audit)**
*   **URL:** `/api/v1/reservations`
*   **Method:** `GET`
*   **Auth Required:** Yes (Admin Only)
*   **Success Response (200 OK):**
    ```json
    {
      "success": true,
      "message": "Reservations retrieved successfully",
      "data": [
        {
          "id": 12,
          "user_id": 2,
          "zone_id": 1,
          "license_plate": "NY-789-AB",
          "status": "active",
          "created_at": "2026-06-28T14:15:00Z",
          "zone": {
            "id": 1,
            "name": "North Sector EV Bay",
            "type": "ev_charging"
          }
        }
      ]
    }
    ```

---

## 🧪 Step-by-Step cURL Testing Workflow

Follow this sequence to test all features of the application locally from your command line.

### Step 1: Register an Admin and a Driver
First, create the users. One will act as the manager/admin and the other as the customer/driver.

```bash
# Register Admin
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "Admin User", "email": "admin@spotsync.com", "password": "AdminSecurePassword", "role": "admin"}'

# Register Driver
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "Driver Sam", "email": "sam@driver.com", "password": "SamSecurePassword", "role": "driver"}'
```

### Step 2: Authenticate and Extract Tokens
Perform login requests to obtain the JWT bearer tokens.

```bash
# Login as Admin (Save the token from the response)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@spotsync.com", "password": "AdminSecurePassword"}'

# Login as Driver (Save the token from the response)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "sam@driver.com", "password": "SamSecurePassword"}'
```
*Note: For the subsequent authenticated commands, replace `<ADMIN_TOKEN>` and `<DRIVER_TOKEN>` with the token string returned under `"data" -> "token"`.*

### Step 3: Create a Parking Zone (Admin Only)
Using the Admin token, create a parking zone. Let's make one with a total capacity of `2`.

```bash
curl -X POST http://localhost:8080/api/v1/zones \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <ADMIN_TOKEN>" \
  -d '{"name": "Downtown EV Station", "type": "ev_charging", "total_capacity": 2, "price_per_hour": 6.00}'
```

### Step 4: Check Availability (Public)
Request the list of zones to verify the zone is created and has 2 available spots.

```bash
curl -X GET http://localhost:8080/api/v1/zones
```

### Step 5: Reserve a Spot (Driver)
Make a booking in the newly created zone (assuming the zone ID is `1`).

```bash
curl -X POST http://localhost:8080/api/v1/reservations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <DRIVER_TOKEN>" \
  -d '{"zone_id": 1, "license_plate": "EV-COOL-2"}'
```

### Step 6: Verify Decreased Availability
Get the details for the specific zone to see that the `available_spots` value has decreased to `1`.

```bash
curl -X GET http://localhost:8080/api/v1/zones/1
```

### Step 7: View Bookings
Inspect bookings as a driver and audit all bookings as an admin.

```bash
# Get logged-in driver's personal list
curl -X GET http://localhost:8080/api/v1/reservations/my-reservations \
  -H "Authorization: Bearer <DRIVER_TOKEN>"

# Audit all system bookings (Admin only)
curl -X GET http://localhost:8080/api/v1/reservations \
  -H "Authorization: Bearer <ADMIN_TOKEN>"
```

### Step 8: Cancel Reservation (Driver)
Cancel the booking to free up space in the lot (assuming the reservation ID is `1`).

```bash
curl -X DELETE http://localhost:8080/api/v1/reservations/1 \
  -H "Authorization: Bearer <DRIVER_TOKEN>"
```
Verify via `GET /api/v1/zones/1` that availability has returned to `2` spots.
