import { Routes } from '@angular/router';
import { authGuard, loginGuard } from './core/auth/auth.guard';
import { AssetDetailComponent } from './features/assets/asset-detail/asset-detail.component';
import { DashboardComponent } from './features/dashboard/dashboard.component';
import { LoginComponent } from './features/login/login.component';

export const routes: Routes = [
    { 
        path: 'login', 
        component: LoginComponent,
        canActivate: [loginGuard]
    },
    { 
        path: 'dashboard', 
        component: DashboardComponent,
        canActivate: [authGuard]
    },
    {
        path: 'assets/:id',
        component: AssetDetailComponent,
        canActivate: [authGuard]
    },
    { 
        path: '', 
        redirectTo: '/dashboard', 
        pathMatch: 'full' 
    },
    {
        path: '**',
        redirectTo: '/dashboard'
    }
];
