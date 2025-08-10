    // search-by-id.component.ts
import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { LabelService } from 'src/app/services/label.service';
import { ToastService } from 'src/app/services/toast.service';
@Component({
  selector: 'app-search-by-id',
  standalone: true,
  imports: [CommonModule, FormsModule ],
  template: `
    <section class="bg-white shadow-md rounded-xl p-6 max-w-md mx-auto">
      <div class="flex gap-3">
        <input
          [(ngModel)]="searchId"
          type="text"
          placeholder="Enter ID to search"
          class="flex-grow border border-gray-300 rounded-lg px-3 py-2"
        />
        <button
          (click)="search()"
          [disabled]="!searchId.trim()"
          class="inline-flex items-center justify-center px-4 py-2 bg-blue-600 text-white font-medium rounded-lg shadow hover:bg-blue-700 transition disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-5 w-5 mr-2"
            viewBox="0 0 20 20"
            fill="currentColor"
          >
            <path
              d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1111.196 3.243l4.328 4.329a1 1 0 01-1.415 1.415l-4.328-4.328A6 6 0 012 8z"
            />
          </svg>
          Search
        </button>
      </div>
    </section>
  `,
})
export class SearchByIdComponent {
  searchId: string = '';

  constructor(private controller: LabelService, private toastService: ToastService) {}

  search() {
    if (!this.searchId.trim()) {
      return;
    }

    console.log('Searching for ID:', this.searchId);
    // TODO: Add your search logic here, or emit event to parent component
    
    this.searchId = '';
    this.getPrintJobById(this.searchId);
    }

    getPrintJobById(id: string) {
      this.controller.getPrintJobById(this.searchId).subscribe({
        next: (response: any) => {
          this.toastService.success('Print job found successfully!');
            console.log('Print job:', response);
        },
        error: (error: any) => {
          this.toastService.error('Error getting print job');
          console.error('Error getting print job:', error);
        },
      });
    }
}
