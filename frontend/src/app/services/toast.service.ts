import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

export interface Toast {
  id: string;
  message: string;
  type: 'success' | 'error' | 'info' | 'warning';
  duration?: number;
}

@Injectable({
  providedIn: 'root'
})
export class ToastService {
  private toasts = new BehaviorSubject<Toast[]>([]);
  toasts$ = this.toasts.asObservable();

  show(message: string, type: 'success' | 'error' | 'info' | 'warning' = 'info', duration: number = 5000) {
    const toast: Toast = {
      id: Date.now().toString(),
      message,
      type,
      duration
    };

    const currentToasts = this.toasts.value;
    this.toasts.next([...currentToasts, toast]);

    if (duration > 0) {
      setTimeout(() => {
        this.remove(toast.id);
      }, duration);
    }
  }

  success(message: string, duration?: number) {
    this.show(message, 'success', duration);
  }

  error(message: string, duration?: number) {
    this.show(message, 'error', duration);
  }

  info(message: string, duration?: number) {
    this.show(message, 'info', duration);
  }

  warning(message: string, duration?: number) {
    this.show(message, 'warning', duration);
  }

  remove(id: string) {
    const currentToasts = this.toasts.value;
    this.toasts.next(currentToasts.filter(toast => toast.id !== id));
  }

  clear() {
    this.toasts.next([]);
  }
} 