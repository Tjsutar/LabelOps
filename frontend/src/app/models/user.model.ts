export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  role: 'admin' | 'user' | 'operator';
  is_active: boolean;
  last_login?: string;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
  expires_at: number;
}

export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  role?: string;
}

export interface UserUpdateRequest {
  first_name?: string;
  last_name?: string;
  email?: string;
}

export interface AuditLog {
  id: string;
  user_id: string;
  action: string;
  resource: string;
  resource_id?: string;
  details: string;
  ip_address: string;
  user_agent: string;
  created_at: string;
  user_email?: string;
  user_name?: string;
}

export interface SystemStats {
  total_users: number;
  active_users: number;
  total_labels: number;
  printed_labels: number;
  pending_labels: number;
  failed_labels: number;
  total_print_jobs: number;
  failed_print_jobs: number;
} 