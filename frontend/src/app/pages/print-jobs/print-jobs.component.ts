import { Component, OnInit } from "@angular/core";
import { CommonModule } from "@angular/common";
import { LabelService } from "../../services/label.service";
import { ToastService } from "../../services/toast.service";
import { PrinterService } from "src/app/services/printer.service";
import { SearchByIdComponent } from "src/app/pages/print-jobs/search-by.component";

interface PrintJob {
  id: string;
  label_id: string;
  user_id: string;
  actual_label_id: string;
  heat_no: string;
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
  imports: [CommonModule, SearchByIdComponent],
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
  // Limit max entries shown on UI
  maxEntries = 8;

  constructor(
    private labelService: LabelService,
    private toastService: ToastService,
    private printerService: PrinterService
  ) {}

  ngOnInit() {
    this.loadPrintJobs();

  }

  loadPrintJobs() {
    this.loading = true;
    this.error = null;

    this.labelService.getPrintJobs().subscribe({
      next: (response: any) => {
        const list = response.print_jobs || response;
        this.printJobs = Array.isArray(list) ? list.slice(0, this.maxEntries) : [];
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
        this.toastService.success("Print job retry initiated successfully!");
        
        this.loadPrintJobs();
        console.log(response);
        this.getPrinter();
        this.printerService.printZPL(response.zpl_content);
        // this.viewZPL(response);
      },
      error: (err) => {
        console.error("Error retrying print job:", err);
        this.toastService.error("Failed to retry print job. Please try again.");
      },
    });
  }

  viewZPL(job: PrintJob) {
    // this.selectedZPL = job.zpl_content;
    // this.showZPLModal = true;
    // this.getPrinter();
    // this.printerService.printZPL(job.zpl_content);
    // console.log(job.zpl_content);
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

  getPrinter() {
    this.printerService.getDefaultPrinter().subscribe({
      next: (printer) => {
        console.log("Default printer:", printer);
      },
      error: (error) => {
        console.error("Error getting default printer:", error);
      },
    });
  }

  searchById(id: string) {
    this.labelService.getPrintJobById(id).subscribe({
      next: (response: any) => {
        // if API returns single object, wrap in array so table renders correctly
        const list = Array.isArray(response) ? response : [response];
        this.printJobs = list.slice(0, this.maxEntries);
        this.toastService.success('Print job found successfully!');
      },
      error: (error) => {
        this.toastService.error('Error getting print job');
        console.error('Error getting print job:', error);
      }
    });
  }

  searchByHeatNo(heatNo: string) {
    this.labelService.getPrintJobByHeatNo(heatNo).subscribe({
      next: (response: any) => {
        // if API returns single object, wrap in array so table renders correctly
        const list = Array.isArray(response) ? response : [response];
        this.printJobs = list.slice(0, this.maxEntries);
        this.toastService.success('Print job found successfully!');
      },
      error: (error) => {
        this.toastService.error('Error getting print job');
        console.error('Error getting print job:', error);
      }
    });
  }
}
