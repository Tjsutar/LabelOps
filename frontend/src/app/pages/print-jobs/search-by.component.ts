import { Component, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-search-by',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <section class="bg-white shadow-md rounded-xl p-3 space-y-1 max-w-3xl mx-auto mb-5">
      <div class="flex gap-4">
        <input
          [(ngModel)]="searchHeat"
          (keyup.enter)="onSearchClick()"
          type="text"
          placeholder="Enter Heat No to search"
          class="flex-grow border border-gray-300 rounded-lg px-3 py-2 mr-4"
        />
        <button
          (click)="onSearchClick()"
          [disabled]="!searchHeat.trim()"
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
  searchId = '';
  searchHeat ='';
  @Output() searchIdEvent = new EventEmitter<string>();
@Output() searchHeatnoEvent = new EventEmitter<string>();

  onSearchClick() {
    // if (this.searchId.trim()) {
    //   this.searchIdEvent.emit(this.searchId.trim());
    // }
    if (this.searchHeat.trim()){
      this.searchHeatnoEvent.emit((this.searchHeat.trim()))
    }

  }
}
