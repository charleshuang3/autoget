import { Router } from '@lit-labs/router';
import { LitElement, html } from 'lit';
import { customElement } from 'lit/decorators.js';

import './views/search_view';
import './views/indexer_view';

@customElement('app-router')
export class AppRouter extends LitElement {
  private router = new Router(this, [
    {
      path: '/',
      render: () => html`<div>Loading...</div>`,
      enter: async () => {
        let newUrl = "/search";
        try {
          const response = await fetch('/api/v1/indexers');
          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
          }
          const indexers = await response.json();
          if (indexers && indexers.length > 0) {
            newUrl = `/indexers/${indexers[0]}`;
          }
        } catch (error) {
          console.error('Failed to fetch indexers, redirecting to search', error);
        }
        this.router.goto(newUrl);
        history.replaceState(null, '', newUrl);
        return false;
      },
    },
    { path: '/search', render: () => html`<search-view></search-view>` },
    {
      path: '/indexers/:id',
      render: ({ id }) => html`<indexer-view .index_id=${id}></indexer-view>`,
    },
  ]);

  render() {
    return this.router.outlet();
  }
}
