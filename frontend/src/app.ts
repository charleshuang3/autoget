import { html, LitElement } from 'lit';
import { provide } from '@lit/context';
import { customElement, property } from 'lit/decorators.js';
import { indexersContext } from './context.ts';

import './router.ts';

@customElement('app-root')
export class App extends LitElement {
  @provide({ context: indexersContext })
  @property({ attribute: false })
  public indexers: string[] = [];

  connectedCallback() {
    super.connectedCallback();
    this._fetchIndexers();
  }

  private async _fetchIndexers() {
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

  render() {
    return html` <app-router></app-router> `;
  }
}
