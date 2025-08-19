import { Router } from '@lit-labs/router';
import { LitElement, html } from 'lit';
import { customElement } from 'lit/decorators.js';

import './views/search_view';
import './views/indexer_view';

@customElement('app-router')
export class AppRouter extends LitElement {
  private router = new Router(this, [
    { path: '/search', render: () => html`<search-view></search-view>` },
    { path: '/indexers/:id', render: ({ id }) => html`<indexer-view index_id=${id}></indexer-view>` },
  ]);

  render() {
    return this.router.outlet();
  }
}
