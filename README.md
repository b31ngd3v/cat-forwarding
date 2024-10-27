# 🐾 Cat Forwarding

> **Open the purr-tals for cats to connect worldwide, forwarding TCP traffic!**

Cat Forwarding is a port-forwarding tool that exposes a port on your local machine to the outside world. It forwards TCP traffic to a remote address, allowing your local cats to communicate globally. Think of it as a secret tunnel for cats to chat beyond their usual boundaries!

### 📦 Backend Repo

Check out the backend implementation of Cat Forwarding [here](https://github.com/b31ngd3v/cf-backend).

## 🚀 Features

- Forward TCP traffic from a local port to a remote address.
- Simple and lightweight, yet effective.
- Ideal for testing, debugging, and, of course, feline communication!

## 🛠️ Installation

### Download Binaries

1. **Download the latest binary** from the [releases page](https://github.com/b31ngd3v/cat-forwarding/releases) for your operating system.

2. **Set up the binary:**

   - **macOS & Linux**:

     ```bash
     chmod +x cat-forwarding
     sudo mv cat-forwarding /usr/local/bin/cat-forwarding
     ```

   - **Windows**:
     - Move the downloaded binary to a folder in your `PATH` (e.g., `C:\Windows\System32`).

### Build from Source

If you want to build Cat Forwarding from source, follow these steps:

1. **Clone the repository:**

   ```bash
   git clone https://github.com/b31ngd3v/cat-forwarding.git
   cd cat-forwarding
   ```

2. **Build the binary:**
   ```bash
   make build
   ```

## 🐾 Usage

To start forwarding TCP traffic from a local port to a remote address:

```bash
cat-forwarding [port]
```

## 🧪 Running Tests

To run the unit tests for Cat-Forwarding, use the following command:

```bash
make test
```

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/b31ngd3v/cat-forwarding/LICENSE) file for details.
