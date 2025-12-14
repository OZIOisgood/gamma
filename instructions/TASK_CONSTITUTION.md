# Task Constitution

## 1. Philosophy
Every significant change or feature implementation must be tracked as a discrete task. This ensures traceability, context preservation, and a clear history of decisions.

## 2. Task Structure
All tasks are stored in the `tasks/` directory at the project root.

### Naming Convention
Folder name: `YYYYMMDDHHMMSS_snake_case_description`
Example: `20251214210000_fix_dialog_icon`

### Directory Contents
Each task folder must contain:
1.  **`README.md`**: The problem statement, requirements, and context. (Renamed from `task.md` for better tool support).
2.  **`RESOLUTION.md`**: A summary of the changes made, technical decisions, and verification steps.

## 3. Workflow

### 1. Creation
Use the helper script to generate a new task:
```bash
./scripts/create_task.sh "my task name"
```
This creates the folder and template files.

### 2. Definition
Fill out `README.md` with:
-   **User Story/Requirement**: As a <role>, I want <feature>...
-   **Context**: Links to relevant files, existing issues, or design docs.
-   **Acceptance Criteria**: What defines "done".

### 3. Execution
-   Implement the changes.
-   Follow the **Project Constitution** for coding standards and QA.

### 4. Resolution
Fill out `RESOLUTION.md` with:
-   **Summary**: What was changed.
-   **Verification**: How it was tested (e.g., "Ran `make dashboard-build`", "Checked UI").
-   **Next Steps**: Any follow-up tasks.

### 5. Completion
-   Mark the item as done in `ISSUES.md` if applicable.
-   Commit the task folder along with the code changes.
