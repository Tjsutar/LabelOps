import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LabelService } from '../../services/label.service';
import { Label } from '../../models/label.model';
import { LabelComponent, LabelData } from '../../components/label/label.component';
import { ToastService } from '../../services/toast.service';

@Component({
  selector: 'app-labels-list',
  standalone: true,
  imports: [CommonModule, LabelComponent],
  templateUrl: './labels-list.component.html'
})
export class LabelsListComponent implements OnInit {
  labels: Label[] = [];
  loading = true;
  error: string | null = null;
  selectedLabel: LabelData | null = null;
  showLabelPreview = false;

  constructor(
    private labelService: LabelService,
    private toastService: ToastService
  ) {}

  ngOnInit() {
    this.loadLabels();
  }

  loadLabels() {
    this.loading = true;
    this.error = null;
    
    this.labelService.getLabels().subscribe({
      next: (response: any) => {
        this.labels = response.labels || response;
        console.log('Frontend: Loaded labels:', this.labels);
        if (this.labels.length > 0) {
          console.log('Frontend: First label structure:', this.labels[0]);
          console.log('Frontend: First label id field:', this.labels[0].id);
          console.log('Frontend: First label label_id field:', this.labels[0].label_id);
        }
        this.loading = false;
      },
      error: (err) => {
        this.error = 'Failed to load labels. Please try again.';
        this.loading = false;
        console.error('Error loading labels:', err);
      }
    });
  }

  printLabel(labelId: string | undefined) {
    if (!labelId) {
      this.toastService.error('Label ID is missing! Cannot print label.');
      return;
    }
    
    console.log('Frontend: Attempting to print label with ID:', labelId);
    console.log('Frontend: Label ID type:', typeof labelId);
    console.log('Frontend: Label ID length:', labelId?.length);
    
    this.labelService.printLabel(labelId).subscribe({
      next: (response) => {
        console.log('Print job created:', response);
        this.toastService.success(`Print job created successfully! Job ID: ${response.print_job_id.substring(0, 8)}...`);
        // Refresh the labels to update status
        this.loadLabels();
      },
      error: (err) => {
        console.error('Error printing label:', err);
        console.error('Error details:', err.error);
        this.toastService.error('Failed to create print job. Please try again.');
      }
    });
  }

  printSelectedLabel() {
    if (this.selectedLabel?.ID) {
      this.printLabel(this.selectedLabel.ID);
      this.closeLabelPreview();
    }
  }

  previewLabel(label: any) {
    // Convert the label data to the format expected by the label component
    this.selectedLabel = {
      ID: label.id || label.ID,
      HEAT_NO: label.HEAT_NO || label.heat_no,
      PRODUCT_HEADING: label.PRODUCT_HEADING || label.product_heading,
      SECTION: label.SECTION || label.section,
      GRADE: label.GRADE || label.grade,
      ISI_TOP: label.ISI_TOP || label.isi_top,
      ISI_BOTTOM: label.ISI_BOTTOM || label.isi_bottom,
      MILL: label.MILL || label.mill,
      DATE1: label.DATE1 || label.date1,
      TIME1: label.TIME1 || label.time1,
      LENGTH: label.LENGTH || label.length,
      label_id: label.label_id
    };
    this.showLabelPreview = true;
  }

  closeLabelPreview() {
    this.showLabelPreview = false;
    this.selectedLabel = null;
  }
}
