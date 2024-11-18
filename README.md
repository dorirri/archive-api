# Doodocs Backend Challenge

## About
A Go-based service that provides three main functionalities through RESTful API endpoints:

### Archive Analysis
Analyzing ZIP archives (file structure, sizes, and MIME types)

### Archive Creation
Creating new ZIP archives from valid file types (.docx, .xml, jpeg, .png)

### Email Distribution
Sending files (.docx, .pdf) via email to multiple recipients

## Docker Setup
```bash
docker build -t archive-service .
docker run -p 8080:8080 archive-service
```
