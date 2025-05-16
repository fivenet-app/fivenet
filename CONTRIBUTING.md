# Contributing to Fivenet

## Getting Started

To contribute to Fivenet, please follow these steps:

1. Fork the repository on GitHub.
2. Clone the forked repository to your local machine using `git clone https://github.com/fivenet-app/fivenet-app.github.io.git`.
3. Create a new branch for your feature or bug fix using `git checkout -b <branch_name>`.
4. Make your changes to the codebase, be sure to check the [code style section](#code-style) and adhere to our coding standards.
5. Test your changes thoroughly to ensure they work as expected.
    - If you are adding new features or fixing bugs, please include tests to cover your changes.
    - Ensure that all existing tests pass before submitting your changes.
6. Commit your changes and add a descriptive commit message explaining what you have done in the commit using `git commit`.
    - Please follow the [commit message formatting](#commit-message-formatting).
7. Push your changes to your forked repository using `git push origin <branch_name>`.
8. Create a pull request from your branch to the main branch of the original repository.

## Commit Message Formatting

When submitting a pull request, please ensure that your commit messages follow the guidelines below. This helps maintain a clear and consistent history for the project.

- Follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification for commit messages.
- Use the correct prefix for your commit type:
    - `feat`: A new feature.
    - `fix`: A bug fix.
    - `docs`: Documentation only changes.
    - `style`: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc).
    - `refactor`: A code change that neither fixes a bug nor adds a feature.
    - `perf`: A code change that improves performance.
    - `test`: Adding missing or correcting existing tests.
    - `chore`: Changes to auxiliary tools and libraries used for, e.g., file generation (proto generated files).
    - `ci`: Changes to our CI configuration files and scripts (e.g., GitHub Actions workflows).
- Include a brief description of what you have done in the commit message (if necessary).
- If applicable, include a reference to an issue or ticket number that this commit addresses.

## Coding Style

When writing code, please follow the following guidelines:

- Follow the existing coding style and conventions used in the project.
- Use meaningful variable and function names.
- Write clear and concise comments to explain complex code or logic.
    - Comments should be in English and should explain the "why" behind the code, not just the "what".
- Ensure that your code is properly formatted and adheres to the project's style guide.
    - Make sure to use the project's linter and formatter to check your code (e.g., `eslint` + `prettier` for frontend code, gofmt for backend code).
- Use consistent indentation and spacing.
- Avoid using magic numbers or hard-coded values; use constants or configuration files instead.
