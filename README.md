## Xom Telegram Bot

Xom is a Go-based Telegram bot that simplifies file conversions, handling audio, images, and some more soon...

I built this project to learn Golang, focusing on messaging queues, goroutines, and how to deploy Go applications. You can get it on [Telegram](https://t.me/xomFileConverterBot).

### Tech Stack

- Go
- Telegram Bot API
- [RabbitMQ](https://www.cloudamqp.com/): Message queue
- [Neon](https://neon.tech/): PostgeSQL
- [Railway](https://railway.app?referralCode=Xf2JoB): Deployment

### **Features**

- Audio Conversion: Convert audio files to a variety of formats, including:
  - mp4, mp3, wav, flac, ogg, aac, wma, m4a

- Image Conversion: Convert image files to multiple formats, such as:
  - jpg, jpeg, png, gif, pdf, webp, bmp, tif, tiff, ico, avif

- Future Enhancements:
  - Video file conversion: Planned but pending due to workload.
  - Document file conversion: Planned but on hold due to other priorities.

### **Commands**

- `/start`: Greets the user and initiates the bot.
- `/help`: Provides a list of available commands and features.
- File Upload: Simply upload a file, and Xom will automatically recognize the file type and offer appropriate conversion options.

### **Installation**

1. Clone the repository:
   ```bash
   git clone https://github.com/41x3n/Xom.git
   cd xom
   ```
2. Install the required Go modules:
   ```bash
   go mod tidy
   ```  
3. Set up environment variables by creating a `.env` file in the root directory:

   ```
   APP_ENV=your_environment
   DSN=your_database_source_name
   TELEGRAM_BOT_TOKEN=your_telegram_bot_token
   RABBITMQ_URL=your_rabbitmq_url
   CONTEXT_TIMEOUT=your_context_timeout
   ```
4. Run the bot:
   ```bash
   go run cmd/xom.go
   ```

### **Usage**

Start a conversation with Xom on Telegram, upload your file, and let the bot guide you through the conversion process.

### **Contributing**

Contributions are welcome! If you have ideas or improvements, feel free to open an issue or submit a pull request.

