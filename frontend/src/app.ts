import { html, LitElement } from 'lit';
import { provide } from '@lit/context';
import { customElement, property } from 'lit/decorators.js';
import { indexersContext, indexerDetailsContext, type IndexerDetails, type Category } from './context.ts';

import './router.ts';

@customElement('app-root')
export class App extends LitElement {
  @provide({ context: indexersContext })
  @property({ attribute: false })
  public indexers: string[] = [];

  @provide({ context: indexerDetailsContext })
  @property({ attribute: false })
  public indexerDetails: IndexerDetails = {
    categories: (indexer: string): Promise<Category[]> => {
      return this.fetchIndexerCategories(indexer);
    },
  };

  private catchedIndexerCategories: Map<string, Category[]> = new Map();

  connectedCallback() {
    super.connectedCallback();
    this.fetchIndexers();
  }

  private async fetchIndexers() {
    try {
      const response = await fetch('/api/v1/indexers');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      this.indexers = await response.json();
    } catch (error) {
      console.error('Failed to fetch indexers:', error);
      this.indexers = []; // Set to empty array on error
    }
  }

  private async fetchIndexerCategories(indexer: string): Promise<Category[]> {
    if (this.catchedIndexerCategories.has(indexer)) {
      return this.catchedIndexerCategories.get(indexer) || [];
    }

    try {
      const response = await fetch(`/api/v1/indexers/${indexer}/categories`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const categories = await response.json();
      this.catchedIndexerCategories.set(indexer, categories);
      return categories;
    } catch (error) {
      console.error('Failed to fetch indexers:', error);
      return []; // Set to empty array on error
    }
  }

  render() {
    return html` <app-router></app-router> `;
  }
}
