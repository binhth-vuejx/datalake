export interface Iframe {
  id: string;
  workspace_id: string;
  title: string;
  description: string | null;
  icon: string | null;
  iframe_url: string | null;
  iframe_script: string | null;
  created_at: string;
  updated_at: string;
}

export interface CreateIframeRequest {
  title: string;
  description?: string;
  icon?: string;
  iframe_url?: string;
  iframe_script?: string;
}

export interface UpdateIframeRequest {
  title?: string;
  description?: string | null;
  icon?: string | null;
  iframe_url?: string | null;
  iframe_script?: string | null;
}

export interface ListIframesResponse {
  iframes: Iframe[];
  total: number;
}
