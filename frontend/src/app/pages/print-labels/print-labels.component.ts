import { Component, EventEmitter, Output } from "@angular/core";
import { CommonModule } from "@angular/common";
import { FormsModule } from "@angular/forms";
import { HttpClient, HttpHeaders } from "@angular/common/http";

@Component({
  selector: "app-print-labels",
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: "./print-labels.html",
})
export class PrintLabelsComponent {
  pathOption = "default";
  jsonUrl = "";
  jsonData: any = null;
  responseMessage = "";
  @Output() pageSizeChange = new EventEmitter<number>();
  @Output() refresh = new EventEmitter<void>();

  constructor(private http: HttpClient) {}

  private mapIncomingData(data: any[]): any[] {
    if (!Array.isArray(data)) return [];
    return data.map((item: any) => {
      // Map DATE1 -> DATE; keep other fields as-is
      const mapped: any = { ...item };
      if (mapped.DATE1 && !mapped.DATE) {
        mapped.DATE = mapped.DATE1;
      }
      return mapped;
    });
  }

  async fetchJson() {
    // Use custom URL only if pathOption is 'custom' and jsonUrl is set
    const url =
      this.pathOption === "custom" && this.jsonUrl
        ? this.jsonUrl
        : "/assets/mmnew.json";

    console.log("Fetching JSON from:", url);
    

    try {
      const response = await fetch(url);
      console.log("Full API response (raw):", JSON.stringify(response, null, 2));
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const raw = await response.json();
      console.log("Fetched JSON data:", raw);
      const data = this.mapIncomingData(raw);
      this.jsonData = data;
      if (!Array.isArray(data) || data.length === 0) {
        this.responseMessage = "❌ No data found in JSON file";
        return null;
      }
      this.responseMessage = "✅ JSON data fetched successfully!";
      // Emit page size based on input JSON length
      if (Array.isArray(data)) {
        try {
          this.pageSizeChange.emit(data.length);
          console.log("Emitted pageSizeChange:", data.length);
        } catch (e) {
          console.warn("Failed to emit pageSizeChange", e);
        }
      }
      return data;
    } catch (error: any) {
      console.error("Fetch error:", error);
      this.responseMessage = `❌ Failed to fetch JSON: ${error.message}`;
      return null;
    }
  }

  uploadJson(jsonData: any) {
    console.log("Starting uploadJson()");
  
    if (!jsonData) {
      console.warn("No JSON data loaded - aborting upload");
      this.responseMessage = "❌ No JSON data loaded - aborting upload";
      return;
    }
  
    const token = localStorage.getItem("token");
    if (!token) {
      console.warn("No auth token found");
      this.responseMessage = "❌ No auth token found";
      return;
    }
  
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    });
  
    // Backend expects { labels: [...] }
    const payload = { labels: Array.isArray(jsonData) ? jsonData : [] };

    console.log("Sending payload:", payload);

    this.http
      .post<any>("http://localhost:8080/api/v1/labels/batch", payload, { headers })
      .subscribe({
        next: (response) => {
          console.log("Upload successful:", response);
  
          const dupCount = Number(response?.duplicate_count ?? 0);
          const newCount = Number(response?.new_count ?? 0);
          const totalProcessed = Number(response?.total_processed ?? 0);
  
          // ✅ Decide the message based on API data
          if (dupCount === 0 && newCount > 0) {
            this.responseMessage = `✅ Labels printed successfully! (${newCount} new)`;
          } else if (dupCount > 0 && newCount === 0) {
            this.responseMessage = `ℹ️ All ${dupCount} labels were already printed.`;
          } else if (dupCount > 0 && newCount > 0) {
            this.responseMessage = `⚠️ ${dupCount} labels already printed, ${newCount} new labels printed.`;
          } else if (totalProcessed === 0) {
            this.responseMessage = `❌ No labels processed.`;
          } else {
            this.responseMessage = response?.message || "✅ Operation completed.";
          }

          // Notify parent to reload labels list
          try {
            this.refresh.emit();
          } catch (e) {
            console.warn("Failed to emit refresh event", e);
          }
        },
        error: (err) => {
          console.error("Upload failed:", err);
          this.responseMessage = "❌ Submission failed";
        },
      });
  }
  
  
  async fetchAndUpload() {
    const data = await this.fetchJson();
    if (data) {
      this.uploadJson(data);
    }
    
  }
}
