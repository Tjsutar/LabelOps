import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';
import { 
  LabelData, 
  Label, 
  LabelBatchRequest, 
  LabelBatchResponse, 
  LabelFilter,
  LabelStats 
} from '../models/label.model';

@Injectable({
  providedIn: 'root'
})
export class LabelService {
  constructor(private http: HttpClient) {}

  processBatch(batchRequest: LabelBatchRequest): Observable<LabelBatchResponse> {
    return this.http.post<LabelBatchResponse>(`${environment.apiUrl}/labels/batch`, batchRequest);
  }

  getLabels(filter?: LabelFilter): Observable<{ labels: Label[], count: number }> {
    let params = new HttpParams();
    
    if (filter) {
      if (filter.status) params = params.set('status', filter.status);
      if (filter.grade) params = params.set('grade', filter.grade);
      if (filter.section) params = params.set('section', filter.section);
      if (filter.heat_no) params = params.set('heat_no', filter.heat_no);
      if (filter.is_duplicate !== undefined) params = params.set('is_duplicate', filter.is_duplicate.toString());
      if (filter.limit) params = params.set('limit', filter.limit.toString());
      if (filter.offset) params = params.set('offset', filter.offset.toString());
    }

    return this.http.get<{ labels: Label[], count: number }>(`${environment.apiUrl}/labels`, { params });
  }

  getLabelById(id: string): Observable<Label> {
    return this.http.get<Label>(`${environment.apiUrl}/labels/${id}`);
  }

  printLabel(id: string): Observable<{ message: string, print_job_id: string, zpl_content: string }> {
    return this.http.post<{ message: string, print_job_id: string, zpl_content: string }>(
      `${environment.apiUrl}/labels/print`, 
      { id: id }
    );
  }

  exportLabelsCSV(filter?: LabelFilter): Observable<Blob> {
    let params = new HttpParams();
    
    if (filter) {
      if (filter.status) params = params.set('status', filter.status);
      if (filter.grade) params = params.set('grade', filter.grade);
      if (filter.section) params = params.set('section', filter.section);
      if (filter.heat_no) params = params.set('heat_no', filter.heat_no);
      if (filter.is_duplicate !== undefined) params = params.set('is_duplicate', filter.is_duplicate.toString());
    }

    console.log(params);

    return this.http.get(`${environment.apiUrl}/labels/export/csv`, { 
      params, 
      responseType: 'blob' 
    });
  }

  getPrintJobs(): Observable<any[]> {
    return this.http.get<any[]>(`${environment.apiUrl}/print-jobs`);
  }

  getPrintJobById(id: string): Observable<any> {
    return this.http.get<any>(`${environment.apiUrl}/print-jobs/${id}`);
  }

  retryPrintJob(id: string, ): Observable<any> {
    return this.http.post<any>(`${environment.apiUrl}/print-jobs/retry`, { job_id: id });
  }

  // Export print jobs as CSV
  exportPrintJobsCSV(status?: string): Observable<Blob> {
    let params = new HttpParams();
    
    if (status) {
      params = params.set('status', status);
    }

    return this.http.get(`${environment.apiUrl}/print-jobs/export/csv`, {
      params,
      responseType: 'blob'
    });
  }

  getStats(): Observable<LabelStats> {
    return this.http.get<LabelStats>(`${environment.apiUrl}/admin/stats`);
  }

  // Helper method to download CSV
  downloadCSV(blob: Blob, filename: string): void {
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = filename;
    link.click();
    window.URL.revokeObjectURL(url);
  }
} 