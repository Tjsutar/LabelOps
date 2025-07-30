import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LabelService } from '../../services/label.service';
import { Label } from '../../models/label.model';
import { LabelComponent, LabelData } from '../../components/label/label.component';

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

  constructor(private labelService: LabelService) {}

  ngOnInit() {
    this.loadLabels();
  }

  loadLabels() {
    this.loading = true;
    this.error = null;
    
    this.labelService.getLabels().subscribe({
      next: (response: any) => {
        this.labels = response.labels || response;
        this.loading = false;
      },
      error: (err) => {
        this.error = 'Failed to load labels. Please try again.';
        this.loading = false;
        console.error('Error loading labels:', err);
      }
    });
  }

  printLabel(labelId: string) {
    this.labelService.printLabel(labelId).subscribe({
      next: (response) => {
        console.log('Print job created:', response);
        // You can add a success message here
      },
      error: (err) => {
        console.error('Error printing label:', err);
        // You can add an error message here
      }
    });
  }

  previewLabel(label: any) {
    // Convert the label data to the format expected by the label component
    this.selectedLabel = {
      UUID: label.id || label.ID,
      HEAT_NO: label.HEAT_NO || label.heat_no,
      ID: label.ID || label.id,
      PRODUCT_HEADING: label.PRODUCT_HEADING || label.product_heading,
      SECTION: label.SECTION || label.section,
      GRADE: label.GRADE || label.grade,
      ISI_TOP: label.ISI_TOP || label.isi_top,
      ISI_BOTTOM: label.ISI_BOTTOM || label.isi_bottom,
      MILL: label.MILL || label.mill,
      DATE1: label.DATE1 || label.date1,
      TIME1: label.TIME1 || label.time1,
      LENGTH: label.LENGTH || label.length
    };
    this.showLabelPreview = true;
  }

  closeLabelPreview() {
    this.showLabelPreview = false;
    this.selectedLabel = null;
  }
}
