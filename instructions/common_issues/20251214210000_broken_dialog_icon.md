# Broken Dialog Close Icon

**Date:** 2025-12-14 21:00:00
**Status:** Fixed (Global Configuration)

## Description
The default Taiga UI dialog close button attempts to load `x.svg` from the assets folder. However, the project is configured to use **Material Design** icons (specifically the sharp variant), and the default Taiga UI icons are not copied to the assets folder during the build. This results in a broken image icon.

## Solution
We configured the `tuiCommonIconsProvider` globally in `app.config.ts` to use the Material Sharp close icon (`@tui.material.sharp.close`) instead of the default `close` icon. This fixes the issue for all dialogs and other components that use the common close icon.

### Code Change
In `web/dashboard/src/app/app.config.ts`:

```typescript
import { tuiCommonIconsProvider } from '@taiga-ui/core';

export const appConfig: ApplicationConfig = {
  providers: [
    // ...
    tuiCommonIconsProvider({
        close: '@tui.material.sharp.close',
    }),
    // ...
  ]
};
```

## Correct Icon Usage
For future reference, when using icons in this project, always use the **Material Sharp** icon set references instead of direct SVG paths or default icon names.

**Correct Pattern:** `@tui.material.sharp.<icon_name>`
