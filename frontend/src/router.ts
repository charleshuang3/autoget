import { Router } from '@lit-labs/router';
import { LitElement, html } from 'lit';
import { customElement, state } from 'lit/decorators.js';
import { fetchIndexers, fetchIndexerCategories } from './utils/api.ts';

import './views/search_view';
import './views/indexer_view';
import './views/404_view';

@customElement('app-router')
export class AppRouter extends LitElement {
  @state()
  private indexers: string[] = [];

  async connectedCallback() {
    super.connectedCallback();
    this.fetchIndexers();
  }

  private async fetchIndexers() {
    this.indexers = await fetchIndexers();
  }

  private router = new Router(this, [
    {
      path: '/',
      render: () => html`<div>Loading...</div>`,
      enter: async () => {
        if (this.indexers.length === 0) {
          await this.fetchIndexers();
        }
        const newUrl = `/indexers/${this.indexers[0]}`;
        this.router.goto(newUrl);
        history.replaceState(null, '', newUrl);
        return false;
      },
    },
    { path: '/search', render: () => html`<search-view></search-view>` },
    {
      path: '/indexers/:id',
      render: ({ id }) => {
        return html`<indexer-view .indexerId=${id || ''} category=""></indexer-view>`;
      },
      enter: async ({ id }) => {
        if (this.indexers.length === 0) {
          await this.fetchIndexers();
        }
        if (id === undefined || !this.indexers.includes(id)) {
          this.router.goto('/404');
          return false;
        }
        const categories = await fetchIndexerCategories(id);
        this.router.goto(`/indexers/${id}/${categories[0].id}`);
        history.replaceState(null, '', `/indexers/${id}/${categories[0].id}`);
        return false;
      },
    },
    {
      path: '/indexers/:id/:category',
      render: ({ id, category }) => {
        return html`<indexer-view .indexerId=${id || ''} .category=${category || ''}></indexer-view>`;
      },
    },
    { path: '*', render: () => html`<not-found-view></not-found-view>` },
  ]);

  render() {
    return this.router.outlet();
  }
}
