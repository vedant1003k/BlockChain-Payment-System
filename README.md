# 🔗 Simple Blockchain App

A lightweight blockchain implementation in Go, complete with wallet generation, transactions, mining, and a web-based frontend interface. Perfect for learning or demoing how blockchains work!

---

## 🧰 Features

- ✅ Create wallets (ECDSA keypairs)
- ✅ Mint tokens to wallets
- ✅ Send transactions
- ✅ Mine blocks (Proof-of-Work)
- ✅ View full blockchain history
- ✅ Simple REST API + React frontend

---

## 🖥️ Tech Stack

- **Backend**: Golang (net/http)
- **Frontend**: React (Vite)
- **Blockchain**: In-memory PoW chain
- **Keys**: ECDSA with SHA256

---

## 🚀 How to Run Locally

### 1. Clone the Repository

```bash
git clone https://github.com/vedant1003k/BlockChain-Payment-System.git
cd BlockChain-Payment-System
```

---

### 2. Start the Backend (Go API)

#### ➤ Step into the backend folder:

```bash
cd blockchain_backend
```

#### ➤ Run the server:

```bash
go run main.go
```

> This starts your backend server on `http://localhost:8080`

---

### 3. Start the Frontend (React App)

#### ➤ Open another terminal:

```bash
cd blockchain
npm install
npm start
```

> This starts your frontend at `http://localhost:3000`

---

## 🧪 Try It Out

1. Open your browser → [http://localhost:5173](http://localhost:3000)
2. Use the interface to:
   - **Create wallets**
   - **Send transactions**
   - **Mine blocks**
   - **View the blockchain**

---

## 🧪 Postman Support

You can also test the backend using **Postman**:

### ➤ Example: Mint Tokens

```http
POST http://localhost:8080/transaction
Content-Type: application/json

{
  "sender": "MINT",
  "receiver": "0xABC123...",
  "amount": 100,
  "signature": ""
}
```

### ➤ Example: Check Balance

```http
GET http://localhost:8080/balance/0xABC123...
```

---

## 📁 Project Structure

```
blockchain-app/
│
├── backend/         # Go server (API + blockchain logic)
│   └── main.go
│
├── frontend/        # React app (UI)
│   ├── src/
│   └── index.html
│
├── README.md
```

---

## 📌 Notes

- All data is stored in memory (no database)
- CORS is enabled for local testing
- Ideal for educational or prototype purposes
