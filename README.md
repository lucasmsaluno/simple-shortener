# Simple Shortener
> Status: ✅ Finished

<br/>

## 💻 About The Project

A simple and functional URL shortener built with **Go** on the backend and **JavaScript** on the frontend. Users can generate short links and get redirected to the original URLs.

<br/>

## 🎯 Features
- Short an URL

<br/>

## 🌍 See the project

If you want to see the project running in your browser, check it out
<br/>
<br/>
👉  [Working Project](https://simple-shortener.vercel.app/)

<br/>

## ⚡Installation

```bash
# Clone the repository
git clone https://github.com/your-username/url-shortener.git
cd url-shortener

# Build and start containers
docker-compose up --build
```
<br/>

## 📂 Folder Structure
```bash
├── backend/                         # Backend service written in Go
│   ├── cmd/                         # Entry point of the application
│   │   └── main.go                  # Starts the HTTP server
│   ├── handlers/                    # HTTP request handlers
│   │   └── shortener.go             # Logic for shortening and redirecting URLs
│   ├── models/                      # Database models and structs
│   │   └── url.go                   # URL model and DB operations
│   ├── utils/                       # Utility functions
│   │   └── generator.go             # Generates random short codes
│   ├── go.mod                       # Go module definition
│   ├── go.sum                       # Dependency checksums
│   ├── Dockerfile                   # Docker build instructions for backend
│   └── docker-compose.yml           # Docker Compose config for services
├── frontend/                        # Frontend files served to users
│   ├── index.html                   # Main HTML page
│   ├── assets/                      # Static assets
│   │   ├── js/                      # JavaScript logic
│   │   │   └── app.js               # Handles form and API calls
│   │   └── styles/                  # CSS styles
│   │       └── style.css            # Styling for the UI
├── .gitignore                       # Git ignored files
└── README.md    
```
<br/>

## ⚙️ Technologies

- GO
- Javascript
- Docker
- PostgreSQL
