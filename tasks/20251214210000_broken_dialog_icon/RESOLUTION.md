# Resolution

## Summary
Implemented a custom confirmation dialog service and component using Taiga UI, replacing the native browser `confirm()`.

## Technical Details
1.  **Service**: Created `ConfirmDialogService` to handle dialog opening via `TuiDialogService`.
2.  **Component**: Created `ConfirmDialogComponent` with a custom template.
3.  **Styling**:
    -   Updated `styles.scss` to define `tui-dialog[data-appearance='confirm-dialog']`.
    -   Applied sharp corners (`border-radius: 0`), black borders, and a hard shadow (`box-shadow: 10px 10px 0 ...`) to match the "Balenciaga" aesthetic.
    -   Removed the default Taiga UI header to control styling fully within the component.

## Verification
- [x] Build passed (`make dashboard-build`)
- [x] Verified code changes in `confirm-dialog.service.ts` and `styles.scss`.
