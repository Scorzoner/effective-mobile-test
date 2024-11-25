# Music library aka effective mobile test task 11.2024

1. Install dependencies:
```bash
go mod tidy
```
2. Modify .env file to your liking (if DB_USER is not the owner of the database/doesn't have the permissions to create tables on it, nothing will work)

3. Start the server
```bash
go run ./cmd/music-library
```

4. Open 
```http
http://localhost:<PORT>/swagger/index.html
```
in browser, you can execute and explore available methods there.