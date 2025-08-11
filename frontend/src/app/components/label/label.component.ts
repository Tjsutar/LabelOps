import {
  Component,
  ElementRef,
  Input,
  ViewChild,
  AfterViewInit,
  OnChanges,
  SimpleChanges,
  OnInit
} from '@angular/core';
import { CommonModule } from '@angular/common';
import * as QRCode from 'qrcode';
import { Label, LabelData } from '../../models/label.model';
import { LabelService } from '../../services/label.service';

@Component({
  selector: 'app-label',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './label.html',
})
export class LabelComponent implements OnInit, AfterViewInit, OnChanges {
  /**
   * You can pass either labelId OR the full label object from the parent.
   * If label is provided, it will be used directly without fetching.
   */
  @Input() labelId?: string;
  @Input() label?: Label | LabelData;
  
  // Minimal shape used by template for display
  private mapToDisplayLabel(input: Label | LabelData): {
    id?: string;
    actual_label_id?: string;
    label_id?: string;
    heat_no?: string;
    section?: string;
    grade?: string;
    product_heading?: string;
    isi_top?: string;
    isi_bottom?: string;
    mill?: string;
    length?: number;
    date?: string;
    time?: string;
    unit?: string;
    weight?: string;
    location?: string;
    pqd?: string;
  } {
    const isLabel = (obj: any): obj is Label => 'id' in obj && ('product_heading' in obj || 'grade' in obj);
    if (isLabel(input)) {
      return {
        id: input.id,
        actual_label_id: input.actual_label_id,
        label_id: input.label_id,
        heat_no: input.heat_no,
        section: input.section,
        grade: input.grade,
        product_heading: input.product_heading,
        isi_top: input.isi_top,
        isi_bottom: input.isi_bottom,
        mill: input.mill,
        length: input.length,
        date: (input as any).date, // field name in Label interface
        time: input.time,
        unit: input.unit,
        weight: input.weight,
        location: input.location,
        pqd: input.pqd,
      };
    }
    const ld = input as LabelData;
    return {
      id: ld.ID,
      actual_label_id: ld.actual_label_id,
      label_id: ld.label_id,
      heat_no: ld.HEAT_NO,
      section: ld.SECTION,
      grade: ld.GRADE,
      product_heading: ld.PRODUCT_HEADING,
      isi_top: ld.ISI_TOP,
      isi_bottom: ld.ISI_BOTTOM,
      mill: ld.MILL,
      length: ld.LENGTH,
      date: ld.DATE1,
      time: ld.TIME,
      unit: ld.UNIT,
      weight: ld.WEIGHT,
      location: ld.LOCATION,
      pqd: ld.PQD,
    };
  }

  labelData!: ReturnType<typeof this.mapToDisplayLabel>;
  @Input() qrUrl: string = '';

  @ViewChild('canvas1', { static: false }) canvas1Ref!: ElementRef<HTMLCanvasElement>;
  @ViewChild('canvas2', { static: false }) canvas2Ref!: ElementRef<HTMLCanvasElement>;

  constructor(private labelService: LabelService) {}

  private viewInitialized = false;
  private qrRetry = 0;

  ngOnInit(): void {
    if (this.label) {
      // Parent provided the full label object
      this.labelData = this.mapToDisplayLabel(this.label);
    } else if (this.labelId) {
      // Parent provided only the ID — fetch from API
      this.loadLabelData(this.labelId);
    }
  }

  private buildSecondQrPayload(): string {
    const d = this.labelData || {} as any;
    const unit = d.unit || 'SAIL-BSP';
    const mill = d.mill || '';
    const heat = d.heat_no || '';
    const section = d.section || '';
    const grade = d.grade || '';
    const id = d.label_id || '';
    const length = (typeof d.length === 'number' && d.length > 0) ? String(d.length) : 'STD';
    const weight = d.weight || '';
    const location = d.location || '';
    const pqd = d.pqd || '';
    const date = d.date || '';
    const time = d.time || '';

    return `UNIT:${unit};MILL:${mill};HEAT:${heat};SECTION:${section};GRADE:${grade};ID:${id};LENGTH:${length};WEIGHT:${weight};LOCATION:${location};PQD:${pqd};DATE:${date};TIME:${time};`;
  }

  ngAfterViewInit(): void {
    this.viewInitialized = true;
    // Defer to next tick to ensure ViewChilds are attached
    if (this.labelData) {
      setTimeout(() => this.generateQrCodes(), 0);
    }
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['label'] && this.label) {
      this.labelData = this.mapToDisplayLabel(this.label);
      if (this.viewInitialized) {
        this.generateQrCodes();
      }
    } else if (changes['labelId'] && this.labelId) {
      this.loadLabelData(this.labelId);
    }
  }

  private loadLabelData(id: string) {
    this.labelService.getLabelById(id).subscribe({
      next: (data) => {
        this.labelData = this.mapToDisplayLabel(data);
        if (this.viewInitialized) {
          this.generateQrCodes();
        }
      },
      error: (err) => {
        console.error('❌ Failed to load label data', err);
      }
    });
  }

  private async generateQrCodes() {
    // Ensure data is present
    if (!this.labelData) {
      return;
    }

    // Build QR URL if not provided
    if (!this.qrUrl) {
      this.qrUrl = `https://madeinindia.qcin.org/product-details/${
        this.labelData.label_id || 'default'
      }/MM_${this.labelData.heat_no || 'default'}_${this.labelData.id || 'default'}`;
      console.log(this.labelData.label_id);
    }

    // Ensure canvas references exist and are attached
    const c1 = this.canvas1Ref?.nativeElement;
    const c2 = this.canvas2Ref?.nativeElement;
    if (!this.qrUrl || !c1 || !c2) {
      // Avoid noisy logs before view is initialized
      if (this.viewInitialized) {
        if (this.qrRetry < 3) {
          this.qrRetry++;
          setTimeout(() => this.generateQrCodes(), 50);
          return;
        }
        console.error('❌ Missing QR URL or canvas references');
      }
      return;
    }

    try {
      // First QR: external URL format
      await QRCode.toCanvas(c1, this.qrUrl, {
        width: 128,
        margin: 1,
      });
      // Second QR: formatted row data payload
      const payload = this.buildSecondQrPayload();
      await QRCode.toCanvas(c2, payload, {
        width: 128,
        margin: 1,
      });

      this.qrRetry = 0;
      console.log('✅ QR code rendered successfully');
    } catch (err) {
      console.error('❌ Failed to generate QR code', err);
    }
  }
}
