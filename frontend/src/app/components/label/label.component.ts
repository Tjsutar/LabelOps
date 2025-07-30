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

export interface LabelData {
  UUID?: string;
  HEAT_NO?: string;
  ID?: string;
  PRODUCT_HEADING?: string;
  SECTION?: string;
  GRADE?: string;
  ISI_TOP?: string;
  ISI_BOTTOM?: string;
  MILL?: string;
  DATE1?: string;
  TIME1?: string;
  LENGTH?: string;
  // Add other properties as needed
}

@Component({
  selector: 'app-label',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div
      class="w-[384px] h-[480px] self-auto border border-black rounded-lg p-2 m-2 bg-white"
    >
      <!-- Header -->
      <div class="flex justify-between items-center">
        <!-- MADE IN INDIA -->
        <div
          class="flex flex-col items-center justify-center text-center font-bold label-condensed border-2 border-black px-2 py-1 w-20 h-20 leading-tight"
        >
          <div>MADE</div>
          <div>IN</div>
          <div>INDIA</div>
        </div>

        <!-- Logo + Text -->
        <div class="text-right leading-tight flex flex-col items-center justify-center">
          <div class="flex items-center space-x-3">
            <div class="flex-shrink-0">
              <img
                src="assets/images/SAIL_LOGO.png"
                alt="ISI Mark"
                class="w-10 h-10 filter grayscale contrast-200 brightness-0"
              />
            </div>
            <div>
              <div class="text-[14px] font-[700]">
                स्टील अथॉरिटी ऑफ इंडिया लिमिटेड
              </div>
              <div class="font-bold text-[12px]">
                STEEL AUTHORITY OF INDIA LIMITED
              </div>
            </div>
          </div>
          <div class="italic label-eb-garamond text-sm self-end">
            भिलाई इस्पात संयंत्र
          </div>
          <div class="font-bold text-sm self-end">BHILAI STEEL PLANT</div>
        </div>
      </div>

      <!-- Section Title -->
      <div class="text-center font-bold text-[20px] border-2 border-black mt-2 mb-1">
        {{ labelData?.PRODUCT_HEADING || 'TMT BAR' }}
      </div>

      <!-- Label Info & QR -->
      <div class="flex justify-around w-auto h-auto bg-white text-black">
        <div>
          <!-- HEAT NO. -->
          <div class="font-bold text-[24px] label-eb-garamond leading-none">
            <div class="font-bold leading-none">HEAT NO.</div>
            <div class="font-bold leading-none">{{ labelData?.HEAT_NO || 'C103262' }}</div>
          </div>

          <!-- SECTION -->
          <div class="mt-3 text-[16px] label-eb-garamond leading-none">
            <div class="font-bold leading-none">SECTION</div>
            <div class="font-bold leading-none">{{ labelData?.SECTION || 'TMT BAR 25' }}</div>
          </div>

          <!-- GRADE -->
          <div class="mt-3 leading-none label-eb-garamond">
            <div class="font-bold text-[22px]">GRADE</div>
            <div class="font-bold">{{ labelData?.GRADE || 'IS 1786 FE550D' }}</div>
          </div>

          <!-- ID -->
          <div class="mt-3 leading-[16px] label-eb-garamond">
            <div class="font-bold">ID</div>
            <div class="font-bold">{{ labelData?.ID || '2025014374' }}</div>
          </div>
        </div>

        <!-- ISI logo and info -->
        <div class="flex flex-col leading-none label-eb-garamond mt-1.5">
          <div class="font-bold ml-0.5 text-[11px]">{{ labelData?.ISI_TOP || 'IS 1786:2008' }}</div>
          <div class="flex items-center justify-center ml-1 w-14 h-10">
            <img
              src="https://upload.wikimedia.org/wikipedia/commons/thumb/e/e8/Isi_mark.svg/1200px-Isi_mark.svg.png"
              alt="ISI Mark"
            />
          </div>
          <div class="font-bold ml-1 text-[11px]">{{ labelData?.ISI_BOTTOM || 'CML 187244' }}</div>
        </div>

        <!-- First QR Code -->
        <div>
          <canvas
            #canvas1
            class="h-32 w-32 mt-2 object-contain "
          ></canvas>
        </div>
      </div>

      <!-- Spacer -->
      <div class="h-0 mt-2 border border-gray-800"></div>

      <!-- Footer Info -->
      <div class="flex justify-between space-x-2">
        <!-- Second QR Code -->
        <div class="mt-2">
          <canvas
            #canvas2
            class="h-32 w-32 object-contain "
          ></canvas>
        </div>

        <!-- Footer Text Info -->
        <div class="flex flex-col w-52 mt-2">
          <div class="text-xl mb-2 mt-1 font-bold text-center">{{ labelData?.MILL || 'MM' }}</div>
          <div
            class="ml-5 text-sm grid grid-cols-[auto_minmax(0,1fr)] gap-x-2 font-bold label-condensed"
          >
            <div>LENGTH</div>
            <div>: {{ labelData?.LENGTH || 'STD' }}</div>
            <div>DATE</div>
            <div>: {{ labelData?.DATE1 || '25-JUN-25' }}</div>
            <div>TIME</div>
            <div>: {{ labelData?.TIME1 || '04:48' }}</div>
          </div>
        </div>
      </div>
    </div>
  `,
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
      this.qrUrl = `https://madeinindia.qcin.org/product-details/${this.labelData.UUID || 'default'}/MM_${this.labelData.HEAT_NO || 'default'}_${this.labelData.ID || 'default'}`;
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