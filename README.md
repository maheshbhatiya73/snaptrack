# snaptrack: Scalable Backup & Server Management System
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8)](https://golang.org/)  
[![React Version](https://img.shields.io/badge/React-18.2+-61DAFB)](https://reactjs.org/)  
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)  
[![Build Status](https://img.shields.io/github/workflow/status/yourusername/gobackup/build)](https://github.com/maheshbhatiya73/gobackup/actions)  
[![Slack](https://img.shields.io/badge/Slack-Join%20Us-blue)](https://slack.golangbridge.org/)

**snaptrack** is a robust, secure, and user-friendly backup and server management system built with **Go** for the backend and **React** for the frontend. It provides scheduled and manual backups, real-time server monitoring via WebSocket, and secure authentication using Linux PAM. Designed for Linux systems, it requires root or sudo privileges to manage critical operations like backups, services, and firewall rules.

---

## Table of Contents
- [Features](#features)
- [Architecture](#architecture)
- [System Flow](#system-flow)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
  - [Frontend Authentication](#frontend-authentication)
  - [Real-Time Dashboard](#real-time-dashboard)
  - [Backup Management](#backup-management)
- [PAM Authentication](#pam-authentication)
- [Upcoming Features](#upcoming-features)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgements](#acknowledgements)

---

## Features

- **Backup Management**:
  - Create, update, delete, and list backups (manual, full, incremental).
  - View detailed backup logs (ID, Status, Message, Created At).
  - Supports multiple file types: `zip`, `tar`, `tar.gz`.
  - Scheduled backups with customizable cron expressions.
- **Real-Time Monitoring**:
  - Dashboard displays CPU, RAM, disk, network, and uptime via WebSocket.
  - Powered by Go's `gopsutil` for system metrics.
- **Secure Authentication**:
  - Linux PAM-based authentication requiring root or sudo privileges.
  - Frontend login with JWT-based session management.
- **User-Friendly Frontend**:
  - Built with React, Tailwind CSS, and Framer Motion for smooth animations.
  - Responsive table for backups with pagination and status badges.
  - Modal for viewing detailed logs with a clean, modern UI.
- **Extensible Design**:
  - Modular Go backend with REST API and WebSocket support.

---

## Architecture

GoBackup follows a client-server architecture:

- **Backend (Go)**:
  - REST API for CRUD operations on backups (`/api/backups`).
  - WebSocket endpoint (`/ws`) for real-time server metrics.
  - PAM authentication for secure access.
  - mongodb for persistent storage.
- **Frontend (React)**:
  - `/main/`: Real-time server status (CPU, RAM, disk, network, uptime).
  - `/main/backups`: Backup management with create, update, delete, and logs view.
  - Tailwind CSS for styling, Headless UI for modals, and React Tooltip for UX.

---

## System Flow

```mermaid
graph TD
    A[User] -->|Login with Credentials| B[Frontend: /login]
    B -->|POST /api/auth| C[Backend: PAM Auth]
    C -->|Validate Root/Sudo| D{Auth Success?}
    D -->|Yes| E[Issue JWT]
    D -->|No| F[Error: Access Denied]
    E --> G[Frontend: /main/]
    G -->|WebSocket /ws| H[Backend: Stream Metrics]
    G -->|/main/backups| I[Backup Management]
    I -->|GET /api/backups| J[List Backups]
    I -->|POST /api/backups| K[Create Backup]
    I -->|PUT /api/backups/:id| L[Update Backup]
    I -->|DELETE /api/backups/:id| M[Delete Backup]
    I
