# Contributing to Yathahi

Thank you for your interest in contributing to **Yathahi**! We welcome bug fixes, new features, documentation improvements, and ideas from the community.

## 1. Licensing & Recognition (No CLA Required)

Yathahi is permanently licensed under the **Apache License 2.0**, and we do **not** require you to sign a Contributor License Agreement (CLA) or assign away your copyright.

* **Inbound = Outbound:** By submitting a Pull Request to this repository, you agree that your contributions will be licensed under the same [Apache License 2.0](LICENSE) that covers the project.
* **Standard Apache 2.0 Terms:** All rights, permissions, and patent grants for your contribution are governed strictly by the official [Apache License 2.0](LICENSE) without any extra agreements or conditions.
* **Full Public Credit:** Your authorship is permanently preserved in our Git commit history, release notes, and the `CONTRIBUTORS.md` file.

## 2. How to Submit a Contribution

1. **Fork the Repository:** Click the **Fork** button at the top right of the main repository page to create your own copy.
2. **Clone Your Fork:**
   ```bash
   git clone https://github.com<your-username>/yathahi.git
   cd yathahi
   
   ```
3. **Create a New Branch:** Always make your changes on a dedicated branch rather than `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```
4. **Make Your Changes & Commit:** Write clean, well-tested code and commit it using your verified GitHub email so your profile gets full credit:
   ```bash
   git add .
   git commit -m "Describe your change clearly"
   ```
5. **Push to Your Fork:**
   ```bash
   git push origin feature/your-feature-name
   ```
6. **Open a Pull Request:** Go to your fork on GitHub, click **Contribute -> Open pull request**, and submit it against the `main` branch of `yathahi/yathahi`. Once our automated checks pass and a maintainer reviews your code, it will be merged!

## 3. Keeping Your Fork Up to Date

Before starting new work, sync your fork with the latest changes from the main repository by clicking **Sync fork -> Update branch** on your fork's GitHub page, then run:

```bash
git pull origin main
```
