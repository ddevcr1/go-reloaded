# Go-Reloaded

**Go-Reloaded** is a command-line tool written in Go that is designed to automatically edit and format text. The program reads text from a specified file, applies a set of transformation rules (case correction, number system conversion, punctuation) to it, and writes the result to a new file.

## Features

The program performs the following transformations:

### 1. Numbering Systems
*   **(hex)**: Converts the word before the command from hexadecimal to decimal.
*   **(bin)**: Converts the word before the command from binary to decimal.

### 2. Text Case Transformation
*   **(up) / (up, <number>)**: Converts the previous word (or the specified number of words) to UPPERCASE.
*   **(low) / (low, <number>)**: Converts the previous word (or the specified number of words) to lowercase.
*   **(cap) / (cap, <number>)**: Converts the previous word (or the specified number of words) to Capitalized Format.

### 3. Punctuation
*   Punctuation marks `, . ! ? : ;` are attached to the previous word and followed by a space.
*   Punctuation groups (e.g., `...` or `!?`) are handled correctly as a single unit.
*   Single quotes `'` wrap the text inside them without internal spaces (e.g., `' word '` becomes `'word'`).

### 4. Grammar
*   The article **a** is automatically replaced with **an** if the following word starts with a vowel (`a, e, i, o, u`) or the letter **h**.

---

## Usage

### Running the program
The program takes two arguments: the path to the input file and the path to the output file.

```bash
go run . sample.txt result.txt
```

## Project Structure
```bash
.
├── modif/
│   ├── modif.go       # Text processing logic (library)
│   └── modif_test.go  # Unit tests for transformation rules
├── main.go            # Entry point (CLI and file handling)
├── go.mod             # Go module definition
└── README.md          # Documentation
```

## Testing

To ensure that all transformation rules work correctly, unit tests have been implemented. You can run them using the following command:

```bash
go test ./modif -v
```

## Requirements
Go version 1.18 or higher.
No external dependencies (Standard Go packages only).
