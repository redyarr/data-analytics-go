

# Student Analytics gRPC Service

## Overview

This project implements a simple gRPC-based server-client application for managing and analyzing student data. The system provides functionalities to add student information, calculate various analytics like average grades, and obtain gender-related statistics. It stores data in a MySQL database using GORM for database interactions.

## Features

- **Add a Student:** Add a student with details such as name, age, grade, and gender.
- **Get a Student:** Retrieve student details by name.
- **Get Average Grade:** Calculate and return the average grade of all students.
- **Get Gender Percentage:** Calculate and return the percentage of male and female students.
- **Get Max/Min Age by Gender:** Retrieve the maximum and minimum ages of students by gender.
- **Get Average Grade by Gender:** Retrieve the average grade of male and female students.
- **Get Combined Data:** Return all analytics data in a single response (average grades, age stats, gender stats, etc.).

## Technology Stack

- **Programming Language:** Go
- **gRPC:** For RPC communication between client and server
- **GORM:** Object-Relational Mapper (ORM) for MySQL database interactions
- **MySQL:** Database for storing student data

## Project Structure

```
.
├── model/
│   └── students.go        # Model definition for the Student entity
├── proto/
│   └── student.proto      # Protocol Buffers definition for gRPC services
├── server.go              # gRPC server implementation
└── client.go              # gRPC client implementation
```

## Setup Instructions

### Prerequisites

- Go installed on your machine.
- MySQL database installed and running.
- `protoc` (Protocol Buffers compiler) installed.

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/student-analytics-grpc.git
cd student-analytics-grpc
```

### 2. Install Go Dependencies

```bash
go mod tidy
```

### 3. Create MySQL Database

Create a MySQL database called `students_db` and adjust the `dsn` in `server.go` to match your local setup.

```sql
CREATE DATABASE students_db;
```

### 4. Run the Server

```bash
go run server.go
```

This starts the gRPC server on port `3000`.

### 5. Run the Client

```bash
go run client.go
```

This invokes the `GetCombinedData` RPC and prints out the analytics in a JSON format.

## gRPC Endpoints

### AddStudent

- **RPC:** `AddStudent(Student) -> Empty`
- **Description:** Adds a new student to the database.

### GetStudent

- **RPC:** `GetStudent(StudentRequest) -> Student`
- **Description:** Fetches student data by name.

### GetAverageGrade

- **RPC:** `GetAverageGrade(Empty) -> AverageGradeResponse`
- **Description:** Returns the average grade of all students.

### GetGenderPercentage

- **RPC:** `GetGenderPercentage(Empty) -> PercentageResponse`
- **Description:** Calculates and returns the percentage of male and female students.

### GetMaxAgeByGender

- **RPC:** `GetMaxAgeByGender(Empty) -> MaxAgeByGenderResponse`
- **Description:** Returns the maximum age of students by gender.

### GetMinAgeByGender

- **RPC:** `GetMinAgeByGender(Empty) -> MinAgeByGenderResponse`
- **Description:** Returns the minimum age of students by gender.

### GetCombinedData

- **RPC:** `GetCombinedData(Empty) -> CombinedResponse`
- **Description:** Returns combined analytics data (students, average grades, max/min ages, and gender percentages).

## License

This project is licensed under the MIT License.

