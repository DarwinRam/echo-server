# 🌀 Go TCP Echo Server – Enhanced

This is a feature-rich TCP Echo Server written in Go. It builds upon a basic echo server and includes powerful enhancements like concurrency, logging, command handling, personality responses, inactivity timeout, and overflow protection.

---

## 🚀 How to Run the Server

Make sure you have **Go** installed. Then:

# bash
go run main.go --port=4000

# open new cmd and run

nc localhost 4000

---

## 🧠 Most Educationally Enriching Feature:

The server personality mode was the most educationally enriching. Designing custom responses like hello → Hi there! and bye → Goodbye! required a creative approach to input handling, response logic, and connection lifecycle — all while maintaining clean concurrency.

---

## 🔍 Feature That Required the Most Research:

The inactivity timeout feature required the most research. Learning how SetReadDeadline() works in Go's net.Conn interface, and properly handling read timeouts without crashing the server, involved understanding lower-level networking concepts and Go-specific behavior around timeouts and goroutine cleanup.

---

## 🧩 Supported Features
✅ Concurrency using goroutines

✅ Connection/disconnection logging

✅ Graceful client disconnect handling

✅ Trimmed input + clean echo

✅ Per-client message logging to file

✅ Command-line port flag (--port)

✅ Inactivity timeout (30s)

✅ Overflow protection (1024 bytes max)

✅ Server personality mode (hello, bye, etc.)

✅ Command protocol (/time, /echo, /quit)

---

## Video Link

