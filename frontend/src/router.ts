import { Router } from '@lit-labs/router';
import { LitElement, html } from 'lit';
import { customElement, state } from 'lit/decorators.js';
import { fetchIndexers } from './utils/api.ts';

import './views/search_view';
import './views/indexer_view';

@customElement('app-router')
export class AppRouter extends LitElement {
  @state()
  private _indexers: string[] = [];

  async connectedCallback() {
    super.connectedCallback();
    this.fetchIndexers();
  }

  private async fetchIndexers() {
    this._indexers = await fetchIndexers();
  }

  private router = new Router(this, [
    {
      path: '/',
      render: () => html`<div>Loading...</div>`,
      enter: async () => {
        if (this._indexers.length === 0) {
          await this.fetchIndexers();
        }
        const newUrl = `/indexers/${this._indexers[0]}`;
        this.router.goto(newUrl);
        history.replaceState(null, '', newUrl);
        return false;
      },
    },
    { path: '/search', render: () => html`<search-view></search-view>` },
    {
      path: '/indexers/:id',
      render: ({ id }) => {
        return html`<indexer-view .indexerId=${id || ''}></indexer-view>`;
      },
    },
    {
      path: '/indexers/:id/:category',
      render: ({ id, category }) => {
        return html`<indexer-view .indexerId=${id || ''} .category=${category || ''}></indexer-view>`;
      },
    },
  ]);

  render() {
    return this.router.outlet();
  }
}
