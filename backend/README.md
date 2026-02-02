# Echo Chat App - Polyglot Persistence Backend

This backend implements a **Polyglot Persistence** architecture using multiple databases for different purposes:

- **MySQL** - Stores structured relational data (users, friends, groups)
- **MongoDB** - Stores chat messages and conversation history
- **Redis** - Caching layer for sessions, online status, and performance optimization

## 📁 Project Structure

```
backend
├── internal/
│   ├── delivery/
│   │   └── http/
│   │       ├── controller/   # Controller (Terima input, panggil usecase)
│   │       ├── middleware/   # Penjaga pintu (Auth, Log, dll)
│   │       └── router/       # Definisi rute API
│   ├── usecase/              # Business Logic (Aturan main aplikasi)
│   ├── repository/           # Data Access (Query MongoDB/BSON)
│   └── models/               # Struct data & Interface (Kontrak)
├── config/                   # Database & Env loading
└── main.go                   # Entry point & Dependency Injection 
└── go.mod
```

## 🗄️ Database Architecture

### MySQL (Relational Data)
**Purpose**: Store structured data with relationships

**Tables**:
- `users` - User accounts
- `friendships` - Friend relationships
- `groups` - Chat groups
- `group_members` - Group membership

**Why MySQL?**
- Strong ACID guarantees for user data
- Complex relationships (friends, group members)
- Data integrity constraints
- Efficient joins for relational queries

### MongoDB (Document Store)
**Purpose**: Store chat messages and conversation history

**Collections**:
- `chat_messages` - Individual messages
- `conversations` - Conversation summaries

**Why MongoDB?**
- Flexible schema for different message types
- High write throughput for chat messages
- Efficient queries for time-series data
- Easy to scale horizontally
- Natural fit for nested data (attachments, read receipts)

### Redis (Cache & Real-time)
**Purpose**: Caching and real-time data

**Use Cases**:
- User session storage
- Online/offline status
- Unread message counts
- Conversation list caching
- Rate limiting (future)
- Pub/Sub for real-time messaging (future)

**Why Redis?**
- In-memory storage for ultra-fast access
- Built-in expiration for cache invalidation
- Pub/Sub for WebSocket integration
- Atomic operations for counters

## 🚀 Getting Started

### Prerequisites

1. **Install MySQL**
   ```bash
   # Download from https://dev.mysql.com/downloads/mysql/
   # Or use Docker:
   docker run -d -p 3306:3306 --name mysql \
     -e MYSQL_ROOT_PASSWORD=password \
     -e MYSQL_DATABASE=echo_chat \
     mysql:8.0
   ```

2. **Install MongoDB**
   ```bash
   # Download from https://www.mongodb.com/try/download/community
   # Or use Docker:
   docker run -d -p 27017:27017 --name mongodb mongo:latest
   ```

3. **Install Redis**
   ```bash
   # Download from https://redis.io/download
   # Or use Docker:
   docker run -d -p 6379:6379 --name redis redis:latest
   ```

### Installation

1. **Install Go dependencies**
   ```bash
   cd backend
   go mod download
   ```

2. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

3. **Run the application**
   ```bash
   go run main.go
   ```

The server will:
- ✅ Connect to MySQL, MongoDB, and Redis
- ✅ Auto-migrate MySQL tables
- ✅ Create MongoDB indexes
- ✅ Start on port 8080

## 📚 API Examples

### User Operations (MySQL)

**Create User**
```bash
POST /api/v1/users
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123",
  "full_name": "John Doe"
}
```

**Get User**
```bash
GET /api/v1/users/1
```

**Add Friend**
```bash
POST /api/v1/friends
{
  "user_id": 1,
  "friend_id": 2
}
```

**Accept Friend Request**
```bash
PUT /api/v1/friends/1/accept
```

**Get Friends**
```bash
GET /api/v1/users/1/friends
```

### Message Operations (MongoDB)

