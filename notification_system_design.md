# Campus Notification Platform

## Stage 1: REST API Design

### Core Actions

| Action | Description |
|--------|------------|
| Create Notification | Admin sends notification to students |
| List Notifications | Student fetches notifications with pagination |
| Get Notification | Fetch a single notification |
| Mark as Read | Mark a notification as read |
| Mark All Read | Mark all notifications as read |
| Delete Notification | Delete a notification |
| Get Unread Count | Get unread notifications count |

---

### API Endpoints

Create Notification  
POST /api/notifications  

Headers:
Content-Type: application/json  
Authorization: Bearer <admin_token>  

Body:
{
  "studentIds": ["stu_001", "stu_002"],
  "type": "Placement",
  "title": "CSX Corporation Hiring",
  "message": "CSX Corporation is visiting campus on 25-Apr-2026.",
  "priority": "high"
}

Response:
{
  "status": "success",
  "data": {
    "notificationId": "uuid",
    "createdAt": "timestamp",
    "recipientCount": 2
  }
}

---

List Notifications  
GET /api/students/{studentId}/notifications?page=1&limit=20&type=Placement&isRead=false  

---

Get Notification  
GET /api/students/{studentId}/notifications/{notificationId}  

---

Mark as Read  
PATCH /api/students/{studentId}/notifications/{notificationId}/read  

---

Mark All Read  
PATCH /api/students/{studentId}/notifications/read-all  

---

Delete Notification  
DELETE /api/students/{studentId}/notifications/{notificationId}  

---

Get Unread Count  
GET /api/students/{studentId}/notifications/unread-count  

---

### Data Model

Notification:
- id: UUID  
- studentId: string  
- type: Placement | Event | Result  
- title: string  
- message: string  
- isRead: boolean  
- createdAt: timestamp  
- updatedAt: timestamp  

---

### Real-Time Notifications (SSE)

Endpoint:
GET /api/students/{studentId}/notifications/stream  

Flow:
- Client opens connection  
- Server pushes notifications  
- Missed events handled using Last-Event-ID  

---

## Stage 2: Database Design

### Database: PostgreSQL

Reasons:
- Structured data  
- Strong indexing  
- ACID compliance  
- Scalable  

---

### Schema

CREATE TYPE notification_type AS ENUM ('Placement', 'Event', 'Result');

CREATE TABLE students (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(255) UNIQUE
);

CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    student_id VARCHAR(50),
    type notification_type,
    title VARCHAR(200),
    message TEXT,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (student_id) REFERENCES students(id)
);

---

### Indexes

CREATE INDEX idx_student_unread
ON notifications (student_id, is_read, created_at DESC);

CREATE INDEX idx_type
ON notifications (type);

CREATE INDEX idx_student_type
ON notifications (student_id, type, created_at DESC);

---

## Stage 3: Query Optimization

Problem Query:
SELECT * FROM notifications
WHERE student_id = 1042 AND is_read = false
ORDER BY created_at DESC;

Issues:
- Full table scan  
- Unnecessary columns  
- No pagination  

---

Optimized Query:
SELECT id, type, title, message, is_read, created_at
FROM notifications
WHERE student_id = 1042 AND is_read = FALSE
ORDER BY created_at DESC
LIMIT 20 OFFSET 0;

---

Index:
CREATE INDEX idx_student_unread_time
ON notifications (student_id, is_read, created_at DESC);

---

Why not index everything:
- Slows writes  
- Uses more storage  
- Hard to maintain  

---

## Stage 4: Caching Strategy

Problem:
Frequent reads overload DB  

Solution: Redis  

Cache Keys:
notif:{studentId}:recent  
notif:{studentId}:unread_count  

Strategy:
- Cache-aside  
- TTL 30–60 sec  
- Invalidate on writes  

Benefits:
- Faster response  
- Reduced DB load  

---

## Stage 5: Notify All Design

Problems:
- Slow sequential execution  
- No retry  
- Failure stops process  

---

Improved Approach:

1. Batch insert into DB  
2. Publish event to queue  
3. Process emails asynchronously  

---

Architecture:
API → DB → Queue → Worker → Email Service  

---

Benefits:
- Scalable  
- Fault tolerant  
- Faster  
- Supports retries  

---
