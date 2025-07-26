import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterLink, RouterLinkActive } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { User } from '../../models/user.model';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [CommonModule, RouterLink, RouterLinkActive],
  template: `
    <header class="bg-white shadow-sm border-b border-gray-200">
      <div class="container mx-auto px-4">
        <div class="flex justify-between items-center h-16">
          <!-- Logo and Navigation -->
          <div class="flex items-center space-x-8">
            <div class="flex items-center">
              <h1 class="text-xl font-bold text-gray-900">LabelOps</h1>
            </div>
            
            <nav class="hidden md:flex space-x-6" *ngIf="currentUser">
              <a 
                routerLink="/dashboard" 
                routerLinkActive="text-primary-600 border-b-2 border-primary-600"
                class="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium transition-colors">
                Dashboard
              </a>
              <a 
                routerLink="/labels" 
                routerLinkActive="text-primary-600 border-b-2 border-primary-600"
                class="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium transition-colors">
                Labels
              </a>
              <a 
                routerLink="/audit" 
                routerLinkActive="text-primary-600 border-b-2 border-primary-600"
                class="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium transition-colors">
                Audit Logs
              </a>
              <a 
                routerLink="/admin" 
                routerLinkActive="text-primary-600 border-b-2 border-primary-600"
                class="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium transition-colors"
                *ngIf="currentUser?.role === 'admin'">
                Admin
              </a>
            </nav>
          </div>

          <!-- User Menu -->
          <div class="flex items-center space-x-4" *ngIf="currentUser; else loginButton">
            <div class="relative">
              <button 
                (click)="toggleUserMenu()"
                class="flex items-center space-x-2 text-sm text-gray-700 hover:text-gray-900 focus:outline-none">
                <div class="w-8 h-8 bg-primary-600 rounded-full flex items-center justify-center">
                  <span class="text-white font-medium">
                    {{ currentUser.first_name.charAt(0) }}{{ currentUser.last_name.charAt(0) }}
                  </span>
                </div>
                <span class="hidden md:block">{{ currentUser.first_name }} {{ currentUser.last_name }}</span>
              </button>

              <!-- Dropdown Menu -->
              <div 
                *ngIf="showUserMenu"
                class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-1 z-50 border border-gray-200">
                <a 
                  routerLink="/profile"
                  class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">
                  Profile
                </a>
                <button 
                  (click)="logout()"
                  class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">
                  Logout
                </button>
              </div>
            </div>
          </div>

          <ng-template #loginButton>
            <div class="flex items-center space-x-4">
              <a 
                routerLink="/login"
                class="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium">
                Login
              </a>
              <a 
                routerLink="/register"
                class="btn-primary">
                Register
              </a>
            </div>
          </ng-template>
        </div>
      </div>
    </header>
  `,
  styles: []
})
export class HeaderComponent {
  currentUser: User | null = null;
  showUserMenu = false;

  constructor(
    private authService: AuthService,
    private router: Router
  ) {
    this.authService.currentUser$.subscribe(user => {
      this.currentUser = user;
    });
  }

  toggleUserMenu(): void {
    this.showUserMenu = !this.showUserMenu;
  }

  logout(): void {
    this.authService.logout();
    this.showUserMenu = false;
    this.router.navigate(['/login']);
  }
} 