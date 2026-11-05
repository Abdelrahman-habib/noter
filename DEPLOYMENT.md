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
   - Add a PostgreSQL database service
   - Railway will auto-deploy

3. **Configure Environment Variables**

   - Set `SESSION_SECRET` to a random string
   - Railway provides a PostgreSQL `DATABASE_URL`; the `railway.json` start
     command already passes it to the app via `-dsn=$DATABASE_URL` (the app
     only reads `DB_DSN` or the `-dsn` flag, never `DATABASE_URL` directly)

4. **Run Migrations**

   ```bash
   # Connect to Railway CLI
   railway login
   railway link

   # Run migrations (goose needs the "postgres" dialect and connection string;
   # the quoted $DATABASE_URL is expanded by railway's child shell)
   railway run -- sh -c 'goose -dir db/schema/migrations postgres "$DATABASE_URL" up'
   ```

### Railway Configuration

- **Port**: Uses `$PORT` environment variable
- **Database**: PostgreSQL via `$DATABASE_URL`, passed as the Go `-dsn` flag
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
# goose needs driver and connection string (or GOOSE_DRIVER/GOOSE_DBSTRING env vars)
goose -dir db/schema/migrations postgres "$DATABASE_URL" up
```

Production setups should provision their own PostgreSQL roles (a limited
application role plus a migration owner role) similar to
`db/init/01-create-databases-and-users.sql`, with the migration owner
creating the schema.

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
