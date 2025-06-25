# Security Policy

## Supported Versions

Please note that only the latest minor release version will be supported, example:

Version `v2025.5.3` is the current latest release in this example.

| Version      | Supported?         |
| ------------ | ------------------ |
| `v2024.12.3` | :x:                |
| `v2025.1.0`  | :x:                |
| `v2025.5.0`  | :white_check_mark: |
| `v2025.5.3`  | :white_check_mark: |

`v2025.5.x` releases will be supported until the next minor release `v2025.6.x` is released.

Still should you be running an older patch release, you should consider upgrading to the latest patch release of the current minor version (e.g., `v2025.5.1` → `v2025.5.3`).

## Reporting a Vulnerability

Should you encounter a security vulnerability or any security-related issues, please **DO NOT** post it publicly.
Instead, report it privately to the FiveNet developers (Galexrt) on Discord via a direct message with a clear description of the vulnerability.

Please include the following information listed below (as much as you can provide) in the report for us better understand the nature and scope of the possible vulnerability:

- Type of issue (e.g. buffer overflow, SQL injection, cross-site scripting, etc.)
- Full paths of source file(s) related to the manifestation of the issue
- The location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit the issue
