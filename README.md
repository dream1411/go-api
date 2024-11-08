# My Go Project

This is a Go project with live reload capabilities using `fresh` for development.

## Prerequisites

- Go 1.18 or later
- [Fresh](https://github.com/pilu/fresh) for live reloading

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/dream1411/go-api.git
   cd go-api
   go get
   go get -u github.com/gorilla/mux
   go get -u github.com/go-sql-driver/mysql
   go get -u github.com/swaggo/swag/cmd/swag
   go get -u github.com/swaggo/http-swagger
   go install github.com/swaggo/swag/cmd/swag@latest
   go install github.com/pilu/fresh@latest
### คำอธิบาย
- ไฟล์นี้มีการแนะนำการติดตั้ง `fresh` และการตั้งค่า `.env` รวมถึงการใช้งาน `fresh` สำหรับการ live reload
- ให้คำแนะนำเกี่ยวกับการตั้งค่า `runner.conf` และปัญหาที่อาจพบบน Windows
