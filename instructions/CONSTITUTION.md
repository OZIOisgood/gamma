# Project Constitution

## 1. Context & Awareness
- **Source of Truth**: Always consult `README.md`, `Makefile`, and `dev/yaak` requests to understand the current architecture and API capabilities.
- **Project Status**: Pre-release. Data is ephemeral.

## 2. Database & Infrastructure
- **Migrations**: Since data is not yet valuable, prefer **editing existing migrations** to refine the schema rather than creating new ones.
- **Environment**: We frequently restart the Docker environment (`make docker-restart`).

## 3. Frontend & Design (Dashboard)
- **Framework**: Angular with Taiga UI.
- **Aesthetic**: "Balenciaga-like" — sharp, minimalistic, black & white.
- **Reference**: Follow patterns in `web/dashboard/src/styles.scss`.

## 4. Quality Assurance
- **Build Check**: Always run a build (`make build` or `make dashboard-build`) before marking a task as complete.
- **Verification**: Ensure no regressions are introduced.

## 5. Task Management
- **Protocol**: Follow the guidelines in [TASK_CONSTITUTION.md](TASK_CONSTITUTION.md).
- **Issues**: Check `ISSUES.md` for the roadmap and task status.
