export interface JournalEntry {
  id: number;
  title: string;
  content: string;
  createdAt: string;
  updatedAt: string;
  tags: string[];
}

export interface ApiJournalEntry {
  id: number;
  title: string;
  content: string;
  created_at: string;
  updated_at: string;
  tags: string[];
}

export interface CreateEntryRequest {
  title: string;
  content: string;
}

export interface UpdateEntryRequest {
  title?: string;
  content?: string;
}