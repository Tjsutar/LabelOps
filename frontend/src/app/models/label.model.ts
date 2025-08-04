export interface LabelData {
  LOCATION?: string;
  BUNDLE_NO: number;
  PQD: string;
  UNIT: string;
  TIME: string;
  LENGTH: number;
  HEAT_NO: string;
  PRODUCT_HEADING: string;
  ISI_BOTTOM: string;
  ISI_TOP: string;
  CHARGE_DTM: string;
  MILL: string;
  GRADE: string;
  URL_APIKEY: string;
  ID?: string;
  WEIGHT?: string;
  SECTION: string;
  DATE1: string;
}

export interface Label {
  id: string;
  label_id: string;
  location?: string;
  bundle_no: string;
  bundle_type: string;
  pqd: string;
  unit: string;
  time: string;
  length: number;
  heat_no: string;
  product_heading: string;
  isi_bottom: string;
  isi_top: string;
  charge_dtm: string;
  mill: string;
  grade: string;
  url_apikey: string;
  weight?: string;
  section: string;
  date: string;
  user_id: string;
  status: 'pending' | 'printed' | 'failed';
  is_duplicate: boolean;
  created_at: string;
  updated_at: string;
}

export interface LabelBatchRequest {
  labels: LabelData[];
}

export interface LabelBatchResponse {
  new_labels: Label[];
  duplicate_labels: Label[];
  total_processed: number;
  new_count: number;
  duplicate_count: number;
}

export interface LabelFilter {
  status?: string;
  grade?: string;
  section?: string;
  heat_no?: string;
  is_duplicate?: boolean;
  limit?: number;
  offset?: number;
}

export interface LabelStats {
  total_labels: number;
  printed_labels: number;
  pending_labels: number;
  failed_labels: number;
  duplicate_labels: number;
  by_grade: Record<string, number>;
  by_section: Record<string, number>;
} 