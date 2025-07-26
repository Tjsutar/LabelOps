import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterOutlet } from '@angular/router';
import { HeaderComponent } from './components/header/header.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterOutlet, HeaderComponent],
  template: `
    <div class="min-h-screen bg-gray-50">
      <app-header></app-header>
      <main class="container mx-auto px-4 py-8">
        <router-outlet></router-outlet>
      </main>
    </div>
  `,
  styles: []
})
export class AppComponent {
  title = 'LabelOps';
} 