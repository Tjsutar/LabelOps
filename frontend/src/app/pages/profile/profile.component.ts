import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../services/auth.service';
import { User } from '../../models/user.model';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="p-5 max-w-2xl mx-auto">
      <div class="bg-white p-8 rounded-lg shadow-md">
        <h2 class="text-3xl font-bold text-gray-800 mb-8 text-center">User Profile</h2>
        
        <div *ngIf="currentUser" class="space-y-5">
          <div class="flex justify-between items-center p-4 bg-gray-50 rounded-lg">
            <label class="font-semibold text-gray-700">Name:</label>
            <span class="text-gray-800">{{ currentUser.first_name }} {{ currentUser.last_name }}</span>
          </div>
          <div class="flex justify-between items-center p-4 bg-gray-50 rounded-lg">
            <label class="font-semibold text-gray-700">Email:</label>
            <span class="text-gray-800">{{ currentUser.email }}</span>
          </div>
          <div class="flex justify-between items-center p-4 bg-gray-50 rounded-lg">
            <label class="font-semibold text-gray-700">Role:</label>
            <span class="text-gray-800">{{ currentUser.role }}</span>
          </div>
          <div class="flex justify-between items-center p-4 bg-gray-50 rounded-lg">
            <label class="font-semibold text-gray-700">Status:</label>
            <span [class]="'font-medium ' + (currentUser.is_active ? 'text-green-600' : 'text-red-600')">
              {{ currentUser.is_active ? 'Active' : 'Inactive' }}
            </span>
          </div>
        </div>
        
        <div *ngIf="!currentUser" class="text-center py-10 text-gray-600">
          <p>Loading profile...</p>
        </div>
      </div>
    </div>
  `
})
export class ProfileComponent implements OnInit {
  currentUser: User | null = null;

  constructor(private authService: AuthService) {}

  ngOnInit() {
    this.authService.currentUser$.subscribe(user => {
      this.currentUser = user;
    });
  }
} 