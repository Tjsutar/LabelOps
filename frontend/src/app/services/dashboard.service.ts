import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

export interface DashboardStats {
  overview: {
    total_labels: number;
    printed_labels: number;
    pending_labels: number;
    failed_labels: number;
    duplicate_labels: number;
  };
  breakdown: {
    by_grade: { [key: string]: number };
    by_section: { [key: string]: number };
  };
  activity: {
    recent_labels_24h: number;
    active_users_7d: number;
  };
  performance: {
    print_success_rate: string;
    total_print_jobs: number;
  };
  timestamp: string;
}

@Injectable({
  providedIn: 'root'
})
export class DashboardService {
  private apiUrl = `${environment.apiUrl}/dashboard`;

  constructor(private http: HttpClient) {}

  getDashboardStats(): Observable<DashboardStats> {
    return this.http.get<DashboardStats>(`${this.apiUrl}/stats`);
  }
}
