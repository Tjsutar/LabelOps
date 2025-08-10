import { Component, OnInit } from "@angular/core";
import { CommonModule } from "@angular/common";
import { LabelService } from "../../services/label.service";
import { ToastService } from "../../services/toast.service";

interface PrintJob {
  id: string;
  label_id: string;
  user_id: string;
  actual_label_id: string;
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
  templateUrl: "./print-jobs.html",
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
      case "success":
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
