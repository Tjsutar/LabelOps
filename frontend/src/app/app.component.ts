import { Component } from "@angular/core";
import { CommonModule } from "@angular/common";
import { RouterOutlet } from "@angular/router";
import { HeaderComponent } from "./components/header/header.component";
import { ToastComponent } from "./components/toast/toast.component";

@Component({
  selector: "app-root",
  standalone: true,
  imports: [CommonModule, RouterOutlet, HeaderComponent, ToastComponent],
  template: `
    <div class="h-screen flex flex-col bg-gray-50">
      <div class="sticky top-0 z-50 bg-white shadow">
        <app-header></app-header>
      </div>

      <main
        class="flex-1 overflow-y-auto container mx-auto px-4 py-4"
        style="-ms-overflow-style: none; scrollbar-width: none;"
      >
        <router-outlet></router-outlet>
      </main>

      <app-toast></app-toast>
    </div>
  `,
  styles: [],
})
export class AppComponent {
  title = "LabelOps";
}
