import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { DashboardService, DashboardStats } from '../../services/dashboard.service';
import { ToastService } from '../../services/toast.service';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './dashboard.component.html',
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