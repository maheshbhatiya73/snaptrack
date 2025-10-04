#!/bin/bash
set -e

INSTALL_DIR="/opt/snaptrack"
CONFIG_DIR="/etc/snaptrack"
CONFIG_FILE="$CONFIG_DIR/config.yaml"

echo "Starting SnapTrack minimal installation..."

# Clone SnapTrack repo
sudo mkdir -p "$INSTALL_DIR"
sudo chown "$USER":"$USER" "$INSTALL_DIR"
if [ ! -d "$INSTALL_DIR/.git" ]; then
    git clone https://github.com/maheshbhatiya73/snaptrack.git "$INSTALL_DIR"
else
    echo "Repo already exists, pulling latest..."
    cd "$INSTALL_DIR" && git pull
fi

# Create sample config.yaml
sudo mkdir -p "$CONFIG_DIR"
if [ ! -f "$CONFIG_FILE" ]; then
    sudo bash -c "cat > $CONFIG_FILE <<EOL
env: production

server:
  host: \"0.0.0.0\"
  port: \"8080\"
  frontend_path: \"/opt/snaptrack/frontend/dist\"
  backend_base_url: \"http://localhost:8080/api\"

cors:
  origins: \"*\"

security:
  jwt_secret: \"changeme-secret\"

database:
  host: \"localhost\"
  port: 5432
  user: \"snapuser\"
  password: \"snap-pass\"
  dbname: \"snaptrack\"
EOL"
    echo "Sample config created at $CONFIG_FILE"
else
    echo "Config already exists at $CONFIG_FILE, skipping..."
fi

# Build Go backend
echo "Building SnapTrack backend..."
cd "$INSTALL_DIR"
go build -o snaptrack-api ./cmd/main.go

# Setup systemd service
echo "Setting up systemd service..."
SERVICE_FILE="/etc/systemd/system/snaptrack.service"
sudo bash -c "cat > $SERVICE_FILE <<EOL
[Unit]
Description=SnapTrack Service
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/snaptrack-api --config $CONFIG_FILE
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOL"

sudo systemctl daemon-reload
sudo systemctl enable snaptrack
sudo systemctl restart snaptrack

echo "SnapTrack installation complete!"
echo "Backend running at http://0.0.0.0:8080"
echo "Config file: $CONFIG_FILE"
