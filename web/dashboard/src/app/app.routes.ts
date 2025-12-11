import { Routes } from '@angular/router';
import { authGuard, loginGuard } from './core/auth/auth.guard';
import { MainLayoutComponent } from './core/layout/main-layout/main-layout.component';
import { AssetDetailComponent } from './features/assets/asset-detail/asset-detail.component';
import { DashboardComponent } from './features/dashboard/dashboard.component';
import { LoginComponent } from './features/login/login.component';
import { RealmDetailComponent } from './features/realms/realm-detail/realm-detail.component';
import { RealmsComponent } from './features/realms/realms.component';

export const routes: Routes = [
    { 
        path: 'login', 
        component: LoginComponent,
        canActivate: [loginGuard]
    },
    {
        path: 'dashboard',
        redirectTo: '/default/dashboard',
        pathMatch: 'full'
    },
    {
        path: ':realm',
        component: MainLayoutComponent,
        canActivate: [authGuard],
        children: [
            { 
                path: 'dashboard', 
                component: DashboardComponent
            },
            {
                path: 'assets/:id',
                component: AssetDetailComponent
            },
            {
                path: 'realms',
                component: RealmsComponent
            },
            {
                path: 'realms/:id',
                component: RealmDetailComponent
            }
        ]
    },
    { 
        path: '', 
        redirectTo: '/default/dashboard', 
        pathMatch: 'full' 
    },
    {
        path: '**',
        redirectTo: '/default/dashboard'
    }
];
