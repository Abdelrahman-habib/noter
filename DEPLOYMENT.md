# Deployment Guide

This guide covers various deployment strategies for the Noter application.

## 🚀 Railway (Recommended - Free)

Railway offers a free tier with $5 monthly credit, perfect for testing applications.

### Prerequisites

- GitHub account
- Railway account (free)

### Steps

1. **Push to GitHub**

   ```bash
   git add .
   git commit -m "Add deployment configuration"
   git push origin main
   ```

2. **Deploy on Railway**

   - Go to [railway.app](https://railway.app)
   - Sign in with GitHub
   - Click "New Project" → "Deploy from GitHub repo"
   - Select your repository
   - Add MySQL database service
   - Railway will auto-deploy

3. **Configure Environment Variables**

   - Set `SESSION_SECRET` to a random string
   - Database URL is automatically provided

4. **Run Migrations**

   ```bash
   # Connect to Railway CLI
   railway login
   railway link

   # Run migrations
   railway run goose -dir db/schema/migrations up
   ```

### Railway Configuration

- **Port**: Uses `$PORT` environment variable
- **Database**: Automatic MySQL connection via `$DATABASE_URL`
- **HTTPS**: Automatic SSL certificate
- **Domain**: Automatic subdomain (e.g., `yourapp.railway.app`)

## 🌐 Render (Free with Limitations)

### Steps

1. Connect GitHub repository
2. Select "Web Service"
3. Build command: `go build -o bin/noter ./cmd/web`
4. Start command: `./bin/noter -addr=:$PORT -env=production -dsn=$DATABASE_URL`
5. Add PostgreSQL database service

### Limitations

- Sleeps after 15 minutes of inactivity
- Cold start takes ~30 seconds

## ☁️ Fly.io (Free Tier)

### Steps

1. Install Fly CLI
2. Create `fly.toml`:

   ```toml
   app = "your-app-name"

   [build]
     builder = "paketobuildpacks/builder:base"

   [[services]]
     http_checks = []
     internal_port = 4000
     processes = ["app"]
     protocol = "tcp"
     script_checks = []

     [services.concurrency]
       hard_limit = 25
       soft_limit = 20
       type = "connections"

     [[services.ports]]
       force_https = true
       handlers = ["http"]
       port = 80

     [[services.ports]]
       handlers = ["tls", "http"]
       port = 443

     [[services.tcp_checks]]
       grace_period = "1s"
       interval = "15s"
       restart_limit = 0
       timeout = "2s"
   ```

3. Deploy:
   ```bash
   fly launch
   fly deploy
   ```

## 🐳 Docker Deployment

### For VPS/Cloud Providers

1. **Build and push image**:

   ```bash
   docker build -t your-app .
   docker tag your-app your-registry/your-app
   docker push your-registry/your-app
   ```

2. **Deploy with docker-compose**:
   ```bash
   docker-compose -f docker-compose.yml --env-file prod.env up -d
   ```

## 🔧 Required Modifications for Production

### 1. Remove TLS for Cloud Platforms

Most cloud platforms handle HTTPS automatically. Modify the app to work without TLS:

```go
// In main.go, make TLS optional
if config.tlsCert != "" && config.tlsKey != "" {
    err = srv.ListenAndServeTLS(config.tlsCert, config.tlsKey)
} else {
    err = srv.ListenAndServe()
}
```

### 2. Environment Variables

- `DATABASE_URL`: Database connection string
- `SESSION_SECRET`: Random string for session encryption
- `PORT`: Port number (set by platform)

### 3. Database Migrations

Run migrations as part of deployment:

```bash
goose -dir db/schema/migrations up
```

## 💰 Cost Comparison

| Platform     | Free Tier       | Database | HTTPS | Best For            |
| ------------ | --------------- | -------- | ----- | ------------------- |
| Railway      | $5/month credit | ✅       | ✅    | Testing, small apps |
| Render       | 750 hours/month | ✅       | ✅    | Demos, portfolios   |
| Fly.io       | 3 VMs           | ✅       | ✅    | Production apps     |
| DigitalOcean | $100 credit     | ✅       | ✅    | Serious projects    |

## 🎯 Recommendation

For a **testing application**, use **Railway**:

- ✅ Free tier with generous limits
- ✅ Easy setup and deployment
- ✅ Built-in database
- ✅ Automatic HTTPS
- ✅ GitHub integration

## 📝 Next Steps

1. Choose a platform
2. Push your code to GitHub
3. Follow the platform-specific deployment steps
4. Configure environment variables
5. Run database migrations
6. Test your deployed application
