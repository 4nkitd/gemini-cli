package main

var SystemInstruction = `**Template:**

**1. Define the Goal/Problem:**
"I need to clearly state what you want to accomplish using the CLI]. For example:
* 'I need to automate a task of renaming all files in a directory, replacing spaces with underscores.'
* 'I'm encountering an error with '[command]' and I don't know how to solve it.'
* 'I'm trying to understand how to use '[command]' with '[option]'.'
* 'I want to quickly count how many files of a particular extension are in the current directory'

**2. Context (if needed):** 
"I am using macOS and my workflow usually involves the terminal, also I have experience with the zsh."

**3. Provide Specific Information:**
"Here's the command/code/error message I'm working with: '[Paste the command or error here, or describe the scenario precisely]'"

**4. Desired Action:**
"Can you help me by [clearly state what type of help you want]?" For example:
* "suggesting a complete command to achieve my goal,"
* "explaining the error message and possible causes,"
* "providing an alternative way to do this,"
* "explaining the functionality of the command in a concise way"
* "giving a practical example of how this command works"
* "suggesting a few options for better readability"

**5. Format Preference (optional):**
"Please provide the response in a clear and structured format, and prefer one-liners if possible"

**Example using the template:**

"I need to automate a task of renaming all files in a directory, replacing spaces with underscores.
I am using macOS and my workflow usually involves the terminal, also I have experience with the zsh.
Here is an example file name: 'My File Name.txt'.
Can you help me by suggesting a complete command to achieve my goal, and provide the command in a single line."`

var ResponseSchema string = `{
  "type": "object",
  "properties": {
    "response": {
      "type": "string"
    },
    "command": {
      "type": "string"
    }
  },
  "required": [
    "response"
  ]
}`
