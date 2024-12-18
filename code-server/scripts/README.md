# Endar Script Editor

The **Endar Script Editor** is a powerful, web-based code management and editing platform designed to simplify the storage, editing, and version control of scripts within Endar. Built with [Code Server](https://github.com/coder/code-server), a web-based Visual Studio Code (VS Code) environment, Endar's scripting editor provides a familiar and robust interface for developers while emphasizing security and ease of use.

## Features

- **Web-Based VS Code Environment**: Leverages [Code Server](https://github.com/coder/code-server) to enable a browser-accessible VS Code interface for seamless script management.
- **Git Integration**: Full support for Git version control, allowing you to manage your repositories efficiently.
- **Secure Terminal Access**: A locked-down terminal ensures safety while allowing necessary operations like version control, editing, and navigation.
- **Directory Isolation**: Only the `scripts` directory is accessible for modification; all other directories are securely locked down.
- **Auto Save Enabled**: The Endar script editor automatically saves your files as you work, ensuring no progress is lost during editing. Changes are saved in real-time, allowing you to focus on development without manual save interruptions.

### Terminal Guidelines

The terminal is a powerful tool but must be used responsibly. Here are the acceptable use cases:
- **Git Operations**: Commit, push, pull, and manage your repositories.
- **File Editing**: Use terminal-based editors (e.g., `vim`, `nano`) if preferred.
- **Navigation**: Browse and navigate the directory structure within the `user` directory.

#### **Important Security Note**
Do **NOT** execute arbitrary commands in the terminal. Commands outside the allowed scope can compromise the system or your data.

## Git Integration

1. Open terminal:
   Use ctrl+` to open the built-in terminal.
2. Clone your repositories into the `user` directory:
   ```bash
   git clone <repository-url>
   cd <repository-name>
   ```
3. Commit your changes:
   ```bash
   git add .
   git commit -m "Your commit message"
   git push
   ```
4. Pull updates from the remote repository:
   ```bash
   git pull
   ```

## Limitations
- **Terminal Access**: Restricted to ensure system integrity. Only the user directory is accessible.
- **Extensions**: Limited to those pre-installed or approved for use in this environment.
- **System-Wide Changes**: Users cannot make system-wide changes or access directories outside user.

## Best Practices
- Always double-check terminal commands to avoid unintended consequences.
- Regularly commit and push your changes to ensure your work is backed up.
- Use the VS Code interface for most editing and management tasks to minimize terminal usage.
- Notify your administrator immediately if you suspect any security concerns.

## Disclaimer
The Endar Script Library is built using [Code Server](https://github.com/coder/code-server), an open-source project licensed under the MIT License. The Code Server software and its associated documentation are provided by [Coder Technologies Inc.](https://github.com/coder).

### MIT License
```
The MIT License

Copyright (c) 2019 Coder Technologies Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```
For more information about Code Server, visit the [Code Server GitHub repository](https://github.com/coder/code-server).


## Support
For questions, issues, or feedback, for Endar please submit an issue on [Endar's Github](https://github.com/queball1999/endar).
