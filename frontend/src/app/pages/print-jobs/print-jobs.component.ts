import { Component, OnInit } from "@angular/core";
import { CommonModule } from "@angular/common";
import { LabelService } from "../../services/label.service";
import { ToastService } from "../../services/toast.service";

interface PrintJob {
  id: string;
  label_id: string;
  user_id: string;
  status: string;
  zpl_content: string;
  max_retries: number;
  retry_count: number;
  error_message?: string;
  created_at: string;
  updated_at: string;
}

@Component({
  selector: "app-print-jobs",
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="p-5 max-w-7xl mx-auto">
      <!-- Header -->
      <div class="flex justify-between items-center mb-5">
        <h2 class="text-2xl font-bold text-gray-800">Print Jobs</h2>
        <div class="flex gap-3">
          <button
            (click)="exportPrintJobs()"
            [disabled]="loading || exporting"
            class="px-4 py-2 bg-green-600 text-white rounded text-sm hover:bg-green-700 disabled:bg-gray-500 disabled:cursor-not-allowed"
          >
            {{ exporting ? "Exporting..." : "Export CSV" }}
          </button>
          <button
            (click)="loadPrintJobs()"
            [disabled]="loading"
            class="px-4 py-2 bg-blue-600 text-white rounded text-sm hover:bg-blue-700 disabled:bg-gray-500 disabled:cursor-not-allowed"
          >
            {{ loading ? "Loading..." : "Refresh" }}
          </button>
        </div>
      </div>

      <!-- Loading -->
      <div *ngIf="loading" class="text-center py-10 text-gray-600">
        <p>Loading print jobs...</p>
      </div>

      <!-- Error -->
      <div
        *ngIf="error"
        class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-5"
      >
        <p>{{ error }}</p>
        <button
          (click)="loadPrintJobs()"
          class="bg-red-600 text-white px-3 py-1 rounded mt-2 hover:bg-red-700"
        >
          Try Again
        </button>
      </div>

      <!-- Table -->
      <div
        *ngIf="!loading && !error"
        class="bg-white rounded-lg shadow-md overflow-hidden"
      >
        <div class="overflow-x-auto">
          <table class="w-full border-collapse">
            <thead
              class="sticky top-0 z-10 bg-gray-100 text-xs font-bold text-gray-600 uppercase tracking-wider"
            >
              <tr>
                <th class="px-3 py-3 text-left">Job ID</th>
                <th class="px-3 py-3 text-left">Label ID</th>
                <th class="px-3 py-3 text-left">Status</th>
                <th class="px-3 py-3 text-left">Retries</th>
                <th class="px-3 py-3 text-left">Created</th>
                <th class="px-3 py-3 text-left">Updated</th>
                <th class="px-3 py-3 text-left">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr
                *ngFor="let job of printJobs"
                class="border-b border-gray-200 hover:bg-gray-50 text-xs"
              >
                <td class="px-3 py-3">
                  {{ job.id.substring(0, 8) }}...
                </td>
                <td class="px-3 py-3">{{ job.label_id }}</td>
                <td class="px-3 py-3">
                  <span
                    [class]="
                      'font-medium px-2 py-1 rounded text-xs ' +
                      getStatusClass(job.status)
                    "
                  >
                    {{ job.status }}
                  </span>
                </td>
                <td class="px-3 py-3">
                  {{ job.retry_count }}/{{ job.max_retries }}
                </td>
                <td class="px-3 py-3">
                  {{ formatDate(job.created_at) }}
                </td>
                <td class="px-3 py-3">
                  {{ formatDate(job.updated_at) }}
                </td>
                <td class="px-3 py-3">
                  <div class="flex gap-2 flex-wrap">
                    <button
                      *ngIf="
                        job.status === 'failed' &&
                        job.retry_count < job.max_retries
                      "
                      (click)="retryPrintJob(job.id)"
                      class="px-3 py-1 bg-yellow-600 text-white rounded text-xs hover:bg-yellow-700"
                    >
                      Retry
                    </button>
                    <button
                      (click)="viewZPL(job)"
                      class="px-3 py-1 bg-blue-600 text-white rounded text-xs hover:bg-blue-700"
                    >
                      View ZPL
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Empty -->
        <div
          *ngIf="printJobs.length === 0"
          class="text-center py-10 text-gray-600"
        >
          <p>No print jobs found.</p>
        </div>
      </div>

      <!-- ZPL Modal -->
      <div
        *ngIf="showZPLModal"
        class="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50"
        (click)="closeZPLModal()"
      >
        <div
          class="bg-white rounded-lg max-w-4xl max-h-screen overflow-auto relative"
          (click)="$event.stopPropagation()"
        >
          <div
            class="flex justify-between items-center p-5 border-b border-gray-200"
          >
            <h3 class="text-xl font-semibold text-gray-800">ZPL Content</h3>
            <button
              (click)="closeZPLModal()"
              class="text-2xl text-gray-600 hover:text-gray-800 w-8 h-8 flex items-center justify-center"
            >
              &times;
            </button>
          </div>
          <div class="p-5">
            <pre class="bg-gray-100 p-4 rounded text-sm overflow-x-auto">{{
              selectedZPL
            }}</pre>
          </div>
        </div>
      </div>
    </div>
  `,
  styles: [],
})
export class PrintJobsComponent implements OnInit {
  printJobs: PrintJob[] = [];
  loading = true;
  error: string | null = null;
  showZPLModal = false;
  selectedZPL = "";
  exporting = false;

  constructor(
    private labelService: LabelService,
    private toastService: ToastService
  ) {}

  ngOnInit() {
    this.loadPrintJobs();
  }

  loadPrintJobs() {
    this.loading = true;
    this.error = null;

    this.labelService.getPrintJobs().subscribe({
      next: (response: any) => {
        this.printJobs = response.print_jobs || response;
        this.loading = false;
      },
      error: (err) => {
        this.error = "Failed to load print jobs. Please try again.";
        this.loading = false;
        console.error("Error loading print jobs:", err);
      },
    });
  }

  retryPrintJob(jobId: string) {
    this.labelService.retryPrintJob(jobId).subscribe({
      next: (response) => {
        console.log("Print job retried:", response);
        this.toastService.success("Print job retry initiated successfully!");
        this.loadPrintJobs();
      },
      error: (err) => {
        console.error("Error retrying print job:", err);
        this.toastService.error("Failed to retry print job. Please try again.");
      },
    });
  }

  viewZPL(job: PrintJob) {
    this.selectedZPL = job.zpl_content;
    this.showZPLModal = true;
  }

  closeZPLModal() {
    this.showZPLModal = false;
    this.selectedZPL = "";
  }

  getStatusClass(status: string): string {
    switch (status) {
      case "completed":
        return "bg-green-100 text-green-800";
      case "failed":
        return "bg-red-100 text-red-800";
      case "pending":
        return "bg-yellow-100 text-yellow-800";
      case "processing":
        return "bg-blue-100 text-blue-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  }

  formatDate(dateString: string): string {
    return new Date(dateString).toLocaleString();
  }

  exportPrintJobs() {
    this.exporting = true;

    this.labelService.exportPrintJobsCSV().subscribe({
      next: (blob) => {
        console.log(blob);
        const filename = `print_jobs_${
          new Date().toISOString().split("T")[0]
        }.csv`;
        this.labelService.downloadCSV(blob, filename);
        this.toastService.success("Print jobs exported successfully!");
        this.exporting = false;
      },
      error: (error) => {
        console.error("Error exporting print jobs:", error);
        this.toastService.error("Failed to export print jobs");
        this.exporting = false;
      },
    });
  }
}
