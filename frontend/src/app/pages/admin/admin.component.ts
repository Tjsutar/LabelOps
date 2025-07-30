import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-admin',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="p-5 max-w-7xl mx-auto">
      <div class="bg-white p-8 rounded-lg shadow-md text-center">
        <h2 class="text-3xl font-bold text-gray-800 mb-3">Admin Panel</h2>
        <p class="text-gray-600">Admin functionality coming soon...</p>
      </div>
    </div>
  `
})
export class AdminComponent {} 