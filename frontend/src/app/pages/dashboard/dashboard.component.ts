import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { LabelService } from '../../services/label.service';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, RouterLink],
  template: `
    <div class="p-5 max-w-7xl mx-auto">
      <div class="text-center mb-10">
        <h1 class="text-4xl font-bold text-gray-800 mb-3">Welcome to LabelOps</h1>
        <p class="text-lg text-gray-600">Manage your TMT bar labels efficiently</p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-5 mb-10">
        <div class="bg-white p-5 rounded-lg shadow-md text-center">
          <h3 class="text-gray-600 mb-3 text-base">Total Labels</h3>
          <p class="text-3xl font-bold text-blue-600 m-0">{{ totalLabels }}</p>
        </div>
        <div class="bg-white p-5 rounded-lg shadow-md text-center">
          <h3 class="text-gray-600 mb-3 text-base">Printed Today</h3>
          <p class="text-3xl font-bold text-blue-600 m-0">{{ printedToday }}</p>
        </div>
        <div class="bg-white p-5 rounded-lg shadow-md text-center">
          <h3 class="text-gray-600 mb-3 text-base">Pending</h3>
          <p class="text-3xl font-bold text-blue-600 m-0">{{ pendingLabels }}</p>
        </div>
      </div>

      <div class="bg-white p-5 rounded-lg shadow-md">
        <h2 class="text-2xl font-bold text-gray-800 mb-5">Quick Actions</h2>
        <div class="flex gap-4 flex-wrap">
          <a routerLink="/labels" 
             class="px-6 py-3 bg-blue-600 text-white text-decoration-none border-none rounded cursor-pointer text-sm transition-colors hover:bg-blue-700">
            <span>View All Labels</span>
          </a>
          <button 
            (click)="refreshStats()" 
            class="px-6 py-3 bg-blue-600 text-white text-decoration-none border-none rounded cursor-pointer text-sm transition-colors hover:bg-blue-700">
            <span>Refresh Stats</span>
          </button>
        </div>
      </div>
    </div>
  `
})
export class DashboardComponent implements OnInit {
  totalLabels = 0;
  printedToday = 0;
  pendingLabels = 0;

  constructor(private labelService: LabelService) {}

  ngOnInit() {
    this.loadStats();
  }

  loadStats() {
    // For now, we'll set some default values
    // In a real app, you'd fetch these from your backend
    this.totalLabels = 5; // Based on your dummy data
    this.printedToday = 2;
    this.pendingLabels = 3;
  }

  refreshStats() {
    this.loadStats();
  }
} 