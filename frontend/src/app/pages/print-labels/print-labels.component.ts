import { Component } from "@angular/core";
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

  constructor(private http: HttpClient) {}

  async fetchJson() {
    // Use custom URL only if pathOption is 'custom' and jsonUrl is set
    const url =
      this.pathOption === "custom" && this.jsonUrl
        ? this.jsonUrl
        : "/assets/dummy_data.json";

    console.log("Fetching JSON from:", url);

    try {
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      console.log("Fetched JSON data:", data);
      this.jsonData = data;
      this.responseMessage = "✅ JSON data fetched successfully!";
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

    console.log("Sending payload:", jsonData);

    this.http
      .post("http://localhost:8080/api/v1/labels/batch", jsonData, { headers })
      .subscribe({
        next: (response) => {
          console.log("Upload successful:", response);
          this.responseMessage = "✅ Labels printed successfully!";
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