**Send Message**
```bash
POST /api/v1/messages
{
  "content": "Hello!",
  "type": "text",
  "sender_id": 1,
  "recipient_id": 2
}
```

**Send Group Message**
```bash
POST /api/v1/messages
{
  "content": "Hello everyone!",
  "type": "text",
  "sender_id": 1,
  "group_id": 1
}
```

**Get Messages**
```bash
# Direct messages
GET /api/v1/messages?user_id=1&recipient_id=2&limit=50&page=1

# Group messages
GET /api/v1/messages?user_id=1&group_id=1&limit=50&page=1
```

**Get Conversations**
```bash
GET /api/v1/conversations?user_id=1
```

## 🔧 Code Examples

### Querying MySQL
```go
// Get user with friends
var user models.User
config.DB.MySQL.Preload("Friends").First(&user, userID)

// Create friendship
friendship := models.Friendship{
    UserID:   1,
    FriendID: 2,
    Status:   "pending",
}
config.DB.MySQL.Create(&friendship)
```

### Querying MongoDB
```go
// Insert message
collection := config.DB.MongoDB.Collection("chat_messages")
result, err := collection.InsertOne(ctx, message)

// Find messages
filter := bson.M{"sender_id": userID}
cursor, err := collection.Find(ctx, filter)
```

### Using Redis Cache
```go
// Cache user session
config.Cache.CacheUserSession(ctx, sessionID, userID, 24*time.Hour)

// Get cached data
var data MyStruct
err := config.Cache.Get(ctx, "my-key", &data)

// Set online status
config.Cache.SetUserOnlineStatus(ctx, userID, "online")
```

## 🏗️ Architecture Patterns

### 1. **Data Separation**
- **MySQL**: Source of truth for user identity and relationships
- **MongoDB**: Optimized for message storage and retrieval
- **Redis**: Temporary cache, can be rebuilt from MySQL/MongoDB

### 2. **Cross-Database Operations**
When sending a message:
1. Validate sender/recipient exists in **MySQL**
2. Store message in **MongoDB**
3. Invalidate conversation cache in **Redis**

### 3. **Cache Strategy**
- **Cache-Aside Pattern**: Check cache first, fallback to database
- **Write-Through**: Update cache when data changes
- **TTL-based Expiration**: Auto-expire cached data

### 4. **Data Flow Example**

**Sending a Message:**
```
Client → API → MySQL (validate users) → MongoDB (store message) → Redis (invalidate cache) → Response
```

**Getting Conversations:**
```
Client → API → Redis (check cache) → [Cache Miss] → MongoDB (query) → Redis (cache result) → Response
```

## 🔐 Security Considerations (TODO)

- [ ] Hash passwords with bcrypt
- [ ] Implement JWT authentication
- [ ] Add rate limiting with Redis
- [ ] Validate user permissions for group access
- [ ] Sanitize user input
- [ ] Use prepared statements (GORM handles this)

## 📈 Performance Optimization

### MySQL
- Indexed foreign keys
- Connection pooling (10 idle, 100 max)
- Prepared statements via GORM

### MongoDB
- Compound indexes on `sender_id + created_at`
- Index on `group_id` for group messages
- Capped collections for old messages (optional)

### Redis
- Connection pooling
- Pipeline operations for bulk updates
- Pub/Sub for real-time notifications

## 🧪 Testing

```bash
# Run tests
go test ./...

# Test database connections
curl http://localhost:8080/health
```

## 📝 Next Steps

1. **Authentication**: Implement JWT-based auth
2. **WebSocket**: Add real-time messaging with WebSocket
3. **File Upload**: Implement file storage for media messages
4. **Notifications**: Use Redis Pub/Sub for push notifications
5. **Search**: Add full-text search for messages
6. **Analytics**: Track message statistics

## 🤝 Contributing

This is a learning project demonstrating polyglot persistence patterns in Go.

## 📄 License

MIT
