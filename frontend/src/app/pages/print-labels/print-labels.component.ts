import { Component } from "@angular/core";
import { CommonModule } from "@angular/common";
import { HttpClient, HttpHeaders } from "@angular/common/http";

@Component({
  selector: "app-print-label",
  standalone: true,
  imports: [CommonModule],
  template: `
    <section class="bg-white shadow-md rounded-xl p-6 space-y-5 max-w-3xl mx-auto">
      <header class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4">
        <label
          for="jsonUpload"
          class="inline-flex items-center gap-3 cursor-pointer bg-gray-100 hover:bg-gray-200 transition px-4 py-2 rounded-lg text-gray-700 border border-dashed border-gray-300"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-blue-600" fill="none"
            viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M12 4v16m8-8H4" />
          </svg>
          <span>Select JSON File</span>
          <input
            id="jsonUpload"
            type="file"
            (change)="onFileSelected($event)"
            accept=".json"
            class="hidden"
          />
        </label>

        <button
          (click)="uploadJson()"
          [disabled]="!jsonData"
          class="inline-flex items-center justify-center px-5 py-2.5 bg-blue-600 text-white font-medium rounded-lg shadow hover:bg-blue-700 transition disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <svg *ngIf="jsonData" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2"
            viewBox="0 0 20 20" fill="currentColor">
            <path
              d="M3 4a1 1 0 011-1h3a1 1 0 010 2H5v11h10V5h-2a1 1 0 110-2h3a1 1 0 011 1v13a1 1 0 01-1 1H4a1 1 0 01-1-1V4z" />
            <path d="M9 8a1 1 0 012 0v5a1 1 0 11-2 0V8z" />
          </svg>
          Print Labels
        </button>
      </header>

      <div *ngIf="responseMessage" class="flex items-center gap-2 text-sm">
        <svg
          *ngIf="responseMessage.startsWith('✅')"
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5 text-green-600"
          viewBox="0 0 20 20"
          fill="currentColor"
        >
          <path
            fill-rule="evenodd"
            d="M16.707 5.293a1 1 0 010 1.414L8.414 15l-4.121-4.121a1 1 0 011.414-1.414L8.414 12.586l7.293-7.293a1 1 0 011.414 0z"
            clip-rule="evenodd"
          />
        </svg>
        <svg
          *ngIf="responseMessage.startsWith('❌')"
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5 text-red-600"
          viewBox="0 0 20 20"
          fill="currentColor"
        >
          <path
            fill-rule="evenodd"
            d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-10.707a1 1 0 00-1.414-1.414L10 8.586 7.707 6.293a1 1 0 00-1.414 1.414L8.586 10l-2.293 2.293a1 1 0 001.414 1.414L10 11.414l2.293 2.293a1 1 0 001.414-1.414L11.414 10l2.293-2.293z"
            clip-rule="evenodd"
          />
        </svg>
        <span [ngClass]="{
           'text-green-600 font-medium': responseMessage.startsWith('✅'),
          'text-red-600 font-medium': responseMessage.startsWith('❌')
        }">
          {{ responseMessage }}
        </span>
      </div>
    </section>
  `,
})
export class PrintLabelsComponent {
  jsonData: any = null;
  responseMessage = "";

  constructor(private http: HttpClient) {}

  onFileSelected(event: Event) {
    const file = (event.target as HTMLInputElement)?.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = () => {
      try {
        this.jsonData = JSON.parse(reader.result as string);
        this.responseMessage = "✅ JSON data loaded successfully";
      } catch (err) {
        this.responseMessage = "❌ Invalid JSON file";
      }
    };
    reader.readAsText(file);
  }

  uploadJson() {
    const token = localStorage.getItem("token");
    if (!token) {
      this.responseMessage = "No auth token found";
      return;
    }

    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    });

    this.http
      .post("http://localhost:8080/api/v1/labels/batch", this.jsonData, {
        headers,
      })
      .subscribe({
        next: () => (this.responseMessage = "✅ Labels printed successfully!"),
        error: (err) => {
          console.error(err);
          this.responseMessage = "❌ Submission failed";
        },
      });
  }
}
