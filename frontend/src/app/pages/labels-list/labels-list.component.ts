import { Component, OnInit } from "@angular/core";
import { CommonModule } from "@angular/common";
import { FormsModule } from "@angular/forms";
import { LabelService } from "../../services/label.service";
import { Label, LabelData } from "../../models/label.model";
import { LabelComponent } from "../../components/label/label.component";
import { ToastService } from "../../services/toast.service";
import { PrintLabelsComponent } from "../print-labels/print-labels.component";

@Component({
  selector: "app-labels-list",
  standalone: true,
  imports: [CommonModule, FormsModule, LabelComponent, PrintLabelsComponent],
  templateUrl: "./labels-list.component.html",
})
export class LabelsListComponent implements OnInit {
  labels: Label[] = [];
  loading = true;
  error: string | null = null;
  selectedLabel: Label | LabelData | null = null;
  showLabelPreview = false;
  exporting = false;
  // Pagination state
  pageSize = 8;
  offset = 0;
  totalCount = 0;
  loadingMore = false;

  constructor(
    private labelService: LabelService,
    private toastService: ToastService
  ) {}

  ngOnInit() {
    this.loadLabels();
  }

  loadLabels(reset: boolean = true) {
    if (reset) {
      this.loading = true;
      this.offset = 0;
      this.labels = [];
    } else {
      this.loadingMore = true;
    }
    this.error = null;

    const limit = this.pageSize && this.pageSize > 0 ? this.pageSize : undefined;
    this.labelService.getLabels({ limit, offset: this.offset }).subscribe({
      next: (response: { labels: Label[]; count: number }) => {
        const chunk = Array.isArray(response?.labels) ? response.labels : [];
        this.totalCount = Number(response?.count ?? (this.offset + chunk.length));

        // Append new chunk
        this.labels = [...this.labels, ...chunk];
        this.offset += chunk.length;

        this.loading = false;
        this.loadingMore = false;
      },
      error: () => {
        this.error = "Failed to load labels. Please try again.";
        this.loading = false;
        this.loadingMore = false;
      },
    });
  }


  printLabel(labelId: string | undefined) {
    if (!labelId) {
      this.toastService.error("Label ID is missing! Cannot print label.");
      return;
    }

    console.log("Frontend: Attempting to print label with ID:", labelId);
    console.log("Frontend: Label ID type:", typeof labelId);
    console.log("Frontend: Label ID length:", labelId?.length);

    this.labelService.printLabel(labelId).subscribe({
      next: (response) => {
        console.log("Print job created:", response);
        this.toastService.success(
          `Print job created successfully! Job ID: ${response.print_job_id.substring(
            0,
            8
          )}...`
        );
        // Refresh the labels to update status
        this.loadLabels();
      },
      error: (err) => {
        console.error("Error printing label:", err);
        console.error("Error details:", err.error);
        this.toastService.error(
          "Failed to create print job. Please try again."
        );
      },
    });
  }

  printSelectedLabel() {
    const anyLabel: any = this.selectedLabel as any;
    const id: string | undefined = anyLabel?.id || anyLabel?.ID;
    if (id) {
      this.printLabel(id);
      this.closeLabelPreview();
    } else {
      this.toastService.error("Label ID is missing! Cannot print label.");
    }
  }

  previewLabel(label: Label | LabelData) {
    // Use the row data directly; LabelComponent will normalize it
    this.selectedLabel = label;
    this.showLabelPreview = true;
  }

  closeLabelPreview() {
    this.showLabelPreview = false;
    this.selectedLabel = null;
  }

  exportLabels() {
    this.exporting = true;

    this.labelService.exportLabelsCSV().subscribe({
      next: (blob) => {
        console.log(blob);
        const filename = `labels_${new Date().toISOString().split("T")[0]}.csv`;
        this.labelService.downloadCSV(blob, filename);
        this.toastService.success("Labels exported successfully!");
        this.exporting = false;
      },
      error: (error) => {
        console.error("Error exporting labels:", error);
        this.toastService.error("Failed to export labels");
        this.exporting = false;
      },
    });
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

  canLoadMore(): boolean {
    return this.offset < this.totalCount && !this.loading && !this.loadingMore;
  }

  loadMore() {
    if (!this.canLoadMore()) return;
    this.loadLabels(false);
  }

 
}
