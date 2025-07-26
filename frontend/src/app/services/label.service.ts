import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';
import { 
  TMTBarData, 
  TMTBarLabel, 
  TMTBarBatchRequest, 
  TMTBarBatchResponse, 
  TMTBarFilter,
  TMTBarStats 
} from '../models/tmt-bar.model';

@Injectable({
  providedIn: 'root'
})
export class LabelService {
  constructor(private http: HttpClient) {}

  processBatch(batchRequest: TMTBarBatchRequest): Observable<TMTBarBatchResponse> {
    return this.http.post<TMTBarBatchResponse>(`${environment.apiUrl}/labels/batch`, batchRequest);
  }

  getLabels(filter?: TMTBarFilter): Observable<{ labels: TMTBarLabel[], count: number }> {
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

    return this.http.get<{ labels: TMTBarLabel[], count: number }>(`${environment.apiUrl}/labels`, { params });
  }

  getLabelById(id: string): Observable<TMTBarLabel> {
    return this.http.get<TMTBarLabel>(`${environment.apiUrl}/labels/${id}`);
  }

  printLabel(id: string): Observable<{ message: string, print_job_id: string, zpl_content: string }> {
    return this.http.post<{ message: string, print_job_id: string, zpl_content: string }>(
      `${environment.apiUrl}/labels/${id}/print`, 
      {}
    );
  }

  exportLabelsCSV(filter?: TMTBarFilter): Observable<Blob> {
    let params = new HttpParams();
    
    if (filter) {
      if (filter.status) params = params.set('status', filter.status);
      if (filter.grade) params = params.set('grade', filter.grade);
      if (filter.section) params = params.set('section', filter.section);
      if (filter.heat_no) params = params.set('heat_no', filter.heat_no);
      if (filter.is_duplicate !== undefined) params = params.set('is_duplicate', filter.is_duplicate.toString());
    }

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

  retryPrintJob(id: string): Observable<any> {
    return this.http.post<any>(`${environment.apiUrl}/print-jobs/${id}/retry`, {});
  }

  getStats(): Observable<TMTBarStats> {
    return this.http.get<TMTBarStats>(`${environment.apiUrl}/admin/stats`);
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