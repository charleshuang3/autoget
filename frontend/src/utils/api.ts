export interface Category {
  id: string;
  name: string;
  subCategories: Category[];
}

export async function fetchIndexers(): Promise<string[]> {
  try {
    const response = await fetch('/api/v1/indexers');
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Failed to fetch indexers:', error);
    return []; // Set to empty array on error
  }
}

export async function fetchIndexerCategories(indexer: string): Promise<Category[]> {
  try {
    const response = await fetch(`/api/v1/indexers/${indexer}/categories`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Failed to fetch indexers:', error);
    return []; // Set to empty array on error
  }
}
