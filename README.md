# momo

## 📌 Overview

This project is a **VPN management system written in Go** with a built-in **Telegram bot** for user interaction.  
It allows users to **subscribe, recharge, and manage their VPN accounts via Telegram**, while the backend automatically manages **server assignment, and configuration updates**.

## ✨ Features

- 🔒 **VPN Management**

  - Start, stop, and monitor VPNs
  - Adapter compatibility checks (whether a VPN can run over another)
  - Automatic migration when a server is drained

- 🤖 **Telegram Bot Integration**

  - Users can interact directly through Telegram
  - Supports charging/renewal of VPN accounts
  - Sends updated VPN configs after server drained

- ⚡ **Dynamic Server Assignment**

  - Assigns users to servers based on capacity and metrics
  - Evicts and reassigns users when an admin drains a server
  - Ensures high availability and balanced load

- 📊 **Admin Tools**
  - Drain server command
  - Monitor active users and server utilization
  - Centralized logging

### Prerequisites

- Go 1.24+
- A Telegram Bot token from [@BotFather](https://t.me/BotFather)
- Admin/root privileges for adapter and VPN handling

### Installation

```bash
git clone https://github.com/mohamadrezamomeni/momo.git
cd momo
go mod tidy
```

### Configuration

Create a `config.json` file:

```yaml
log:
  access_file: "./access.log"
  error_file: "./error.log"
db:
  dialect: sqlite3
  path: "example.db"
  migrations: ./repository/sqllite/migrations

port_assignment:
  start_port: 3000
  end_port: 4000
worker_server:
  address: "localhost"
  port: "50051"
http:
  port: "8008"
admin:
  username: "admin username"
  password: "the password will be used for login admin"
  firstname: "admin's firstname"
  lastname: "admin's lastname"

auth:
  secret_key: "the random charcters for token that will be given for web app"
  expire_time_by_minutes: 20

encrypt:
  encrypt_key: "qFQGz0yE0nYBxqi9"

telegram:
  token: "the token that was given by telegram"
notification:
  telegram_token: "the token that was given by telegram"
```

### Build & Run

```bash
go build -o ./bin/telegram ./cmd/telegram/telegram.go
go build -o ./bin/notification ./cmd/notification/notification.go
go build -o ./bin/schaduler ./cmd/schaduler/schaduler.go
go build -o ./bin/worker ./cmd/worker/worker.go
go build -o ./bin/admin_api ./cmd/admin_api/admin_api.go
./telegram
./notification
./schaduler
./worker
./admin_api
```
