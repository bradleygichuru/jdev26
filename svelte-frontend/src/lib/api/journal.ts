import type { JournalEntry, CreateEntryRequest, UpdateEntryRequest, ApiJournalEntry } from '$lib/types/journal';

const API_BASE = 'http://localhost:8080/api';

class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = 'ApiError';
  }
}

async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE}${endpoint}`;

  const config: RequestInit = {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  };

  try {
    const response = await fetch(url, config);

    if (!response.ok) {
      throw new ApiError(response.status, `HTTP error! status: ${response.status}`);
    }

    return await response.json();
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }
    throw new ApiError(0, 'Network error or failed to parse response');
  }
}

function transformApiEntry(apiEntry: ApiJournalEntry): JournalEntry {
  return {
    id: apiEntry.id,
    title: apiEntry.title,
    content: apiEntry.content,
    createdAt: apiEntry.created_at,
    updatedAt: apiEntry.updated_at,
    tags: apiEntry.tags || []
  };
}

export const journalApi = {
  // Get all journal entries
  async getEntries(): Promise<JournalEntry[]> {
    const response = await apiRequest<{success: boolean; data: ApiJournalEntry[]}>('/entries');
    return response.data.map(transformApiEntry);
  },

  // Get a single journal entry by ID
  async getEntry(id: number): Promise<JournalEntry> {
    const response = await apiRequest<{success: boolean; data: ApiJournalEntry}>(`/entries/${id}`);
    return transformApiEntry(response.data);
  },

  // Create a new journal entry
  async createEntry(entry: CreateEntryRequest): Promise<JournalEntry> {
    const response = await apiRequest<{success: boolean; data: ApiJournalEntry}>('/entries', {
      method: 'POST',
      body: JSON.stringify(entry),
    });
    return transformApiEntry(response.data);
  },

  // Update an existing journal entry
  async updateEntry(id: number, updates: UpdateEntryRequest): Promise<JournalEntry> {
    const response = await apiRequest<{success: boolean; data: ApiJournalEntry}>(`/entries/${id}`, {
      method: 'PUT',
      body: JSON.stringify(updates),
    });
    return transformApiEntry(response.data);
  },

  // Delete a journal entry
  async deleteEntry(id: number): Promise<void> {
    const response = await apiRequest<{success: boolean; message?: string}>(`/entries/${id}`, {
      method: 'DELETE',
    });
    // API might return success/message for delete, but we don't need to process it
  },
};