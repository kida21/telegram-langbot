# Telegram Language Learning Bot

A production‑grade Telegram bot for vocabulary learning, built in Go with a layered architecture (repository → service → handler).  
It supports importing words, translating them into target languages, and generating example sentences using **Google Gemini 2.5 API**.

---

## ✨ Features

- `/importword <word> <source_lang> <target_lang>`  
  Imports a word, translates it, and stores both the translation and an example sentence.
- **Gemini 2.5 integration** for translation + example generation.
- **PostgreSQL database** for persistent vocabulary storage.
- **Docker Compose orchestration** for bot + database.
- Defensive error handling and logging for reliable production use.

---

## 🛠️ Architecture

- **Go backend** with clean layering:
  - `repository` → database access
  - `service` → translation + example generation
  - `handler` → Telegram command handling
- **Gemini API** used via REST calls (`generateContent`).
- **Docker Compose** manages bot and Postgres containers.

---

## ⚙️ Setup

### 1. Clone the Repo
```bash
git clone https://github.com/yourusername/telegram-langbot.git
cd telegram-langbot