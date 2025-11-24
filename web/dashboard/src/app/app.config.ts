import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
import { ApplicationConfig, provideBrowserGlobalErrorListeners } from '@angular/core';
import { provideAnimations } from "@angular/platform-browser/animations";
import { provideRouter } from '@angular/router';
import { provideEventPlugins } from "@taiga-ui/event-plugins";
import { tuiPasswordOptionsProvider } from '@taiga-ui/kit';

import { routes } from './app.routes';

export const appConfig: ApplicationConfig = {
  providers: [
        provideAnimations(),
        provideBrowserGlobalErrorListeners(),
        provideHttpClient(withInterceptorsFromDi()),
    provideRouter(routes),
        provideEventPlugins(),
        tuiPasswordOptionsProvider({
            icons: {
                hide: '@tui.material.sharp.visibility_off',
                show: '@tui.material.sharp.visibility',
            },
        }),
    ]
};
