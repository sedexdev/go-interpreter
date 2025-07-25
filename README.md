# 📘 C-- Interpreter: Written in Go

This is an interpreter for the fictitious language C--, written in Go

-   C-- syntax corresponds to that of C++
-   C-- only has one type, _integer_, used for numerical and boolean (0 or 1) values

The available operators are:

```bash
+     =     (
-     ==    )
/     !=    {
*     <     }
%     >     ,
      <=
      >=
      &&
      ||
```

## ✨ Features

-   ✅ Declare variables
-   ✅ Perform calculations using an arbitrary number of brackets
-   ✅ Print a list of statements using the 'print' keyword
-   ✅ Declare if/else statements
-   ✅ Declare while loops

## 📦 Installation

### Prerequisites

```bash
Go >= 1.11
```

Go installation files and instructions can be found on the [official website](https://go.dev/doc/install)

### Local Setup

```bash
# Clone the repository
git clone https://github.com/sedexdev/go-interpreter/internal.git
```

## ⚙️ Configuration

This project uses [Go modules](https://go.dev/blog/using-go-modules) so no additional config is required to run the code locally.

## 🛠️ Usage

After cloning the repo:

```bash
cd go-interpreter
go run main.go
```

### C-- Programs

> This interpreter does not support dedicated file types - code samples are strings declared in .go files

Program declared in `internal/testcode/testcode.go`

```val = 104
while (val >= 2) {
    if (val % 2 == 0) {
        next = val / 2
    } else {
        next = 3 * val + 1
    }
    print val, next
    val = next
}
```

💡 _Please experiment with your own code to run your own C-- programs!_

## 📂 Project Structure

```
go-interpreter/
│
├── internal/                  # Internal module source files
├── go.mod                     # Go module file
├── main.go                    # App entry point
└── README.md                  # This README.md file
```

## 🐛 Reporting Issues

Found a bug or need a feature? Open an issue [here](https://github.com/sedexdev/go-interpreter/internal/issues).

## 🧑‍💻 Authors

-   **Andrew Macmillan** – [@sedexdev](https://github.com/sedexdev)

## 📜 License

-   **!!** This project is not licensed due to source code and derivative work being included from the book below

## 📣 Acknowledgements

-   Many thanks to Thorsten Ball for his wonderful book, without which this project would have been significantly harder
-   [_Writing An Interpreter In Go_](https://interpreterbook.com/)
