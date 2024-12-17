# Gema CLI

Gema CLI is an AI assistant for your terminal that can run commands based on your queries. It leverages the power of AI to understand your requests and suggest appropriate commands to execute.

## DEMO

[![Gema CLI Demo](https://img.youtube.com/vi/yklcZ4dNYWg/0.jpg)](https://youtu.be/yklcZ4dNYW)

## Installation

To install Gema CLI, clone the repository and build the executable:

```sh
git clone https://github.com/4nkitd/gemini-cli.git
cd gemini-cli
go build -o ?
```

## Usage

Before using Gema CLI, you need to set the following environment variables:

- `GENAI_DEFAULT_MODEL`: The default AI model to use.
- `GENAI_API_KEY`: Your API key for accessing the AI service.

You can set these environment variables in your terminal:

```sh
export GENAI_DEFAULT_MODEL="models/gemini-2.0-flash-exp" 
export GENAI_API_KEY=your_api_key
```

To use Gema CLI, simply run the executable with your query:

```sh
./? "your query here"
```

Gema CLI will process your query, provide a response, and suggest a command to run. You will be prompted to confirm whether you want to execute the suggested command.

## Example

```sh
./? "list all files in the current directory"
```

Output:

```
Response:
Here is the command to list all files in the current directory:

Suggested Command to RUN: ls -la

Run command (y for yes, n for no):
```

If you type `y`, the command will be executed, and the output will be displayed.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.
