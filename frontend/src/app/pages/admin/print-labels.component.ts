import { Component } from "@angular/core";
import { CommonModule } from "@angular/common";
import { HttpClient, HttpHeaders } from "@angular/common/http";

@Component({
  selector: "app-print-label",
  standalone: true,
  imports: [CommonModule],
  template: `
    <div
      class="p-2 my-8 bg-white rounded shadow flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 text-sm"
    >
      <div class="flex items-center gap-2 flex-wrap">
        <input
          type="file"
          (change)="onFileSelected($event)"
          accept=".json"
          class="border rounded px-2 py-1 text-sm"
        />
        <p *ngIf="responseMessage" class="text-green-600 font-medium sm:ml-4">
          {{ responseMessage }}
        </p>
      </div>
      <button
        (click)="uploadJson()"
        class="bg-blue-600 text-white px-3 py-1.5 rounded hover:bg-blue-700 disabled:opacity-50"
        [disabled]="!jsonData"
      >
        Print Labels
      </button>
    </div>
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
    const token = localStorage.getItem("token"); // Or however you're storing it
    if (!token) {
      this.responseMessage = "❌ No auth token found";
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
