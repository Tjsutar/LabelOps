import {
  Component,
  ElementRef,
  Input,
  ViewChild,
  AfterViewInit,
  OnChanges,
  SimpleChanges,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import * as QRCode from 'qrcode';
import { LabelData } from '../../models/label.model';

@Component({
  selector: 'app-label',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './label.html',
})
export class LabelComponent implements AfterViewInit, OnChanges {
  @Input() labelData!: LabelData;

  @Input() qrUrl: string = '';

  @ViewChild('canvas1', { static: false }) canvas1Ref!: ElementRef<HTMLCanvasElement>;
  @ViewChild('canvas2', { static: false }) canvas2Ref!: ElementRef<HTMLCanvasElement>;

  ngAfterViewInit(): void {
    this.generateQrCodes();
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['qrUrl'] && !changes['qrUrl'].firstChange) {
      this.generateQrCodes();
    }
    if (changes['labelData'] && !changes['labelData'].firstChange) {
      this.generateQrCodes();
    }
  }

  private async generateQrCodes() {
    // Generate QR URL if not provided
    if (!this.qrUrl && this.labelData) {
      this.qrUrl = `https://madeinindia.qcin.org/product-details/${this.labelData.ID || 'default'}/MM_${this.labelData.HEAT_NO || 'default'}_${this.labelData.ID || 'default'}`;
    }

    if (!this.qrUrl || !this.canvas1Ref || !this.canvas2Ref) {
      console.error('❌ Missing QR URL or canvas references');
      return;
    }

    try {
      await QRCode.toCanvas(this.canvas1Ref.nativeElement, this.qrUrl, {
        width: 128,
        margin: 1,
      });

      await QRCode.toCanvas(this.canvas2Ref.nativeElement, this.qrUrl, {
        width: 128,
        margin: 1,
      });

      console.log('✅ QR code rendered successfully');
    } catch (err) {
      console.error('❌ Failed to generate QR code', err);
    }
  }
} 