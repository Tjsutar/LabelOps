import { Routes } from '@angular/router';
import { authGuard } from './guards/auth.guard';
import { adminGuard } from './guards/admin.guard';

export const routes: Routes = [
  {
    path: '',
    redirectTo: '/dashboard',
    pathMatch: 'full'
  },
  {
    path: 'login',
    loadComponent: () => import('./pages/login/login.component').then(m => m.LoginComponent)
  },
  {
    path: 'register',
    loadComponent: () => import('./pages/register/register.component').then(m => m.RegisterComponent)
  },
  {
    path: 'dashboard',
    loadComponent: () => import('./pages/dashboard/dashboard.component').then(m => m.DashboardComponent),
    canActivate: [authGuard]
  },
  {
    path: 'labels',
    loadComponent: () => import('./pages/labels-list/labels-list.component').then(m => m.LabelsListComponent),
    canActivate: [authGuard]
  },
  {
    path: 'print-jobs',
    loadComponent: () => import('./pages/print-jobs/print-jobs.component').then(m => m.PrintJobsComponent),
    canActivate: [authGuard]
  },
  {
    path: 'audit',
    loadComponent: () => import('./pages/audit/audit.component').then(m => m.AuditComponent),
    canActivate: [authGuard]
  },
  {
    path: 'admin',
    loadComponent: () => import('./pages/print-labels/print-labels.component').then(m => m.PrintLabelsComponent),
    canActivate: [authGuard, adminGuard]
  },
  {
    path: 'profile',
    loadComponent: () => import('./pages/profile/profile.component').then(m => m.ProfileComponent),
    canActivate: [authGuard]
  },
  {
    path: 'font-test',
    loadComponent: () => import('./components/font-test/font-test.component').then(m => m.FontTestComponent)
  },
  {
    path: '**',
    redirectTo: '/dashboard'
  }
]; 