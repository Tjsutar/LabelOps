import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { DashboardService, DashboardStats } from '../../services/dashboard.service';
import { ToastService } from '../../services/toast.service';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, RouterLink],
  template: `
    <div class="p-5 max-w-7xl mx-auto">
      <div class="text-center mb-10">
        <h1 class="text-4xl font-bold text-gray-800 mb-3">Welcome to LabelOps</h1>
        <p class="text-lg text-gray-600">Manage your TMT bar labels efficiently</p>
        <p class="text-sm text-gray-500" *ngIf="lastUpdated">Last updated: {{ lastUpdated | date:'medium' }}</p>
      </div>

      <!-- Loading State -->
      <div *ngIf="loading" class="text-center py-10">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <p class="mt-2 text-gray-600">Loading dashboard data...</p>
      </div>

      <!-- Main Dashboard Content -->
      <div *ngIf="!loading && dashboardStats">
        <!-- Overview Cards -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-5 mb-8">
          <div class="bg-white p-5 rounded-lg shadow-md text-center border-l-4 border-blue-500">
            <h3 class="text-gray-600 mb-2 text-sm font-medium">Total Labels</h3>
            <p class="text-3xl font-bold text-blue-600 m-0">{{ dashboardStats.overview.total_labels | number }}</p>
          </div>
          <div class="bg-white p-5 rounded-lg shadow-md text-center border-l-4 border-green-500">
            <h3 class="text-gray-600 mb-2 text-sm font-medium">Printed Labels</h3>
            <p class="text-3xl font-bold text-green-600 m-0">{{ dashboardStats.overview.printed_labels | number }}</p>
          </div>
          <div class="bg-white p-5 rounded-lg shadow-md text-center border-l-4 border-yellow-500">
            <h3 class="text-gray-600 mb-2 text-sm font-medium">Pending Labels</h3>
            <p class="text-3xl font-bold text-yellow-600 m-0">{{ dashboardStats.overview.pending_labels | number }}</p>
          </div>
          <div class="bg-white p-5 rounded-lg shadow-md text-center border-l-4 border-red-500">
            <h3 class="text-gray-600 mb-2 text-sm font-medium">Failed Labels</h3>
            <p class="text-3xl font-bold text-red-600 m-0">{{ dashboardStats.overview.failed_labels | number }}</p>
          </div>
          <div class="bg-white p-5 rounded-lg shadow-md text-center border-l-4 border-purple-500">
            <h3 class="text-gray-600 mb-2 text-sm font-medium">Duplicate Labels</h3>
            <p class="text-3xl font-bold text-purple-600 m-0">{{ dashboardStats.overview.duplicate_labels | number }}</p>
          </div>
        </div>

        <!-- Activity & Performance Row -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          <!-- Recent Activity -->
          <div class="bg-white p-6 rounded-lg shadow-md">
            <h2 class="text-xl font-bold text-gray-800 mb-4">Recent Activity</h2>
            <div class="space-y-4">
              <div class="flex justify-between items-center p-3 bg-blue-50 rounded">
                <span class="text-gray-700">Labels Created (24h)</span>
                <span class="font-bold text-blue-600">{{ dashboardStats.activity.recent_labels_24h | number }}</span>
              </div>
              <div class="flex justify-between items-center p-3 bg-green-50 rounded">
                <span class="text-gray-700">Active Users (7d)</span>
                <span class="font-bold text-green-600">{{ dashboardStats.activity.active_users_7d | number }}</span>
              </div>
            </div>
          </div>

          <!-- Performance Metrics -->
          <div class="bg-white p-6 rounded-lg shadow-md">
            <h2 class="text-xl font-bold text-gray-800 mb-4">Performance</h2>
            <div class="space-y-4">
              <div class="flex justify-between items-center p-3 bg-indigo-50 rounded">
                <span class="text-gray-700">Print Success Rate</span>
                <span class="font-bold text-indigo-600">{{ dashboardStats.performance.print_success_rate }}</span>
              </div>
              <div class="flex justify-between items-center p-3 bg-gray-50 rounded">
                <span class="text-gray-700">Total Print Jobs</span>
                <span class="font-bold text-gray-600">{{ dashboardStats.performance.total_print_jobs | number }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Breakdown Charts Row -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          <!-- By Grade -->
          <div class="bg-white p-6 rounded-lg shadow-md">
            <h2 class="text-xl font-bold text-gray-800 mb-4">Labels by Grade</h2>
            <div class="space-y-2" *ngIf="getObjectKeys(dashboardStats.breakdown.by_grade).length > 0; else noGradeData">
              <div *ngFor="let grade of getObjectKeys(dashboardStats.breakdown.by_grade)" 
                   class="flex justify-between items-center p-2 hover:bg-gray-50 rounded">
                <span class="text-gray-700">{{ grade }}</span>
                <span class="font-semibold text-blue-600">{{ dashboardStats.breakdown.by_grade[grade] | number }}</span>
              </div>
            </div>
            <ng-template #noGradeData>
              <p class="text-gray-500 text-center py-4">No grade data available</p>
            </ng-template>
          </div>

          <!-- By Section -->
          <div class="bg-white p-6 rounded-lg shadow-md">
            <h2 class="text-xl font-bold text-gray-800 mb-4">Labels by Section</h2>
            <div class="space-y-2" *ngIf="getObjectKeys(dashboardStats.breakdown.by_section).length > 0; else noSectionData">
              <div *ngFor="let section of getObjectKeys(dashboardStats.breakdown.by_section)" 
                   class="flex justify-between items-center p-2 hover:bg-gray-50 rounded">
                <span class="text-gray-700">{{ section }}</span>
                <span class="font-semibold text-green-600">{{ dashboardStats.breakdown.by_section[section] | number }}</span>
              </div>
            </div>
            <ng-template #noSectionData>
              <p class="text-gray-500 text-center py-4">No section data available</p>
            </ng-template>
          </div>
        </div>

        <!-- Quick Actions -->
        <div class="bg-white p-6 rounded-lg shadow-md">
          <h2 class="text-2xl font-bold text-gray-800 mb-5">Quick Actions</h2>
          <div class="flex gap-4 flex-wrap">
            <a routerLink="/labels" 
               class="px-6 py-3 bg-blue-600 text-white text-decoration-none border-none rounded cursor-pointer text-sm transition-colors hover:bg-blue-700">
              <span>View All Labels</span>
            </a>
            <a routerLink="/print-jobs" 
               class="px-6 py-3 bg-green-600 text-white text-decoration-none border-none rounded cursor-pointer text-sm transition-colors hover:bg-green-700">
              <span>Print Jobs</span>
            </a>
            <a routerLink="/audit" 
               class="px-6 py-3 bg-purple-600 text-white text-decoration-none border-none rounded cursor-pointer text-sm transition-colors hover:bg-purple-700">
              <span>Audit Logs</span>
            </a>
            <button 
              (click)="refreshStats()" 
              [disabled]="loading"
              class="px-6 py-3 bg-gray-600 text-white border-none rounded cursor-pointer text-sm transition-colors hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed">
              <span>{{ loading ? 'Refreshing...' : 'Refresh Stats' }}</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Error State -->
      <div *ngIf="!loading && !dashboardStats" class="text-center py-10">
        <div class="text-red-500 mb-4">
          <svg class="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
        </div>
        <h3 class="text-xl font-semibold text-gray-800 mb-2">Failed to Load Dashboard</h3>
        <p class="text-gray-600 mb-4">There was an error loading the dashboard data.</p>
        <button 
          (click)="refreshStats()" 
          class="px-6 py-3 bg-blue-600 text-white border-none rounded cursor-pointer text-sm transition-colors hover:bg-blue-700">
          Try Again
        </button>
      </div>
    </div>
  `
})
export class DashboardComponent implements OnInit {
  dashboardStats: DashboardStats | null = null;
  loading = false;
  lastUpdated: Date | null = null;

  constructor(
    private dashboardService: DashboardService,
    private toastService: ToastService
  ) {}

  ngOnInit() {
    this.loadStats();
  }

  loadStats() {
    this.loading = true;
    this.dashboardService.getDashboardStats().subscribe({
      next: (stats) => {
        this.dashboardStats = stats;
        this.lastUpdated = new Date();
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading dashboard stats:', error);
        this.toastService.error('Failed to load dashboard statistics');
        this.loading = false;
      }
    });
  }

  refreshStats() {
    this.loadStats();
  }

  getObjectKeys(obj: any): string[] {
    return Object.keys(obj || {});
  }
}