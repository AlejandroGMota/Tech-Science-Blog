# === Development ===
install:
	cd frontend && npm install

dev-front:
	cd frontend && npm run dev

dev-back:
	cd backend && ADMIN_PASS=dev123 go run ./cmd/api/

dev:
	trap 'kill 0' EXIT; make dev-back & make dev-front & wait

build-front:
	cd frontend && npx vite build

build-admin:
	cd frontend && npx vite build --config vite.admin.config.js
	mv backend/admin-dist/admin.html backend/admin-dist/index.html 2>/dev/null || true

build-back:
	cd backend && CGO_ENABLED=0 go build -o tech-blog-api ./cmd/api/

build: build-front build-admin build-back

clean:
	rm -rf backend/public-dist backend/admin-dist backend/tech-blog-api

# === VM Setup (run inside VM) ===
vm-deps:
	sudo apt update && sudo apt install -y git curl ufw nginx certbot python3-certbot-nginx

vm-go:
	wget https://go.dev/dl/go1.22.0.linux-arm64.tar.gz
	sudo tar -C /usr/local -xzf go1.22.0.linux-arm64.tar.gz
	rm go1.22.0.linux-arm64.tar.gz
	echo 'export PATH=$$PATH:/usr/local/go/bin' >> ~/.bashrc

vm-node:
	curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
	sudo apt install -y nodejs

vm-clone:
	cd ~ && git clone https://github.com/AlejandroGMota/Tech-Science-Blog.git

vm-service:
	sudo cp deploy/tech-blog.service /etc/systemd/system/
	sudo systemctl daemon-reload
	sudo systemctl enable tech-blog
	sudo systemctl restart tech-blog

vm-firewall:
	sudo ufw allow 22/tcp
	sudo ufw allow 80/tcp
	sudo ufw allow 443/tcp
	sudo ufw --force enable

vm-nginx:
	sudo cp deploy/nginx-blog.conf /etc/nginx/sites-available/blog
	sudo ln -sf /etc/nginx/sites-available/blog /etc/nginx/sites-enabled/
	sudo rm -f /etc/nginx/sites-enabled/default
	sudo nginx -t && sudo systemctl reload nginx

vm-ssl:
	sudo certbot --nginx -d blog.alejandrogmota.com

vm-setup: vm-deps vm-go vm-node vm-clone vm-firewall vm-nginx vm-service

# === VM Deploy (run inside VM) ===
vm-deploy:
	cd ~/Tech-Science-Blog && git pull origin main
	make build
	sudo systemctl restart tech-blog

vm-status:
	sudo systemctl status tech-blog

vm-logs:
	sudo journalctl -u tech-blog -f

vm-restart:
	sudo systemctl restart tech-blog
