# Telegram Language Learning Bot

A production-grade Telegram bot for vocabulary learning, built in Go with a layered architecture (repository → service → handler). It supports importing words, translating them into target languages, and generating example sentences using Google Gemini 2.5 API.

---

## ✨ Features

- `/importword <word> <source_lang> <target_lang>` imports a word, translates it, and stores both the translation and an example sentence.  
- Gemini 2.5 integration for translation + example generation.  
- PostgreSQL database for persistent vocabulary storage.  
- Docker Compose orchestration for bot + database.  
- Defensive error handling and logging for reliable production use.

---

## 🛠️ Architecture

- Go backend with clean layering:
  - `bot` -> for bootstraping
  - `repository` → database access  
  - `service` → translation + example generation  
  - `handler` → Telegram command handling  
- Gemini API used via REST calls (`generateContent`).  
- Docker Compose manages bot and Postgres containers.

---

## ⚙️ Setup (Step-by-Step)

### 1. Clone the Repo  
`git clone https://github.com/yourusername/telegram-langbot.git`  
`cd telegram-langbot`

### 2. Create `.env` File  
Create a `.env` file in the project root and add:

TELEGRAM_BOT_TOKEN=your_telegram_bot_token
GEMINI_API_KEY=your_real_gemini_api_key_here
API_BASE_URL=https://generativelanguage.googleapis.com

POSTGRES_USER=langbot
POSTGRES_PASSWORD=securepassword
POSTGRES_DB=langbot_db
POSTGRES_HOST=postgres
POSTGRES_PORT=5432


### 3. Start Services with Docker Compose  
`docker-compose up --build`

This starts both the Telegram bot and PostgreSQL database.

### 4. Use the Bot  
Open Telegram and send:

`/importword hello en es`

Expected response:  
`Word 'hello' imported successfully: Hola (en → es). Example: Hola, ¿cómo estás?`

---

## 🧾 Bot Command

### `/importword`
Syntax:  
`/importword <word> <source_lang> <target_lang>`

Example:  
`/importword hello en es`

---

## 🧠 Gemini 2.5 Integration

The bot uses Google Gemini 2.5 API to translate the word and generate an example sentence, then stores both in PostgreSQL.

---

## 🗃️ Database Schema

The bot uses a simple vocabulary table:

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| word | TEXT | Original word |
| source_lang | TEXT | Language of word |
| target_lang | TEXT | Translation language |
| translated_word | TEXT | Translation result |
| example_sentence | TEXT | Generated example |
| created_at | TIMESTAMP | Timestamp |

---



