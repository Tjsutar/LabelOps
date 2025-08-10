import { Injectable } from '@angular/core';
import { Observable, from, of, throwError } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import * as ZplImage from 'zpl-image';


declare var BrowserPrint: any;

@Injectable({
  providedIn: 'root'
})
export class PrinterService {
  private browserPrint: any;

  constructor() {
    this.initializeBrowserPrint();
  }

  private initializeBrowserPrint(): void {
    // Initialize Zebra Browser Print SDK
    if (typeof BrowserPrint !== 'undefined') {
      this.browserPrint = BrowserPrint;
    } else {
      console.warn('BrowserPrint SDK not loaded');
    }
  }

  getPrinters(): Observable<any[]> {
    if (!this.browserPrint) {
      return throwError(() => new Error('BrowserPrint SDK not available'));
    }

    return from(this.browserPrint.getPrinters()).pipe(
      map((printers: any) => printers || []),
      catchError(error => {
        console.error('Error getting printers:', error);
        return throwError(() => new Error('Failed to get printers'));
      })
    );
  }

  getDefaultPrinter(): Observable<any> {
    if (!this.browserPrint) {
      return throwError(() => new Error('BrowserPrint SDK not available'));
    }

    return from(this.browserPrint.getDefaultPrinter()).pipe(
      catchError(error => {
        console.error('Error getting default printer:', error);
        return throwError(() => new Error('Failed to get default printer'));
      })
    );
  }

  printZPL(zplContent: string, printer?: any): Observable<any> {
    if (!this.browserPrint) {
      return throwError(() => new Error('BrowserPrint SDK not available'));
    }

    const targetPrinter = printer || this.browserPrint.getDefaultPrinter();
    
    return from(targetPrinter.send(zplContent)).pipe(
      map((result: any) => {
        console.log('Print result:', result);
        return result;
      }),
      catchError(error => {
        console.error('Error printing ZPL:', error);
        return throwError(() => new Error('Failed to print label'));
      })
    );
  }

  printPRN(prnContent: string, printer?: any): Observable<any> {
    if (!this.browserPrint) {
      return throwError(() => new Error('BrowserPrint SDK not available'));
    }

    const targetPrinter = printer || this.browserPrint.getDefaultPrinter();
    
    return from(targetPrinter.send(prnContent)).pipe(
      map((result: any) => {
        console.log('Print result:', result);
        return result;
      }),
      catchError(error => {
        console.error('Error printing PRN:', error);
        return throwError(() => new Error('Failed to print label'));
      })
    );
  }

  isBrowserPrintAvailable(): boolean {
    return typeof BrowserPrint !== 'undefined' && this.browserPrint !== null;
  }

  // Method to download ZPL as file
  downloadZPL(zplContent: string, filename: string = 'label.zpl'): void {
    const blob = new Blob([zplContent], { type: 'text/plain' });
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = filename;
    link.click();
    window.URL.revokeObjectURL(url);
  }

  // Method to download PRN as file
  downloadPRN(prnContent: string, filename: string = 'label.prn'): void {
    const blob = new Blob([prnContent], { type: 'application/octet-stream' });
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = filename;
    link.click();
    window.URL.revokeObjectURL(url);
  }

  // Test printer connection
  testPrinter(printer?: any): Observable<boolean> {
    if (!this.browserPrint) {
      return throwError(() => new Error('BrowserPrint SDK not available'));
    }

    const testZPL = '^XA^FO50,50^A0N,50,50^FDTest Print^FS^XZ';
    const targetPrinter = printer || this.browserPrint.getDefaultPrinter();
    
    return from(targetPrinter.send(testZPL)).pipe(
      map(() => true),
      catchError(error => {
        console.error('Printer test failed:', error);
        return throwError(() => new Error('Printer test failed'));
      })
    );
  }

  /**
   * Preview ZPL: convert ZPL string to PNG image data URL
   * @param zplContent ZPL command string
   * @returns Observable<string> base64 image URL
   */
  // previewZPL(zplContent: string): Observable<string> {
  //   try {
  //     const zplImage = new ZplImage({
  //       width: 384, // label width in dots
  //       height: 600, // label height in dots
  //       scale: 2 // scaling factor
  //     });
      
  //     const imageDataUrl = zplImage.render(zplContent);
  //     // imageDataUrl is a base64 PNG data URL string

  //     return of(imageDataUrl);
  //   } catch (error) {
  //     console.error('Error rendering ZPL preview:', error);
  //     return of('');
  //   }
  // }


} 